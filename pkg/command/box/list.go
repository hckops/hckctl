package box

import (
	"fmt"
	"github.com/hckops/hckctl/pkg/command/common/flag"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

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
		Short: "List running boxes",
		RunE:  opts.run,
	}

	return command
}

func (opts *boxListCmdOptions) run(cmd *cobra.Command, args []string) error {
	for _, providerFlag := range boxProviders() {
		if err := listByProvider(providerFlag, opts.configRef); err != nil {
			return err
		}
	}
	return nil
}

func listByProvider(providerFlag flag.ProviderFlag, configRef *config.ConfigRef) error {
	log.Debug().Msgf("list boxes: providerFlag=%v", providerFlag)

	provider, err := toBoxProvider(providerFlag)
	if err != nil {
		return err
	}
	boxClient, err := newDefaultBoxClient(provider, configRef)
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
