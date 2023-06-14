package box

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/hckops/hckctl/pkg/box"
	"github.com/hckops/hckctl/pkg/box/model"
	"github.com/hckops/hckctl/pkg/command/config"
)

type boxExecCmdOptions struct {
	configRef *config.ConfigRef
}

func NewBoxExecCmd(configRef *config.ConfigRef) *cobra.Command {

	opts := &boxExecCmdOptions{
		configRef: configRef,
	}

	command := &cobra.Command{
		Use:   "exec [name]",
		Short: "exec box",
		RunE:  opts.run,
	}

	return command
}

func (opts *boxExecCmdOptions) run(cmd *cobra.Command, args []string) error {

	if len(args) == 1 {
		boxName := args[0]
		log.Debug().Msgf("exec box: boxName=%s", boxName)

		execClient := func(client box.BoxClient, template *model.BoxV1) error {
			return client.Exec(boxName, template.Shell)
		}
		return runRemoteBoxClient(opts.configRef, boxName, execClient)

	} else {
		cmd.HelpFunc()(cmd, args)
	}

	return nil
}
