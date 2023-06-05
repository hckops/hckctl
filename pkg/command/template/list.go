package template

import (
	"fmt"

	"github.com/spf13/cobra"
)

type templateListCmdOptions struct {
	kind   string // TODO filter comma separated list e.g. "box,lab"
	order  string // TODO sort output
	column string // TODO output only specific fields
}

func NewTemplateListCmd() *cobra.Command {

	opts := &templateListCmdOptions{}

	command := &cobra.Command{
		Use:   "list",
		Short: "list available templates",
		RunE:  opts.run,
	}

	return command
}

func (opts *templateListCmdOptions) run(cmd *cobra.Command, args []string) error {
	fmt.Println("not implemented")
	return nil
}
