package livetiming_test

import (
	"bufio"
	"encoding/json"
	"io"
	"os"
	"testing"

	"github.com/mikydna/sports/f1/livetiming"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLivetiming_DecodeLine(t *testing.T) {
	t.Run("Compressed", func(t *testing.T) {
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

			ts, data, err := livetiming.DecodeLine(b, true)
			assert.NoError(t, err)
			assert.True(t, json.Valid(data))

			// position testdata is between +1hr and +2hr
			assert.True(t, ts.Hours() >= 1 && ts.Hours() < 2)
		}
	})

	t.Run("NotCompressed", func(t *testing.T) {
		f, err := os.Open("./_testdata/LapCount.jsonStream")
		require.NoError(t, err)
		defer f.Close()

		r := bufio.NewReader(f)
		for i := 0; ; i++ {
			b, _, err := r.ReadLine()
			if err == io.EOF {
				break
			}
			require.NoError(t, err)

			ts, data, err := livetiming.DecodeLine(b, false)
			assert.NoError(t, err)
			assert.True(t, json.Valid(data))

			// position testdata is between 0 and +2hr
			assert.True(t, ts.Seconds() > 0 && ts.Hours() < 2)
		}
	})
}
