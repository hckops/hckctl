package template

import (
	"github.com/spf13/cobra"
)

func NewTemplateCmd() *cobra.Command {
	var revision string
	var path string

	command := &cobra.Command{
		Use:   "template [NAME]",
		Short: "load and validate a template",
		Run: func(cmd *cobra.Command, args []string) {
			if path != "" {
				RunTemplateLocalCmd(path)
			} else if len(args) == 1 {
				RunTemplateRemoteCmd(args[0], revision)
			} else {
				cmd.HelpFunc()(cmd, args)
			}
		},
	}
	command.Flags().StringVarP(&revision, "revision", "r", "main", "git source version i.e. branch|tag|sha")
	command.Flags().StringVarP(&path, "path", "p", "", "load a template from a local path")
	command.MarkFlagsMutuallyExclusive("revision", "path")
	return command
}
