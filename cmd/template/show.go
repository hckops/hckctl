package template

import (
	"fmt"
	"github.com/spf13/cobra"
)

// TODO format e.g. convert to json
type templateShowCmdOptions struct {
	template *templateCmdOptions
}

func NewTemplateShowCmd(template *templateCmdOptions) *cobra.Command {

	opts := &templateShowCmdOptions{
		template: template,
	}

	command := &cobra.Command{
		Use:   "show",
		Short: "validate and print template",
		RunE:  opts.run,
	}

	return command
}

func (opts *templateShowCmdOptions) run(cmd *cobra.Command, args []string) error {
	fmt.Println("not implemented")
	return nil
}
