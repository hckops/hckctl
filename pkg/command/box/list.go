package box

import (
	"fmt"

	"github.com/spf13/cobra"
)

type boxListCmdOptions struct {
	box *boxCmdOptions
}

func NewBoxListCmd(boxOpts *boxCmdOptions) *cobra.Command {

	opts := &boxListCmdOptions{
		box: boxOpts,
	}

	command := &cobra.Command{
		Use:   "list",
		Short: "TODO list",
		RunE:  opts.run,
	}

	return command
}

func (opts *boxListCmdOptions) run(cmd *cobra.Command, args []string) error {
	fmt.Println("not implemented")
	return nil
}
