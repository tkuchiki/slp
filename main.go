package main

import (
	"log"

	"github.com/tkuchiki/slp/cmd/slp/cmd"
	buildversion "github.com/tkuchiki/slp/internal/version"
)

var version string

func main() {
	command := cmd.NewCommand(buildversion.Resolve(version))
	if err := command.Execute(); err != nil {
		log.Fatal(err)
	}
}
