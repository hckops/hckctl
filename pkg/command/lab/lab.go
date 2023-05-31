package lab

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/hckops/hckctl/pkg/command/common"
)

// TODO command create, list, describe, delete

type labCmdOptions struct {
	global *common.GlobalCmdOptions
}

func NewLabCmd(globalOpts *common.GlobalCmdOptions) *cobra.Command {

	opts := labCmdOptions{
		global: globalOpts,
	}

	command := &cobra.Command{
		Use:   "lab [name]",
		Short: "manage a lab",
		RunE:  opts.run,
	}

	return command
}

func (opts *labCmdOptions) run(cmd *cobra.Command, args []string) error {
	fmt.Println("not implemented")
	return nil
}
