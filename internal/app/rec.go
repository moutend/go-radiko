package app

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
	"time"

	"github.com/moutend/go-radiko/pkg/radiko"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var recCommand = &cobra.Command{
	Use:     "rec",
	Aliases: []string{"r"},
	Short:   "record live stream",
	RunE:    recCommandRunE,
}

func recCommandRunE(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return nil
	}

	output, _ := cmd.Flags().GetString("output")

	if output == "" {
		return fmt.Errorf("output file name must not be empty")
	}

	length, _ := cmd.Flags().GetDuration("length")

	if length == 0 {
		return fmt.Errorf("length must be greater than 0")
	}

	target, _ := cmd.Flags().GetString("target")

	if target == "" {
		return fmt.Errorf("target must not be empty")
	}

	date, err := time.Parse("200601021504", target)

	if err != nil {
		return fmt.Errorf("target date layout is invalid: %w", err)
	}

	stations, err := radiko.GetStations()

	if err != nil {
		return err
	}

	id := strings.ToUpper(args[0])

	matched := stations.Match(func(s radiko.Station) bool {
		return strings.ToUpper(s.ID) == id
	})

	if !matched {
		return fmt.Errorf("cannot find radio station: id=%q", id)
	}

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

	ffmpeg := exec.CommandContext(
		cmd.Context(),
		`ffmpeg`,
		`-headers`, `Referer: http://radiko.jp/`,
		`-headers`, `Pragma: no-cache`,
		`-headers`, fmt.Sprintf("X-Radiko-AuthToken: %s", session.AuthToken),
		`-i`, fmt.Sprintf(`https://radiko.jp/v2/api/ts/playlist.m3u8?station_id=%s&l=15&ft=%s&to=%s`, id, date.Format(`20060102150405`), date.Add(length).Format(`20060102150405`)),
		`-acodec`, `copy`,
		`-vn`,
		`-bsf:a`, `aac_adtstoasc`,
		`-y`, output,
	)

	cmd.Println("Recording past live stream (Ctrl-C to quit)")

	ffmpeg.Run()

	return nil
}

func init() {
	RootCommand.AddCommand(recCommand)

	recCommand.PersistentFlags().StringP("target", "t", "", "target date with 'YYYYMMDDhhmm` layout (e.g. '201901021234')")
	recCommand.PersistentFlags().StringP("output", "o", "output.m4a", "output file name (default is 'output.m4a')")
	recCommand.PersistentFlags().DurationP("length", "l", 0, "recording length (e.g. '10s' is 10 seconds / '10m' is 10 minutes) ")
}
