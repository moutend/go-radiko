package app

import (
	"github.com/moutend/go-radiko/pkg/radiko"
	"github.com/spf13/cobra"
)

var stationCommand = &cobra.Command{
	Use:     "station",
	Aliases: []string{"s"},
	Short:   "list all available stations",
	RunE:    stationCommandRunE,
}

func stationCommandRunE(cmd *cobra.Command, args []string) error {
	stations, err := radiko.GetStations()

	if err != nil {
		return err
	}
	for _, station := range stations {
		cmd.Printf("%s\t%s\n", station.Identifier, station.Name)
	}

	return nil
}

func init() {
	RootCommand.AddCommand(stationCommand)
}
