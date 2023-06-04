package lab

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/hckops/hckctl/pkg/command/common"
)

// TODO command create, list, describe, delete

type labCmdOptions struct {
	configRef *common.ConfigRef
}

func NewLabCmd(configRef *common.ConfigRef) *cobra.Command {

	opts := labCmdOptions{
		configRef: configRef,
	}

	command := &cobra.Command{
		Use:   "lab [name]",
		Short: "manage labs",
		RunE:  opts.run,
	}

	return command
}

func (opts *labCmdOptions) run(cmd *cobra.Command, args []string) error {
	fmt.Println("not implemented")
	return nil
}
