package box

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/hckops/hckctl/pkg/box/model"
	boxFlag "github.com/hckops/hckctl/pkg/command/box/flag"
	commonFlag "github.com/hckops/hckctl/pkg/command/common/flag"
	"github.com/hckops/hckctl/pkg/command/config"
	"github.com/hckops/hckctl/pkg/template"
)

type boxDeleteCmdOptions struct {
	configRef *config.ConfigRef
	all       bool // TODO refactor --providers="docker,kube" or --providers="all"
}

func NewBoxDeleteCmd(configRef *config.ConfigRef) *cobra.Command {

	opts := &boxDeleteCmdOptions{
		configRef: configRef,
	}

	command := &cobra.Command{
		Use:   "delete [name]",
		Short: "Delete one or more running boxes",
		RunE:  opts.run,
	}

	const (
		allFlagName  = "all"
		allFlagUsage = "delete all boxes"
	)
	command.Flags().BoolVarP(&opts.all, allFlagName, commonFlag.NoneFlagShortHand, false, allFlagUsage)

	return command
}

func (opts *boxDeleteCmdOptions) run(cmd *cobra.Command, args []string) error {

	if len(args) == 0 && opts.all {
		// silently fail attempting all the providers
		for _, providerFlag := range boxFlag.BoxProviders() {
			if err := deleteByProvider(providerFlag, opts.configRef); err != nil {
				log.Warn().Err(err).Msgf("ignoring error delete boxes: providerFlag=%v", providerFlag)
			}
		}
		// cleanup cache directory
		if localPath, err := template.DeleteLocalCacheDir(opts.configRef.Config.Template.CacheDir); err != nil {
			return err
		} else {
			log.Debug().Msgf("delete local cache path: localPath=%s", localPath)
			return nil
		}

	} else if len(args) == 1 && !opts.all {
		boxName := args[0]
		log.Debug().Msgf("delete box: boxName=%s", boxName)

		deleteClient := func(invokeOpts *invokeOptions, boxDetails *model.BoxDetails) error {

			if result, err := invokeOpts.client.Delete([]string{boxName}); err != nil {
				return err
			} else if len(result) == 0 {
				// attempt next provider
				return fmt.Errorf("box not found: boxName=%s", boxName)
			}
			fmt.Println(boxName)

			if boxDetails.TemplateInfo.IsCached() {
				log.Debug().Msgf("delete cached template: path=%s", invokeOpts.template.Path)
				return template.DeleteCachedTemplate(invokeOpts.template.Path)
			}
			return nil
		}
		return attemptRunBoxClients(opts.configRef, boxName, deleteClient)

	} else {
		cmd.HelpFunc()(cmd, args)
	}

	return nil
}

func deleteByProvider(providerFlag commonFlag.ProviderFlag, configRef *config.ConfigRef) error {
	log.Debug().Msgf("delete boxes: providerFlag=%v", providerFlag)

	boxClient, err := newDefaultBoxClient(providerFlag, configRef)
	if err != nil {
		return err
	}

	names, err := boxClient.Delete([]string{})
	if err != nil {
		log.Warn().Err(err).Msgf("error deleting boxes: provider=%v", boxClient.Provider())
		return fmt.Errorf("%v delete error", boxClient.Provider())
	}

	fmt.Println(fmt.Sprintf("# %v", boxClient.Provider()))
	for _, name := range names {
		fmt.Println(name)
	}
	fmt.Println(fmt.Sprintf("total: %v", len(names)))
	return nil
}
