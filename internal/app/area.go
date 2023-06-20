package app

import (
	"log"

	"github.com/moutend/go-radiko/pkg/radiko"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var areaCommand = &cobra.Command{
	Use:     "area",
	Aliases: []string{"a"},
	Short:   "print your geolocation",
	RunE:    areaCommandRunE,
}

func areaCommandRunE(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()

	stationID := ""
	username := viper.GetString("RADIKO_USERNAME")
	password := viper.GetString("RADIKO_PASSWORD")

	client := radiko.New(stationID, username, password)

	if yes, _ := cmd.Flags().GetBool("debug"); yes {
		client.SetLogger(log.New(cmd.ErrOrStderr(), "debug: ", 0))
	}
	if err := client.GetAreaName(ctx); err != nil {
		return err
	}

	cmd.Println(client.AreaName)

	return nil
}

func init() {
	RootCommand.AddCommand(areaCommand)
}
