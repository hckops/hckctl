package template

import (
	"github.com/spf13/cobra"

	"github.com/hckops/hckctl/pkg/command/common"
)

type templateCmdOptions struct {
	configRef *common.ConfigRef
	path      string
	revision  string
}

func NewTemplateCmd(configRef *common.ConfigRef) *cobra.Command {

	opts := &templateCmdOptions{
		configRef: configRef,
	}

	command := &cobra.Command{
		Use:   "template [name]",
		Short: "validate and print template",
		RunE:  opts.run,
	}

	const (
		pathFlag     = "path"
		revisionFlag = "revision"
	)
	command.PersistentFlags().StringVarP(&opts.path, pathFlag, "p", "", "local path")
	command.PersistentFlags().StringVarP(&opts.revision, revisionFlag, "r", common.DefaultMegalopolisBranch, "megalopolis version i.e. branch|tag|sha")
	command.MarkFlagsMutuallyExclusive(pathFlag, revisionFlag)

	command.AddCommand(NewTemplateShowCmd(opts)) // default
	command.AddCommand(NewTemplateValidateCmd(opts))

	listCommand := NewTemplateListCmd(opts)
	listCommand.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		// hide flag for this command
		command.PersistentFlags().MarkHidden(pathFlag)
		command.Parent().HelpFunc()(cmd, args)
	})
	command.AddCommand(listCommand)

	return command
}

func (opts *templateCmdOptions) run(cmd *cobra.Command, args []string) error {
	showOpts := &templateShowCmdOptions{
		template: opts,
		format:   yamlFlag,
	}
	return showOpts.run(cmd, args)
}
