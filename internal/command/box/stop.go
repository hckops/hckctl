package box

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	boxFlag "github.com/hckops/hckctl/internal/command/box/flag"
	"github.com/hckops/hckctl/internal/command/common"
	commonFlag "github.com/hckops/hckctl/internal/command/common/flag"
	"github.com/hckops/hckctl/internal/command/config"
	"github.com/hckops/hckctl/pkg/box/model"
	"github.com/hckops/hckctl/pkg/template"
)

type boxStopCmdOptions struct {
	configRef *config.ConfigRef
	allFlag   bool // TODO refactor --providers="docker,kube" or --providers="all"
}

func NewBoxStopCmd(configRef *config.ConfigRef) *cobra.Command {

	opts := &boxStopCmdOptions{
		configRef: configRef,
	}

	command := &cobra.Command{
		Use:   "stop [name]",
		Short: "Stop one or more running boxes",
		Args:  cobra.MaximumNArgs(1),
		RunE:  opts.run,
	}

	const (
		allFlagName  = "all"
		allFlagUsage = "stop all boxes"
	)
	command.Flags().BoolVarP(&opts.allFlag, allFlagName, commonFlag.NoneFlagShortHand, false, allFlagUsage)

	return command
}

func (opts *boxStopCmdOptions) run(cmd *cobra.Command, args []string) error {

	if len(args) == 0 && opts.allFlag {
		loader := common.NewLoader()
		defer loader.Stop()
		loader.Start("stopping boxes")

		// silently fail attempting all the providers
		for _, providerFlag := range boxFlag.BoxProviders() {
			if err := stopByProvider(providerFlag, opts.configRef, loader); err != nil {
				log.Warn().Err(err).Msgf("ignoring error stopping boxes: providerFlag=%s", providerFlag)
			}
		}
		// cleanup cache directory
		if localPath, err := template.DeleteLocalCacheDir(opts.configRef.Config.Template.CacheDir); err != nil {
			return err
		} else {
			log.Debug().Msgf("delete local cache path: localPath=%s", localPath)
			return nil
		}

	} else if len(args) == 1 && !opts.allFlag {
		boxName := args[0]
		log.Debug().Msgf("stop box: boxName=%s", boxName)

		deleteClient := func(invokeOpts *invokeOptions, boxDetails *model.BoxDetails) error {

			if result, err := invokeOpts.client.Delete([]string{boxName}); err != nil {
				return err
			} else if len(result) == 0 {
				// attempt next provider
				return fmt.Errorf("box not found: boxName=%s", boxName)
			}
			invokeOpts.loader.Stop()
			fmt.Println(boxName)

			// cleanup cached template
			if boxDetails.TemplateInfo.IsCached() {
				log.Debug().Msgf("delete cached template: path=%s", invokeOpts.template.Path)
				return template.DeleteCachedTemplate(invokeOpts.template.Path)
			}
			return nil
		}
		return attemptRunBoxClients(opts.configRef, boxName, deleteClient)

	} else {
		cmd.HelpFunc()(cmd, args)
		return nil
	}
}

func stopByProvider(providerFlag commonFlag.ProviderFlag, configRef *config.ConfigRef, loader *common.Loader) error {
	log.Debug().Msgf("stop boxes: providerFlag=%s", providerFlag)

	provider, err := boxFlag.ToBoxProvider(providerFlag)
	if err != nil {
		return fmt.Errorf("%s provider error", providerFlag)
	}

	boxClient, err := newDefaultBoxClient(provider, configRef, loader)
	if err != nil {
		return err
	}

	names, err := boxClient.Delete([]string{})
	if err != nil {
		return fmt.Errorf("%s delete error", boxClient.Provider())
	}

	loader.Reload()
	fmt.Println(fmt.Sprintf("# %s", boxClient.Provider()))
	for _, name := range names {
		fmt.Println(name)
	}
	fmt.Println(fmt.Sprintf("total: %d", len(names)))
	return nil
}
