package pb_test

import (
	"bytes"
	"testing"

	"github.com/mikydna/sports/f1/pb"
	"github.com/stretchr/testify/assert"
)

func TestPB_WriteRead(t *testing.T) {
	frame := &pb.PositionFrame{
		Offset:    1,
		Timestamp: 2,
		Position: map[uint32]*pb.Position{
			99999: {Xyz: 123456, Status: pb.PositionStatus_OFFTRACK},
		},
	}

	// write
	buf := bytes.NewBuffer([]byte{})
	err := pb.Write(buf, frame)
	assert.NoError(t, err)

	// read back
	newFrame := new(pb.PositionFrame)
	err = pb.Read(bytes.NewBuffer(buf.Bytes()), newFrame)
	assert.NoError(t, err)
	assert.Equal(t, frame.Offset, newFrame.Offset)
	assert.Equal(t, frame.Timestamp, newFrame.Timestamp)
	assert.Equal(t, frame.Position[99999].Xyz, newFrame.Position[99999].Xyz)
	assert.Equal(t, frame.Position[99999].Status, newFrame.Position[99999].Status)
}
