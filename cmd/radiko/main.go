package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/moutend/go-radiko/pkg/radiko"
)

func main() {
	log.SetFlags(0)
	log.SetPrefix("error:")

	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	enableDebugOutput := flag.Bool("d", false, "enable debug output")

	flag.Parse()

	session := radiko.NewSession(os.Getenv("RADIKO_USERNAME"), os.Getenv("RADIKO_PASSWORD"))

	if enableDebugOutput != nil && *enableDebugOutput {
		session.SetLogger(log.New(os.Stdout, "debug: ", 0))
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

	fmt.Print(session.AuthToken)

	return nil
}
