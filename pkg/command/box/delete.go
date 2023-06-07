package box

import (
	"fmt"

	"github.com/spf13/cobra"
)

type boxDeleteCmdOptions struct {
	box *boxCmdOptions
}

func NewBoxDeleteCmd(boxOpts *boxCmdOptions) *cobra.Command {

	opts := &boxDeleteCmdOptions{
		box: boxOpts,
	}

	command := &cobra.Command{
		Use:   "delete",
		Short: "TODO delete",
		RunE:  opts.run,
	}

	return command
}

func (opts *boxDeleteCmdOptions) run(cmd *cobra.Command, args []string) error {
	fmt.Println("not implemented")
	return nil
}
