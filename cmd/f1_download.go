package cmd

import (
	"context"
	"net/http"

	"github.com/mikydna/sports/f1"
	"github.com/spf13/cobra"
)

var (
	flagF1DownloadSeasons []int
)

func init() {
	f1CmdDownload.Flags().IntSliceVar(&flagF1DownloadSeasons, "season", []int{2021}, "seasons")
}

var f1CmdDownload = &cobra.Command{
	Use: "download",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		var err error
		if f1Config, err = ParseF1Config(flagF1Config); err != nil {
			return err
		}
		if err := ValidateF1Config(f1Config); err != nil {
			return err
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := context.WithCancel(cmd.Context())
		interrupt(cancel)
		defer cancel()

		c := f1.NewClient(http.DefaultClient, w, f1Config.Repo, f1Config.Workspace)
		cfg := &f1.DownloadConfig{
			Seasons: flagF1DownloadSeasons,
		}

		if err := c.Download(ctx, cfg); err != nil {
			return err
		}

		return nil
	},
}
