package cmd

import (
	"io"
	"os"

	"github.com/spf13/cobra"
)

var (
	w io.Writer = os.Stderr
)

func init() {
	Root.AddCommand(f1Cmd)
}

var Root = &cobra.Command{
	Use:     "gl",
	Aliases: []string{"sports"},
}
