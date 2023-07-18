package box

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/hckops/hckctl/pkg/box"
	"github.com/hckops/hckctl/pkg/box/model"
	"github.com/hckops/hckctl/pkg/command/config"
)

// TODO output format yaml/json
type boxDescribeCmdOptions struct {
	configRef *config.ConfigRef
}

func NewBoxDescribeCmd(configRef *config.ConfigRef) *cobra.Command {

	opts := &boxDescribeCmdOptions{
		configRef: configRef,
	}

	command := &cobra.Command{
		Use:   "describe [name]",
		Short: "Describe a running box",
		RunE:  opts.run,
	}

	return command
}

func (opts *boxDescribeCmdOptions) run(cmd *cobra.Command, args []string) error {

	if len(args) == 1 {
		boxName := args[0]
		log.Debug().Msgf("describe box: boxName=%s", boxName)

		describeClient := func(client box.BoxClient, template *model.BoxV1) error {

			// TODO details
			if _, err := client.Describe(boxName); err != nil {
				// attempt next provider
				return err
			}

			// TODO
			fmt.Println(boxName)
			return nil
		}
		return attemptRunBoxClients(opts.configRef, boxName, describeClient)
	} else {
		cmd.HelpFunc()(cmd, args)
	}

	return nil
}
