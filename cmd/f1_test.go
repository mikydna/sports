package cmd_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/mikydna/sports/cmd"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestCmd_ParseF1Config(t *testing.T) {
	dir, err := os.MkdirTemp(os.TempDir(), "TestCmd_ParseF1Config-")
	require.NoError(t, err)
	defer os.RemoveAll(dir)

	repo := filepath.Join(dir, "repo")
	err = os.Mkdir(repo, 0777)
	require.NoError(t, err)

	workspace := filepath.Join(dir, "workspace")
	err = os.Mkdir(workspace, 0777)
	require.NoError(t, err)

	fp := tempF1Config(t, dir, "test-f1.yml", &cmd.F1Config{
		Repo:      repo,
		Workspace: workspace,
	})

	cfg, err := cmd.ParseF1Config(fp)
	assert.NoError(t, err)
	assert.Equal(t, cfg.Repo, repo)
	assert.NoError(t, cmd.ValidateF1Config(cfg))
}

func tempF1Config(t *testing.T, dir, name string, cfg *cmd.F1Config) string {
	f, err := os.CreateTemp(dir, name)
	require.NoError(t, err)
	err = yaml.NewEncoder(f).Encode(cfg)
	require.NoError(t, err)
	return f.Name()
}
