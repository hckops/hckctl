package template

import (
	"github.com/spf13/cobra"
)

func NewTemplateCmd() *cobra.Command {
	var path string

	command := &cobra.Command{
		Use:   "template [NAME]",
		Short: "load and validate a template",
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
	command.Flags().StringVarP(&path, "path", "p", "", "load a template from a local path")
	return command
}
