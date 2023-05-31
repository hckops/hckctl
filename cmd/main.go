package main

import (
	"log"

	"github.com/hckops/hckctl/pkg/command"
)

func main() {
	// removes default timestamps
	log.SetFlags(0)

	if err := command.NewRootCmd().Execute(); err != nil {
		log.Fatal(err)
	}
}
