package template

import (
	"fmt"
	"github.com/spf13/cobra"
)

// TODO format e.g. convert to json
type templateShowCmdOptions struct {
	template *templateCmdOptions
}

func NewTemplateShowCmd(templateOpts *templateCmdOptions) *cobra.Command {

	opts := &templateShowCmdOptions{
		template: templateOpts,
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
