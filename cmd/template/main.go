package template

import (
	"os"

	"github.com/spf13/cobra"
)

func NewTemplateCmd() *cobra.Command {
	command := &cobra.Command{
		Use:   "template",
		Short: "TODO template",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				// TODO override usage: "hckctl template NAME"
				cmd.HelpFunc()(cmd, args)
				os.Exit(1)
			} else {
				fetchBox(args[0])
			}
		},
	}
	return command
}
