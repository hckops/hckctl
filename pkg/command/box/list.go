package box

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/hckops/hckctl/pkg/box/model"
	"github.com/hckops/hckctl/pkg/command/config"
)

type boxListCmdOptions struct {
	configRef *config.ConfigRef
	providers string // TODO filter by provider (comma separated list), default all
}

func NewBoxListCmd(configRef *config.ConfigRef) *cobra.Command {

	opts := &boxListCmdOptions{
		configRef: configRef,
	}

	command := &cobra.Command{
		Use:   "list",
		Short: "list running boxes",
		RunE:  opts.run,
	}

	return command
}

func (opts *boxListCmdOptions) run(cmd *cobra.Command, args []string) error {
	// TODO model.BoxProviders()
	for _, provider := range []model.BoxProvider{model.Docker} {
		if err := listByProvider(provider); err != nil {
			return err
		}
	}
	return nil
}

func listByProvider(provider model.BoxProvider) error {
	boxClient, err := newDefaultBoxClient(provider)
	if err != nil {
		return err
	}

	fmt.Println(fmt.Sprintf("# %v", provider))
	boxes, err := boxClient.List()
	if err != nil {
		log.Warn().Err(err).Msgf("error listing boxes: provider=%v", provider)
		return fmt.Errorf("%v list error", provider)
	}
	for _, b := range boxes {
		fmt.Println(b.Name)
	}
	fmt.Println(fmt.Sprintf("total: %v", len(boxes)))
	return nil
}
