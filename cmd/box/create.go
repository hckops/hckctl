package box

import (
	"fmt"

	"github.com/spf13/cobra"
)

type boxCreateCmdOptions struct {
	box *boxCmdOptions
}

func NewBoxCreateCmd(boxOpts *boxCmdOptions) *cobra.Command {

	opts := &boxCopyCmdOptions{
		box: boxOpts,
	}

	command := &cobra.Command{
		Use:   "create",
		Short: "TODO create",
		RunE:  opts.run,
	}

	return command
}

func (opts *boxCreateCmdOptions) run(cmd *cobra.Command, args []string) error {
	fmt.Println("not implemented")
	return nil
}
