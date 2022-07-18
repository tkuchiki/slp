package main

import (
	"log"

	"github.com/tkuchiki/slp/cmd/slp/cmd"
)

var version string

func main() {
	if err := cmd.Execute(version); err != nil {
		log.Fatal(err)
	}
}
