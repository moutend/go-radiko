package app

import (
	"encoding/json"

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

	data, err := json.MarshalIndent(stations, "", "  ")

	if err != nil {
		return err
	}

	cmd.Printf("%s\n", data)

	return nil
}

func init() {
	RootCommand.AddCommand(stationCommand)
}
