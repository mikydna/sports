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

	ctx := context.TODO()
	err = s.Export(ctx, "/foo/", frames, f1.ExportFormatProto)
	assert.NoError(t, err)

}
