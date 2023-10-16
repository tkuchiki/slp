package main

import (
	"log"

	"github.com/tkuchiki/slp/cmd/slp/cmd"
)

var version string

func main() {
	command := cmd.NewCommand(version)
	if err := command.Execute(); err != nil {
		log.Fatal(err)
	}
}
