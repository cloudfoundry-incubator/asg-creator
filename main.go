package main

import (
	"log"

	"github.com/cloudfoundry-incubator/asg-creator/commands"
	"github.com/jessevdk/go-flags"
)

func main() {
	parser := flags.NewParser(&commands.ASGCreator, flags.None)

	_, err := parser.Parse()
	if err != nil {
		log.Fatalf("error: %s", err.Error())
	}
}
