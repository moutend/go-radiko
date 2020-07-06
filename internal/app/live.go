package app

import (
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os/exec"
	"strings"

	"github.com/moutend/go-radiko/pkg/radiko"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var liveCommand = &cobra.Command{
	Use:     "live",
	Aliases: []string{"l"},
	Short:   "play live stream",
	RunE:    liveCommandRunE,
}

func liveCommandRunE(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return nil
	}

	playbackVolume, _ := cmd.Flags().GetString("volume")

	if playbackVolume == "" {
		playbackVolume = "100"
	}

	p := make([]byte, 16, 16)

	if _, err := rand.Read(p); err != nil {
		return err
	}

	stations, err := radiko.GetStations()

	if err != nil {
		return err
	}

	station := strings.ToUpper(args[0])
	found := false

	for i, _ := range stations {
		if stations[i].Identifier == station {
			found = true

			break
		}
	}
	if !found {
		return fmt.Errorf("can't find radio station %q", station)
	}

	uuid := hex.EncodeToString(p)
	username := viper.GetString("RADIKO_USERNAME")
	password := viper.GetString("RADIKO_PASSWORD")
	session := radiko.NewSession(username, password)

	if yes, _ := cmd.Flags().GetBool("debug"); yes {
		session.SetLogger(log.New(cmd.ErrOrStderr(), "debug: ", 0))
	}
	if err := session.Login(); err != nil {
		return err
	}
	if err := session.Auth1(); err != nil {
		return err
	}
	if err := session.Auth2(); err != nil {
		return err
	}

	ffplay := exec.CommandContext(cmd.Context(), `ffplay`, `-volume`, playbackVolume, `-i`, `-`)
	ffmpeg := exec.CommandContext(
		cmd.Context(),
		`ffmpeg`,
		`-headers`, `Referer: http://radiko.jp/`,
		`-headers`, `Pragma: no-cache`,
		`-headers`, fmt.Sprintf("X-Radiko-AuthToken: %s", session.AuthToken),
		`-i`, fmt.Sprintf(`https://c-rpaa.smartstream.ne.jp/so/playlist.m3u8?station_id=%s&l=15&lsid=%s&type=c`, station, uuid),
		`-f`, `matroska`, `-`,
	)

	r, w := io.Pipe()
	ffmpeg.Stdout = w
	ffplay.Stdin = r

	cmd.Println("Playing live stream (Ctrl-C to quit)")

	ffmpeg.Start()
	ffplay.Start()

	ffmpeg.Wait()
	w.Close()

	ffplay.Wait()
	return nil
}

func init() {
	RootCommand.AddCommand(liveCommand)
	liveCommand.PersistentFlags().StringP("volume", "v", "", "path to configuration file")
}
