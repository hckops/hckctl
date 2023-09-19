package box

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	boxFlag "github.com/hckops/hckctl/internal/command/box/flag"
	commonFlag "github.com/hckops/hckctl/internal/command/common/flag"
	"github.com/hckops/hckctl/internal/command/config"
	"github.com/hckops/hckctl/pkg/box/model"
)

type boxOpenCmdOptions struct {
	configRef  *config.ConfigRef
	tunnelFlag *boxFlag.TunnelFlag
}

func NewBoxOpenCmd(configRef *config.ConfigRef) *cobra.Command {

	opts := &boxOpenCmdOptions{
		configRef: configRef,
	}

	command := &cobra.Command{
		Use:   "open [name]",
		Short: "Access and tunnel a running box",
		Args:  cobra.ExactArgs(1),
		RunE:  opts.run,
	}

	// --no-exec or --no-tunnel
	opts.tunnelFlag = boxFlag.AddTunnelFlag(command)

	return command
}

func (opts *boxOpenCmdOptions) run(cmd *cobra.Command, args []string) error {
	boxName := args[0]
	log.Debug().Msgf("open box: boxName=%s", boxName)

	connectClient := func(invokeOpts *invokeOptions, _ *model.BoxDetails) error {

		// log only and ignore invalid tunnel flags to avoid false positive during provider attempts
		if err := boxFlag.ValidateTunnelFlag(opts.tunnelFlag, invokeOpts.client.Provider()); err != nil {
			log.Warn().Err(err).Msgf("ignore validation %s", commonFlag.ErrorFlagNotSupported)
		}

		connectOpts := opts.tunnelFlag.ToConnectOptions(&invokeOpts.template.Value.Data, boxName, false)
		return invokeOpts.client.Connect(connectOpts)
	}
	return attemptRunBoxClients(opts.configRef, boxName, connectClient)
}
