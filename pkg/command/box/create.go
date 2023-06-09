package box

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/hckops/hckctl/pkg/command/common"
	"github.com/hckops/hckctl/pkg/command/config"
)

type boxCreateCmdOptions struct {
	configRef  *config.ConfigRef
	sourceFlag *common.SourceFlag
}

func NewBoxCreateCmd(configRef *config.ConfigRef) *cobra.Command {

	opts := &boxCreateCmdOptions{
		configRef: configRef,
	}

	command := &cobra.Command{
		Use:   "create [name]",
		Short: "create a detached box",
		RunE:  opts.run,
	}

	// --revision or --local
	opts.sourceFlag = common.AddTemplateSourceFlag(command)

	return command
}

func (opts *boxCreateCmdOptions) run(cmd *cobra.Command, args []string) error {
	fmt.Println("not implemented")
	return nil
}
