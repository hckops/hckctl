package template

import (
	"fmt"
	"github.com/hckops/hckctl/pkg/command/common"

	"github.com/spf13/cobra"
)

type templateValidateCmdOptions struct {
	kind string
}

func NewTemplateValidateCmd() *cobra.Command {

	opts := &templateValidateCmdOptions{}

	command := &cobra.Command{
		Use:   "validate [path]",
		Short: "validate template",
		RunE:  opts.run,
	}

	command.Flags().StringVarP(&opts.kind, "kind", common.NoneFlagShortHand, "", "expected template kind")

	return command
}

// TODO validate multiple templates in the given path (not only single file) + add regex filter
func (opts *templateValidateCmdOptions) run(cmd *cobra.Command, args []string) error {
	fmt.Println("not implemented")
	return nil
}
