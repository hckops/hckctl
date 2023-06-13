package box

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/hckops/hckctl/pkg/box"
	"github.com/hckops/hckctl/pkg/command/config"
	"github.com/hckops/hckctl/pkg/template/model"
)

type boxDeleteCmdOptions struct {
	configRef *config.ConfigRef
}

func NewBoxDeleteCmd(configRef *config.ConfigRef) *cobra.Command {

	opts := &boxDeleteCmdOptions{
		configRef: configRef,
	}

	command := &cobra.Command{
		Use:   "delete [name]",
		Short: "delete running boxes",
		RunE:  opts.run,
	}

	return command
}

func (opts *boxDeleteCmdOptions) run(cmd *cobra.Command, args []string) error {

	if len(args) == 1 {
		boxName := args[0]
		log.Debug().Msgf("delete remote box: boxName=%s", boxName)

		execClient := func(boxClient box.BoxClient, boxTemplate *model.BoxV1) error {
			return boxClient.Delete(boxName)
		}
		return runBoxClient(opts.configRef, boxName, execClient)

	} else {
		cmd.HelpFunc()(cmd, args)
	}

	return nil
}
