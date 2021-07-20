package cmd

import (
	"context"
	"net/http"

	"github.com/mikydna/sports/f1"
	"github.com/spf13/cobra"
)

var ()

func init() {

}

var f1CmdRepo = &cobra.Command{
	Use: "list",
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

		if err := c.List(ctx, w); err != nil {
			return err
		}

		return nil
	},
}
