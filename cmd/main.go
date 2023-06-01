package main

import (
	"fmt"
	"os"

	"github.com/hckops/hckctl/pkg/command"
)

func main() {
	if err := command.NewRootCmd().Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
