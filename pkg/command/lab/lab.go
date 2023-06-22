package lab

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/hckops/hckctl/pkg/command/config"
)

// TODO command create, list, describe, delete

type labCmdOptions struct {
	configRef *config.ConfigRef
}

func NewLabCmd(configRef *config.ConfigRef) *cobra.Command {

	opts := labCmdOptions{
		configRef: configRef,
	}

	command := &cobra.Command{
		Use:   "lab [name]",
		Short: "Manage labs",
		RunE:  opts.run,
	}

	return command
}

func (opts *labCmdOptions) run(cmd *cobra.Command, args []string) error {
	fmt.Println("not implemented")
	return nil
}
