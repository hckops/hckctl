package template

import (
	"github.com/spf13/cobra"
)

func NewTemplateCmd() *cobra.Command {
	var path string

	command := &cobra.Command{
		Use:   "template [NAME]",
		Short: "Loads and validates a template",
		Run: func(cmd *cobra.Command, args []string) {
			if path != "" {
				RunTemplateLocalCmd(path)
			} else if len(args) == 1 {
				RunTemplateRemoteCmd(args[0])
			} else {
				cmd.HelpFunc()(cmd, args)
			}
		},
	}
	command.Flags().StringVarP(&path, "path", "p", "", "load the template from a local path")
	return command
}
