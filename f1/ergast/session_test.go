package ergast_test

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/mikydna/sports/f1/ergast"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestErgast_Sessions(t *testing.T) {
	f, err := os.Open("./_testdata/session-2019.json")
	require.NoError(t, err)
	defer f.Close()
	testdata, err := ioutil.ReadAll(f)
	require.NoError(t, err)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write(testdata)
		require.NoError(t, err)
		assert.True(t, strings.Contains(r.URL.Path, "/f1/"))
		assert.True(t, strings.HasSuffix(r.URL.Path, "/2019.json"))
	}))
	defer server.Close()

	ctx := context.TODO()
	c := ergast.NewClient(http.DefaultClient, server.URL)

	sessions, err := c.Sessions(ctx, "f1", 2019)
	assert.NoError(t, err)
	assert.Equal(t, sessions.MRData.Total, "21")
}
