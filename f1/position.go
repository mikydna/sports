package f1

import (
	"sort"
	"strings"

	"github.com/mikydna/sports/f1/livetiming"
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
						int32(srcPosition.Y)},
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
