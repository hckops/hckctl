package box

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/hckops/hckctl/pkg/box"
	"github.com/hckops/hckctl/pkg/box/model"
	"github.com/hckops/hckctl/pkg/command/common/flag"
	"github.com/hckops/hckctl/pkg/command/config"
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
	command.Flags().BoolVarP(&opts.all, allFlagName, flag.NoneFlagShortHand, false, allFlagUsage)

	return command
}

func (opts *boxDeleteCmdOptions) run(cmd *cobra.Command, args []string) error {

	if len(args) == 0 && opts.all {
		// silently fail attempting all the providers
		for _, providerFlag := range boxProviders() {
			if err := deleteByProvider(providerFlag, opts.configRef); err != nil {
				log.Warn().Err(err).Msgf("ignoring error delete boxes: providerFlag=%v", providerFlag)
			}
		}
		return nil

	} else if len(args) == 1 && !opts.all {
		boxName := args[0]
		log.Debug().Msgf("delete box: boxName=%s", boxName)

		deleteClient := func(client box.BoxClient, _ *model.BoxV1) error {
			if err := client.Delete(boxName); err != nil {
				return err
			}
			fmt.Println(boxName)
			return nil
		}
		return attemptRunBoxClients(opts.configRef, boxName, deleteClient)

	} else {
		cmd.HelpFunc()(cmd, args)
	}

	return nil
}

func deleteByProvider(providerFlag flag.ProviderFlag, configRef *config.ConfigRef) error {
	log.Debug().Msgf("delete boxes: providerFlag=%v", providerFlag)

	boxClient, err := newDefaultBoxClient(providerFlag, configRef)
	if err != nil {
		return err
	}

	boxes, err := boxClient.DeleteAll()
	if err != nil {
		log.Warn().Err(err).Msgf("error deleting boxes: provider=%v", boxClient.Provider())
		return fmt.Errorf("%v delete error", boxClient.Provider())
	}

	fmt.Println(fmt.Sprintf("# %v", boxClient.Provider()))
	for _, b := range boxes {
		fmt.Println(b.Name)
	}
	fmt.Println(fmt.Sprintf("total: %v", len(boxes)))
	return nil
}
