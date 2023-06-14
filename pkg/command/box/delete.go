package box

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/hckops/hckctl/pkg/box"
	"github.com/hckops/hckctl/pkg/box/model"
	"github.com/hckops/hckctl/pkg/command/config"
)

type boxDeleteCmdOptions struct {
	configRef *config.ConfigRef
	prune     bool
}

func NewBoxDeleteCmd(configRef *config.ConfigRef) *cobra.Command {

	opts := &boxDeleteCmdOptions{
		configRef: configRef,
	}

	command := &cobra.Command{
		Use:   "delete [name]",
		Short: "delete running box",
		RunE:  opts.run,
	}

	return command
}

func (opts *boxDeleteCmdOptions) run(cmd *cobra.Command, args []string) error {

	if len(args) == 1 {
		boxName := args[0]
		log.Debug().Msgf("delete box: boxName=%s", boxName)

		deleteClient := func(client box.BoxClient, _ *model.BoxV1) error {
			return client.Delete(boxName)
		}
		return runRemoteBoxClient(opts.configRef, boxName, deleteClient)

	} else {
		cmd.HelpFunc()(cmd, args)
	}

	return nil
}
