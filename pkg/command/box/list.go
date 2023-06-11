package box

import (
	"fmt"
	"strings"

	"github.com/MakeNowJust/heredoc"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/hckops/hckctl/pkg/box"
	"github.com/hckops/hckctl/pkg/client"
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
	for _, provider := range box.BoxProviders() {
		if err := listByProvider(provider); err != nil {
			return err
		}
	}
	return nil
}

func listByProvider(provider box.BoxProvider) error {
	boxClient, err := box.NewBoxClient(provider)
	if err != nil {
		log.Warn().Err(err).Msgf("error creating client: provider=%v", provider)
		return fmt.Errorf("%v client error", provider)
	}
	// provider not implemented
	if boxClient == nil {
		return nil
	}

	boxClient.Events().Subscribe(func(event client.Event) {
		log.Debug().Msg(event.String())
	})

	fmt.Println(fmt.Sprintf("# %v", provider))
	boxes, err := boxClient.List()
	if err != nil {
		log.Warn().Err(err).Msgf("error listing %v boxes", provider)
		return fmt.Errorf("%v list error", provider)
	}
	for _, b := range boxes {
		fmt.Println(strings.TrimPrefix(b.Name, "/"))
	}
	fmt.Println(fmt.Sprintf("total: %v", len(boxes)))
	return nil
}
