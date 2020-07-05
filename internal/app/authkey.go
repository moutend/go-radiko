package app

import (
	"log"

	"github.com/moutend/go-radiko/pkg/radiko"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var authkeyCommand = &cobra.Command{
	Use:     "authkey",
	Aliases: []string{"a"},
	Short:   "print auth key",
	RunE:    authkeyCommandRunE,
}

func authkeyCommandRunE(cmd *cobra.Command, args []string) error {
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

	cmd.Print(session.AuthKey)

	return nil
}

func init() {
	RootCommand.AddCommand(authkeyCommand)
}
