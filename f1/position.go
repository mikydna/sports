package f1

import (
	"bytes"
	"math"
	"sort"
	"strings"

	"github.com/mikydna/sports/f1/livetiming"
	"github.com/mikydna/sports/f1/pb"
	"github.com/mikydna/sports/lib"
)

type PositionStatus uint8

const (
	PositionStatusUnknown  PositionStatus = iota // unknown
	PositionStatusOnTrack                        // ontrack
	PositionStatusOffTrack                       // offtrack
)

type Position struct {
	CarID  uint32         `json:"carID"`
	XYZ    [3]int32       `json:"xyz"`
	Status PositionStatus `json:"status"`
}

type PositionFrame struct {
	Offset    int32       `json:"offset"`
	Timestamp int64       `json:"t"`
	Position  []*Position `json:"position"`
}

type PositionFrames []*PositionFrame

func (pfs *PositionFrames) FromLivetiming(src livetiming.PositionFrames) error {
	for _, srcFrame := range src {
		offset := srcFrame.Offset
		for _, curr := range srcFrame.Position {
			positions := []*Position{}
			for carID, srcPosition := range curr.Entries {
				positions = append(positions, &Position{
					CarID:  uint32(carID),
					Status: parseLivetimingStatus(srcPosition.Status),
					XYZ: [3]int32{
						int32(srcPosition.X),
						int32(srcPosition.Y),
						int32(srcPosition.Z)},
				})

				sort.SliceStable(positions, func(i, j int) bool {
					return positions[i].CarID < positions[j].CarID
				})
			}

			*pfs = append(*pfs, &PositionFrame{
				Offset:    int32(offset.Milliseconds()),
				Timestamp: curr.Timestamp.Unix(),
				Position:  positions,
			})
		}
	}

	return nil
}

func (pfs PositionFrames) Proto() ([]*pb.PositionFrame, error) {

	// note: morton encoded values need to be uints
	// find the min of { x, y, z } of all frames
	rangeX := [2]int32{math.MaxInt32, math.MinInt32}
	rangeY := [2]int32{math.MaxInt32, math.MinInt32}
	rangeZ := [2]int32{math.MaxInt32, math.MinInt32}

	for _, frame := range pfs {
		for _, curr := range frame.Position {
			rangeX = lib.MinMaxInt32(rangeX, curr.XYZ[0])
			rangeY = lib.MinMaxInt32(rangeY, curr.XYZ[1])
			rangeZ = lib.MinMaxInt32(rangeZ, curr.XYZ[2])
		}
	}

	// encode
	protoFrames := make([]*pb.PositionFrame, len(pfs))
	for i, frame := range pfs {
		protoFrame := &pb.PositionFrame{
			Offset:    frame.Offset,
			Timestamp: frame.Timestamp,
			Position:  map[uint32]*pb.Position{},
		}

		for _, position := range frame.Position {
			protoFrame.Position[position.CarID] = &pb.Position{
				Xyz: lib.EncodeInt32(
					position.XYZ[0]-rangeX[0],
					position.XYZ[1]-rangeX[0],
					position.XYZ[2]-rangeX[0],
				),
				Status: convertToProtobufStatus(position.Status),
			}
		}

		protoFrames[i] = protoFrame
	}

	return protoFrames, nil
}

func (pfs PositionFrames) ProtoBytes() ([]byte, error) {
	protoFrames, err := pfs.Proto()
	if err != nil {
		return nil, err
	}

	buf := bytes.NewBuffer([]byte{})
	for _, curr := range protoFrames {
		if err := pb.Write(buf, curr); err != nil {
			return nil, err
		}
	}

	return buf.Bytes(), nil
}

func convertToProtobufStatus(x PositionStatus) pb.PositionStatus {
	switch x {
	case PositionStatusOnTrack:
		return pb.PositionStatus_ONTRACK
	case PositionStatusOffTrack:
		return pb.PositionStatus_OFFTRACK
	default:
		return pb.PositionStatus_UNKNOWN
	}
}

func parseLivetimingStatus(str string) PositionStatus {
	var status PositionStatus
	switch strings.ToLower(str) {
	case "ontrack":
		status = PositionStatusOnTrack
	case "offtrack":
		status = PositionStatusOffTrack
	default:
		status = PositionStatusUnknown
	}
	return status
}
