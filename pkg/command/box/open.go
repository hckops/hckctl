package box

import (
	"fmt"

	"github.com/spf13/cobra"
)

type boxOpenCmdOptions struct {
	box *boxCmdOptions
}

func NewBoxOpenCmd(boxOpts *boxCmdOptions) *cobra.Command {

	opts := &boxOpenCmdOptions{
		box: boxOpts,
	}

	command := &cobra.Command{
		Use:   "open",
		Short: "TODO open",
		RunE:  opts.run,
	}

	return command
}

func (opts *boxOpenCmdOptions) run(cmd *cobra.Command, args []string) error {
	fmt.Println("not implemented")
	return nil
}
