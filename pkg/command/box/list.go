package box

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/hckops/hckctl/pkg/command/common/flag"
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
	// silently fails attempting all the providers
	for _, providerFlag := range boxProviders() {
		if err := listByProvider(providerFlag, opts.configRef); err != nil {
			log.Warn().Err(err).Msgf("ignoring error list boxes: providerFlag=%v", providerFlag)
		}
	}
	return nil
}

func listByProvider(providerFlag flag.ProviderFlag, configRef *config.ConfigRef) error {
	log.Debug().Msgf("list boxes: providerFlag=%v", providerFlag)

	boxClient, err := newDefaultBoxClient(providerFlag, configRef)
	if err != nil {
		return err
	}

	fmt.Println(fmt.Sprintf("# %v", boxClient.Provider()))
	boxes, err := boxClient.List()
	if err != nil {
		log.Warn().Err(err).Msgf("error listing boxes: provider=%v", boxClient.Provider())
		return fmt.Errorf("%v list error", boxClient.Provider())
	}
	for _, b := range boxes {
		fmt.Println(b.Name)
	}
	fmt.Println(fmt.Sprintf("total: %v", len(boxes)))
	return nil
}
