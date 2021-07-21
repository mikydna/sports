package livetiming_test

import (
	"bufio"
	"io"
	"os"
	"testing"

	"github.com/mikydna/sports/f1/livetiming"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLivetiming_PositionFrame_Unmarshal(t *testing.T) {
	f, err := os.Open("./_testdata/Position.z.jsonStream")
	require.NoError(t, err)
	defer f.Close()

	r := bufio.NewReader(f)
	for i := 0; ; i++ {
		b, _, err := r.ReadLine()
		if err == io.EOF {
			break
		}
		require.NoError(t, err)

		position := new(livetiming.PositionFrame)
		err = position.Unmarshal(b)
		assert.NoError(t, err)

		// position testdata is between +1hr and +2hr
		assert.True(t, position.Offset.Hours() >= 1 && position.Offset.Hours() < 2)
	}
}

func TestLivetiming_PositionFrames_Unmarshal(t *testing.T) {
	b, err := os.ReadFile("./_testdata/Position.z.jsonStream")
	require.NoError(t, err)

	var frames livetiming.PositionFrames
	err = frames.Unmarshal(b)
	require.NoError(t, err)
	assert.Len(t, frames, 99)
	assert.Less(t, frames[0].Offset, frames[len(frames)-1].Offset) // sort asc
}
