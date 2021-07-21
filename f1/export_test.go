package f1_test

import (
	"context"
	"io/ioutil"
	"os"
	"testing"

	"github.com/mikydna/sports/f1"
	"github.com/mikydna/sports/f1/livetiming"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestF1_Export_PositionFrames(t *testing.T) {
	b, err := os.ReadFile("./livetiming/_testdata/Position.z.jsonStream")
	require.NoError(t, err)

	var livetimingFrames livetiming.PositionFrames
	err = livetimingFrames.Unmarshal(b)
	require.NoError(t, err)

	var frames f1.PositionFrames
	err = frames.FromLivetiming(livetimingFrames)
	require.NoError(t, err)

	tmp, err := ioutil.TempDir(os.TempDir(), "TestF1_Export_PositionFrames-")
	require.NoError(t, err)
	s := f1.NewExportService(tmp)

	t.Run("JSON", func(t *testing.T) {
		ctx := context.TODO()
		exportedFile, err := s.Export(ctx, "/foo", frames, f1.ExportFormatJSON)
		assert.NoError(t, err)

		stat, err := os.Stat(exportedFile)
		assert.NoError(t, err)
		assert.Equal(t, "Position.jsonl", stat.Name())
	})

	t.Run("Gob", func(t *testing.T) {
		ctx := context.TODO()
		exportedFile, err := s.Export(ctx, "/foo/", frames, f1.ExportFormatGob)
		assert.NoError(t, err)

		stat, err := os.Stat(exportedFile)
		assert.NoError(t, err)
		assert.Equal(t, "Position.gob", stat.Name())
	})

	t.Run("Protobuf", func(t *testing.T) {
		ctx := context.TODO()
		exportedFile, err := s.Export(ctx, "/foo/", frames, f1.ExportFormatProto)
		assert.NoError(t, err)

		stat, err := os.Stat(exportedFile)
		assert.NoError(t, err)
		assert.Equal(t, "Position.pb", stat.Name())
	})

}
