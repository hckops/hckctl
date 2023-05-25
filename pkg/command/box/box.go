package box

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/hckops/hckctl/pkg/command/common"
)

type boxCmdOptions struct {
	global   *common.GlobalCmdOptions
	path     string
	revision string
}

func NewBoxCmd(globalOpts *common.GlobalCmdOptions) *cobra.Command {

	opts := boxCmdOptions{
		global: globalOpts,
	}

	command := &cobra.Command{
		Use:   "box [name]",
		Short: "attach and tunnel a box",
		RunE:  opts.run,
	}

	command.Flags().StringVarP(&opts.path, "path", "p", "", "load a local template")
	command.Flags().StringVarP(&opts.revision, "revision", "r", "main", "megalopolis version i.e. branch|tag|sha")
	command.MarkFlagsMutuallyExclusive("path", "revision")

	command.AddCommand(NewBoxCopyCmd(&opts))
	command.AddCommand(NewBoxCreateCmd(&opts))
	command.AddCommand(NewBoxDeleteCmd(&opts))
	command.AddCommand(NewBoxExecCmd(&opts))
	command.AddCommand(NewBoxListCmd(&opts))
	command.AddCommand(NewBoxOpenCmd(&opts))
	command.AddCommand(NewBoxTunnelCmd(&opts))

	return command
}

func (opts *boxCmdOptions) run(cmd *cobra.Command, args []string) error {
	fmt.Println("not implemented")
	return nil
}
