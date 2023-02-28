package app

import (
	"encoding/json"
	"log"

	"github.com/moutend/go-radiko/pkg/radiko"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var stationCommand = &cobra.Command{
	Use:     "station",
	Aliases: []string{"s"},
	Short:   "list all available stations",
	RunE:    stationCommandRunE,
}

func stationCommandRunE(cmd *cobra.Command, args []string) error {
	username := viper.GetString("RADIKO_USERNAME")
	password := viper.GetString("RADIKO_PASSWORD")

	client := radiko.New("", username, password)

	if yes, _ := cmd.Flags().GetBool("debug"); yes {
		client.SetLogger(log.New(cmd.ErrOrStderr(), "debug: ", 0))
	}
	if err := client.GetAllStations(cmd.Context()); err != nil {
		return err
	}

	data, err := json.MarshalIndent(client.AllStations, "", "  ")

	if err != nil {
		return err
	}

	cmd.Printf("%s\n", data)

	return nil
}

func init() {
	RootCommand.AddCommand(stationCommand)
}
