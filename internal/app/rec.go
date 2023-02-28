package app

import (
	"fmt"
	"log"
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

	outputFile, _ := cmd.Flags().GetString("output")
	length, _ := cmd.Flags().GetDuration("length")

	if length <= 0 {
		return nil
	}

	dateFlag, _ := cmd.Flags().GetString("target")

	if dateFlag == "" {
		return nil
	}

	date, err := time.Parse(time.RFC3339, dateFlag)

	if err != nil {
		return fmt.Errorf("invalid date: %w", err)
	}

	ctx := cmd.Context()

	stationID := strings.ToUpper(args[0])
	username := viper.GetString("RADIKO_USERNAME")
	password := viper.GetString("RADIKO_PASSWORD")

	client := radiko.New(stationID, username, password)

	if yes, _ := cmd.Flags().GetBool("debug"); yes {
		client.SetLogger(log.New(cmd.ErrOrStderr(), "debug: ", 0))
	}
	if err := client.Rec(ctx, date, length, outputFile); err != nil {
		return err
	}

	return nil
}

func init() {
	RootCommand.AddCommand(recCommand)

	recCommand.PersistentFlags().StringP("target", "t", "", "target date with 'YYYYMMDDhhmm` layout (e.g. '201901021234')")
	recCommand.PersistentFlags().StringP("output", "o", "output.m4a", "output file name (default is 'output.m4a')")
	recCommand.PersistentFlags().DurationP("length", "l", 0, "recording length (e.g. '10s' is 10 seconds / '10m' is 10 minutes) ")
}
