package box

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	boxFlag "github.com/hckops/hckctl/pkg/command/box/flag"
	commonFlag "github.com/hckops/hckctl/pkg/command/common/flag"
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
		Short: "List all running boxes",
		RunE:  opts.run,
	}

	return command
}

func (opts *boxListCmdOptions) run(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		// silently fail attempting all the providers
		for _, providerFlag := range boxFlag.BoxProviders() {
			if err := listByProvider(providerFlag, opts.configRef); err != nil {
				log.Warn().Err(err).Msgf("ignoring error list boxes: providerFlag=%v", providerFlag)
			}
		}
	} else {
		cmd.HelpFunc()(cmd, args)
	}
	return nil
}

func listByProvider(providerFlag commonFlag.ProviderFlag, configRef *config.ConfigRef) error {
	log.Debug().Msgf("list boxes: providerFlag=%v", providerFlag)

	boxClient, err := newDefaultBoxClient(providerFlag, configRef)
	if err != nil {
		return err
	}

	boxes, err := boxClient.List()
	if err != nil {
		log.Warn().Err(err).Msgf("error listing boxes: provider=%v", boxClient.Provider())
		return fmt.Errorf("%v list error", boxClient.Provider())
	}

	fmt.Println(fmt.Sprintf("# %v", boxClient.Provider()))
	for _, b := range boxes {
		fmt.Println(b.Name)
	}
	fmt.Println(fmt.Sprintf("total: %v", len(boxes)))
	return nil
}
