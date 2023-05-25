package box

import (
	"fmt"

	"github.com/spf13/cobra"
)

type boxExecCmdOptions struct {
	box *boxCmdOptions
}

func NewBoxExecCmd(boxOpts *boxCmdOptions) *cobra.Command {

	opts := &boxExecCmdOptions{
		box: boxOpts,
	}

	command := &cobra.Command{
		Use:   "exec",
		Short: "TODO exec",
		RunE:  opts.run,
	}

	return command
}

func (opts *boxExecCmdOptions) run(cmd *cobra.Command, args []string) error {
	fmt.Println("not implemented")
	return nil
}
