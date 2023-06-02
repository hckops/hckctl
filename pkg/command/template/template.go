package template

import (
	"fmt"

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

	command.PersistentFlags().StringVarP(&opts.path, "path", "p", "", "load a local template")
	command.PersistentFlags().StringVarP(&opts.revision, "revision", "r", "main", "megalopolis version i.e. branch|tag|sha")
	command.MarkFlagsMutuallyExclusive("path", "revision")

	command.AddCommand(NewTemplateListCmd(opts))
	command.AddCommand(NewTemplateShowCmd(opts))
	command.AddCommand(NewTemplateValidateCmd(opts))

	return command
}

func (opts *templateCmdOptions) run(cmd *cobra.Command, args []string) error {
	// TODO alias of show
	fmt.Println("not implemented")
	return nil
}
