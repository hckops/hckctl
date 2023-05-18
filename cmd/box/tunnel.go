package box

import (
	"fmt"

	"github.com/spf13/cobra"
)

type boxTunnelCmdOptions struct {
	box *boxCmdOptions
}

func NewBoxTunnelCmd(boxOpts *boxCmdOptions) *cobra.Command {

	opts := &boxTunnelCmdOptions{
		box: boxOpts,
	}

	command := &cobra.Command{
		Use:   "tunnel",
		Short: "TODO tunnel",
		RunE:  opts.run,
	}

	return command
}

func (opts *boxTunnelCmdOptions) run(cmd *cobra.Command, args []string) error {
	fmt.Println("not implemented")
	return nil
}
