package box

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/hckops/hckctl/cmd/common"
)

// TODO filePath vs repoUrl + revision
type boxCmdOptions struct {
	global *common.GlobalCmdOptions
}

func NewBoxCmd(global *common.GlobalCmdOptions) *cobra.Command {

	opts := boxCmdOptions{
		global: global,
	}

	command := &cobra.Command{
		Use:   "box [NAME]",
		Short: "attach and tunnel a box",
		RunE:  opts.run,
	}

	return command
}

func (opts *boxCmdOptions) run(cmd *cobra.Command, args []string) error {
	fmt.Println("not implemented")
	return nil
}
