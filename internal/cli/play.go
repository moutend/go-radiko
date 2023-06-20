package cli

import (
	"log"
	"strings"

	"github.com/moutend/go-radiko/pkg/radiko"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var playCommand = &cobra.Command{
	Use:     "play",
	Aliases: []string{"p"},
	Short:   "play live stream",
	RunE:    playCommandRunE,
}

func playCommandRunE(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return nil
	}

	playbackVolume, _ := cmd.Flags().GetInt("volume")

	if playbackVolume < 0 || playbackVolume > 100 {
		playbackVolume = 100
	}

	ctx := cmd.Context()

	stationID := strings.ToUpper(args[0])
	username := viper.GetString("RADIKO_USERNAME")
	password := viper.GetString("RADIKO_PASSWORD")

	client := radiko.New(stationID, username, password)

	if yes, _ := cmd.Flags().GetBool("debug"); yes {
		client.SetLogger(log.New(cmd.ErrOrStderr(), "debug: ", 0))
	}
	if err := client.Play(ctx, playbackVolume); err != nil {
		return err
	}

	return nil
}

func init() {
	RootCommand.AddCommand(playCommand)

	playCommand.PersistentFlags().IntP("volume", "v", 100, "playback volume (min = 0, max = 100)")
}
