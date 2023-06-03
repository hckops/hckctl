package template

import (
	"fmt"

	"github.com/spf13/cobra"
)

// TODO add folder validation + regex filter
type templateValidateCmdOptions struct {
	template *templateCmdOptions
}

func NewTemplateValidateCmd(templateOpts *templateCmdOptions) *cobra.Command {

	opts := &templateValidateCmdOptions{
		template: templateOpts,
	}

	command := &cobra.Command{
		Use:   "validate",
		Short: "validate template",
		RunE:  opts.run,
	}

	return command
}

func (opts *templateValidateCmdOptions) run(cmd *cobra.Command, args []string) error {
	fmt.Println("not implemented")
	return nil
}
