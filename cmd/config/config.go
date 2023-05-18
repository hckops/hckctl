package config

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/hckops/hckctl/cmd/common"
)

// TODO add commands to "set" a field with dot notation and "reset" all to default
type configCmdOptions struct {
	global *common.GlobalCmdOptions
}

func NewConfigCmd(global *common.GlobalCmdOptions) *cobra.Command {

	opts := configCmdOptions{
		global: global,
	}

	command := &cobra.Command{
		Use:   "config",
		Short: "validate and print current configurations",
		RunE:  opts.run,
	}

	command.AddCommand(NewConfigShowCmd(&opts))

	return command
}

func (opts *configCmdOptions) run(cmd *cobra.Command, args []string) error {
	// TODO alias of show
	fmt.Println("not implemented")
	return nil
}
