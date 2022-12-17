package cmd

import (
	"log"
	"os"

	"github.com/hckops/hckctl/box"
	"github.com/spf13/cobra"
)

func NewBoxCmd() *cobra.Command {
	command := &cobra.Command{
		Use:   "box",
		Short: "TODO box",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				cmd.HelpFunc()(cmd, args)
				os.Exit(1)
			}
		},
	}

	openCmd := &cobra.Command{
		Use:   "open",
		Short: "TODO open",
		Run: func(cmd *cobra.Command, args []string) {
			if GlobalFlags().local {
				log.Println(box.OpenLocalBox())
			} else {
				log.Println(box.OpenBox())
			}
		},
	}

	command.AddCommand(openCmd)
	return command
}
