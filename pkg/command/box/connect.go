package box

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	boxFlag "github.com/hckops/hckctl/pkg/command/box/flag"
	commonFlag "github.com/hckops/hckctl/pkg/command/common/flag"
	"github.com/hckops/hckctl/pkg/command/config"
)

type boxConnectCmdOptions struct {
	configRef  *config.ConfigRef
	tunnelFlag *boxFlag.TunnelFlag
}

func NewBoxConnectCmd(configRef *config.ConfigRef) *cobra.Command {

	opts := &boxConnectCmdOptions{
		configRef: configRef,
	}

	command := &cobra.Command{
		Use:   "connect [name]",
		Short: "Access and tunnel a running box",
		RunE:  opts.run,
	}

	// --tunnel-only or --no-tunnel
	opts.tunnelFlag = boxFlag.AddTunnelFlag(command)

	return command
}

func (opts *boxConnectCmdOptions) run(cmd *cobra.Command, args []string) error {

	if len(args) == 1 {
		boxName := args[0]
		log.Debug().Msgf("connect box: boxName=%s", boxName)

		execClient := func(invokeOpts *invokeOptions) error {

			// log only and ignore invalid tunnel flags to avoid false positive during provider attempts
			if err := boxFlag.ValidateTunnelFlag(invokeOpts.client.Provider(), opts.tunnelFlag); err != nil {
				log.Warn().Err(err).Msgf("ignore validation %s", commonFlag.ErrorFlagNotSupported)
			}

			return invokeOpts.client.Connect(&invokeOpts.template.Value.Data, opts.tunnelFlag.ToTunnelOptions(), boxName)
		}
		return attemptRunBoxClients(opts.configRef, boxName, execClient)
	} else {
		cmd.HelpFunc()(cmd, args)
	}

	return nil
}
