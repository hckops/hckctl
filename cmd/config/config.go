package config

import (
	"fmt"
	"github.com/spf13/cobra"
)

func NewConfigCmd() *cobra.Command {
	command := &cobra.Command{
		Use:   "config",
		Short: "print current configurations",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("not implemented")
		},
	}
	return command
}
