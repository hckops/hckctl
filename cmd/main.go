package main

import (
	"fmt"
	"os"

	"github.com/hckops/hckctl/internal/command"
)

func main() {
	if err := command.NewRootCmd().Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
