package box

import (
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/hckops/hckctl/pkg/box"
	"github.com/hckops/hckctl/pkg/client"
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
		Use:   "exec",
		Short: "exec in a box",
		RunE:  opts.run,
	}

	return command
}

func (opts *boxExecCmdOptions) run(cmd *cobra.Command, args []string) error {

	boxClient, err := box.NewBoxClient(box.Docker)
	if err != nil {
		log.Warn().Err(err).Msg("error creating client")
		return errors.New("client error")
	}
	boxClient.Events().Subscribe(func(event client.Event) {
		log.Debug().Msg(event.String())
	})

	return nil
}
