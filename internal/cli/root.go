package cli

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var RootCommand = &cobra.Command{
	Use:               "radiko",
	Short:             "radiko - command line radiko.jp client",
	PersistentPreRunE: rootPersistentPreRunE,
}

func rootPersistentPreRunE(cmd *cobra.Command, args []string) error {
	if path, _ := cmd.Flags().GetString("config"); path != "" {
		viper.SetConfigFile(path)
	}

	viper.AutomaticEnv()

	if path, _ := cmd.Flags().GetString("config"); path == "" {
		return nil
	}

	return viper.ReadInConfig()
}

func init() {
	RootCommand.PersistentFlags().BoolP("debug", "d", false, "enable debug output")
	RootCommand.PersistentFlags().StringP("config", "c", "", "path to configuration file")
}
