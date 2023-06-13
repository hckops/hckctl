package box

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/hckops/hckctl/pkg/box"
	"github.com/hckops/hckctl/pkg/command/config"
	"github.com/hckops/hckctl/pkg/template/model"
)

type boxExecCmdOptions struct {
	configRef *config.ConfigRef
}

func NewBoxExecCmd(configRef *config.ConfigRef) *cobra.Command {

	opts := &boxExecCmdOptions{
		configRef: configRef,
	}

	command := &cobra.Command{
		Use:   "exec",
		Short: "exec in a box",
		RunE:  opts.run,
	}

	return command
}

func (opts *boxExecCmdOptions) run(cmd *cobra.Command, args []string) error {

	if len(args) == 1 {
		boxName := args[0]
		log.Debug().Msgf("exec remote box: boxName=%s", boxName)

		execClient := func(boxClient box.BoxClient, boxTemplate *model.BoxV1) error {
			return boxClient.Exec(boxName, boxTemplate.Shell)
		}
		return runBoxClient(opts.configRef, boxName, execClient)

	} else {
		cmd.HelpFunc()(cmd, args)
	}

	return nil
}
