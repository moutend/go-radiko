package app

import (
	"github.com/moutend/go-radiko/pkg/radiko"
	"github.com/spf13/cobra"
)

var versionCommand = &cobra.Command{
	Use:     "version",
	Aliases: []string{"v"},
	Short:   "print version",
	RunE:    versionCommandRunE,
}

func versionCommandRunE(cmd *cobra.Command, args []string) error {
	cmd.Println(radiko.Version())

	return nil
}

func init() {
	RootCommand.AddCommand(versionCommand)
}
