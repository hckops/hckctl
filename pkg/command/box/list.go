package box

import (
	"fmt"
	"github.com/MakeNowJust/heredoc"
	"github.com/hckops/hckctl/pkg/box"
	"github.com/hckops/hckctl/pkg/client"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/hckops/hckctl/pkg/command/config"
)

type boxListCmdOptions struct {
	configRef *config.ConfigRef
	provider  string // TODO filter by provider, default all
}

func NewBoxListCmd(configRef *config.ConfigRef) *cobra.Command {

	opts := &boxListCmdOptions{
		configRef: configRef,
	}

	command := &cobra.Command{
		Use:   "list",
		Short: "list available templates",
		Example: heredoc.Doc(`

			# list all running boxes
			hckctl box list
		`),
		RunE: opts.run,
	}

	return command
}

func (opts *boxListCmdOptions) run(cmd *cobra.Command, args []string) error {

	boxClient, err := box.NewBoxClient(box.Docker)
	if err != nil {
		log.Warn().Err(err).Msg("error creating client")
		return errors.New("client error")
	}
	boxClient.Events().Subscribe(func(event client.Event) {
		log.Debug().Msg(event.String())
	})

	fmt.Println(fmt.Sprintf("# %s", box.Docker))
	dockerBoxes, err := boxClient.List()
	for _, dockerBox := range dockerBoxes {
		fmt.Println(dockerBox.Name)
	}
	return nil
}
