package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

const (
	DefaultF1ConfigFile = "./f1.yml"
)

var (
	f1Config *F1Config

	flagF1Config string

	ErrBadF1ConfigFunc = func(x interface{}) error {
		return fmt.Errorf("bad f1 config: %v", x)
	}
)

type F1Config struct {
	Repo      string `yaml:"repo"`
	Workspace string `yaml:"workspace"`
}

func init() {
	f1Cmd.PersistentFlags().StringVar(&flagF1Config, "config", DefaultF1ConfigFile, "config")
	f1Cmd.MarkPersistentFlagFilename("config") // nolint
	f1Cmd.AddCommand(f1CmdDownload)
	f1Cmd.AddCommand(f1CmdRepo)
}

var f1Cmd = &cobra.Command{
	Use: "f1",
}

func ParseF1Config(fp string) (*F1Config, error) {
	f, err := os.Open(fp)
	if err != nil {
		return nil, ErrBadF1ConfigFunc(err)
	}

	var cfg *F1Config
	if err := yaml.NewDecoder(f).Decode(&cfg); err != nil {
		return nil, ErrBadF1ConfigFunc(err)
	}

	return cfg, nil
}

func ValidateF1Config(cfg *F1Config) error {
	if !dirExists(cfg.Repo) {
		return ErrBadF1ConfigFunc("invalid repo")
	}
	if !dirExists(cfg.Workspace) {
		return ErrBadF1ConfigFunc("invalid workspace")
	}
	return nil
}
