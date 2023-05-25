package box

import (
	"fmt"

	"github.com/spf13/cobra"
)

type boxCreateCmdOptions struct {
	box      *boxCmdOptions
	path     string
	revision string
}

func NewBoxCreateCmd(boxOpts *boxCmdOptions) *cobra.Command {

	opts := &boxCreateCmdOptions{
		box: boxOpts,
	}

	command := &cobra.Command{
		Use:   "create",
		Short: "TODO create",
		RunE:  opts.run,
	}

	command.Flags().StringVarP(&opts.path, "path", "p", "", "load a local template")
	command.Flags().StringVarP(&opts.revision, "revision", "r", "main", "megalopolis version i.e. branch|tag|sha")
	command.MarkFlagsMutuallyExclusive("path", "revision")

	return command
}

func (opts *boxCreateCmdOptions) run(cmd *cobra.Command, args []string) error {
	fmt.Println("not implemented")
	return nil
}
