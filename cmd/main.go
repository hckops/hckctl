package main

import (
	"log"

	"github.com/hckops/hckctl/pkg/command"
)

func main() {
	// removes timestamps
	log.SetFlags(0)

	if err := command.NewRoodCmd().Execute(); err != nil {
		log.Fatal(err)
	}
}
