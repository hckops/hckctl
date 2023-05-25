package template

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/hckops/hckctl/pkg/command/common"
)

// TODO order, columns, etc.
type templateListCmdOptions struct {
	template *templateCmdOptions
	kind     string // TODO box, lab
}

func NewTemplateListCmd(templateOpts *templateCmdOptions) *cobra.Command {

	opts := &templateListCmdOptions{
		template: templateOpts,
	}

	command := &cobra.Command{
		Use:   "list",
		Short: "list templates",
		RunE:  opts.run,
	}

	command.Flags().StringVarP(&opts.kind, "kind", common.NoneFlagShortHand, "", "filter by kind")

	return command
}

func (opts *templateListCmdOptions) run(cmd *cobra.Command, args []string) error {
	fmt.Println("not implemented")
	return nil
}
