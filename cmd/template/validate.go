package template

import (
	"fmt"

	"github.com/spf13/cobra"
)

// TODO validate multiple templates in a path and filter by regex
type templateValidateCmdOptions struct {
	template *templateCmdOptions
}

func NewTemplateValidateCmd(template *templateCmdOptions) *cobra.Command {

	opts := &templateValidateCmdOptions{
		template: template,
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
