package box

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	boxFlag "github.com/hckops/hckctl/internal/command/box/flag"
	"github.com/hckops/hckctl/internal/command/common"
	commonFlag "github.com/hckops/hckctl/internal/command/common/flag"
	"github.com/hckops/hckctl/internal/command/config"
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
		loader := common.NewLoader()
		loader.Start("loading boxes")
		defer loader.Stop()

		// silently fail attempting all the providers
		for _, providerFlag := range boxFlag.BoxProviders() {
			if err := listByProvider(providerFlag, opts.configRef, loader); err != nil {
				log.Warn().Err(err).Msgf("ignoring error list boxes: providerFlag=%v", providerFlag)
			}
		}
	} else {
		cmd.HelpFunc()(cmd, args)
	}
	return nil
}

func listByProvider(providerFlag commonFlag.ProviderFlag, configRef *config.ConfigRef, loader *common.Loader) error {
	log.Debug().Msgf("list boxes: providerFlag=%s", providerFlag)

	provider, err := boxFlag.ToBoxProvider(providerFlag)
	if err != nil {
		return fmt.Errorf("%s provider error", providerFlag)
	}

	boxClient, err := newDefaultBoxClient(provider, configRef, loader)
	if err != nil {
		return err
	}

	boxes, err := boxClient.List()
	if err != nil {
		log.Warn().Err(err).Msgf("error listing boxes: provider=%v", boxClient.Provider())
		return fmt.Errorf("%s list error", boxClient.Provider())
	}

	loader.Stop()
	fmt.Println(fmt.Sprintf("# %s", boxClient.Provider()))
	for _, b := range boxes {
		if b.Healthy {
			fmt.Println(b.Name)
		} else {
			fmt.Println(fmt.Sprintf("%s (unhealthy)", b.Name))
		}
	}
	fmt.Println(fmt.Sprintf("total: %d", len(boxes)))
	return nil
}
