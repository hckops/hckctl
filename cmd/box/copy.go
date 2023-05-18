package box

import (
	"fmt"

	"github.com/spf13/cobra"
)

type boxCopyCmdOptions struct {
	box  *boxCmdOptions
	from string
	to   string
}

func NewBoxCopyCmd(boxOpts *boxCmdOptions) *cobra.Command {

	opts := &boxCopyCmdOptions{
		box: boxOpts,
	}

	command := &cobra.Command{
		Use:   "copy",
		Short: "TODO copy",
		RunE:  opts.run,
	}

	return command
}

func (opts *boxCopyCmdOptions) run(cmd *cobra.Command, args []string) error {
	fmt.Println("not implemented")
	return nil
}
