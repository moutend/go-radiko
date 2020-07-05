package main

import (
	"os"

	"github.com/moutend/go-radiko/internal/app"
)

func main() {
	app.RootCommand.SetOutput(os.Stdout)

	if err := app.RootCommand.Execute(); err != nil {
		os.Exit(-1)
	}
}
