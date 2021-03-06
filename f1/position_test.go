package f1_test

import (
	"os"
	"testing"

	"github.com/mikydna/sports/f1"
	"github.com/mikydna/sports/f1/livetiming"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestF1_FromLivetimingPositionFrames(t *testing.T) {
	b, err := os.ReadFile("./livetiming/_testdata/Position.z.jsonStream")
	require.NoError(t, err)

	var livetimingFrames livetiming.PositionFrames
	err = livetimingFrames.Unmarshal(b)
	require.NoError(t, err)
	require.Len(t, livetimingFrames, 99)

	var frames f1.PositionFrames
	err = frames.FromLivetiming(livetimingFrames)
	assert.NoError(t, err)
	assert.Len(t, frames, 466)
	assert.Len(t, frames[0].Position, 20)

	// &{3 [-7276 5671 5671] 1}
	assert.Equal(t, frames[0].Position[0].CarID, uint32(3))
	assert.Equal(t, frames[0].Position[0].XYZ, [3]int32{-7276, 5671, 5671})
	assert.Equal(t, frames[0].Position[0].Status, f1.PositionStatusOnTrack)
}
