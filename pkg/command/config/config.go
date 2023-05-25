package config

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/hckops/hckctl/pkg/command/common"
)

// TODO add command to "set" a field with dot notation and "reset" all to default
type configCmdOptions struct {
	global *common.GlobalCmdOptions
}

func NewConfigCmd(globalOpts *common.GlobalCmdOptions) *cobra.Command {

	opts := configCmdOptions{
		global: globalOpts,
	}

	command := &cobra.Command{
		Use:   "config",
		Short: "print current configurations",
		RunE:  opts.run,
	}

	return command
}

func (opts *configCmdOptions) run(cmd *cobra.Command, args []string) error {
	fmt.Println("not implemented")
	return nil
}
