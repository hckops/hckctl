package box

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/hckops/hckctl/pkg/box"
	"github.com/hckops/hckctl/pkg/box/model"
	"github.com/hckops/hckctl/pkg/command/config"
	"github.com/hckops/hckctl/pkg/util"
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

			details, err := client.Describe(boxName)
			if err != nil {
				// attempt next provider
				return err
			}

			if value, err := util.EncodeYaml(newBoxValue(details)); err != nil {
				return err
			} else {
				fmt.Print(value)
			}
			return nil
		}
		return attemptRunBoxClients(opts.configRef, boxName, describeClient)
	} else {
		cmd.HelpFunc()(cmd, args)
	}

	return nil
}

type BoxValue struct {
	Name          string
	Healthy       bool
	Provider      string
	Size          string
	LocalTemplate *model.LocalTemplateInfo `yaml:"local,omitempty"`
	GitTemplate   *model.GitTemplateInfo   `yaml:"git,omitempty"`
	Env           []model.BoxEnv           `yaml:",omitempty"`
	Ports         []model.BoxPort          `yaml:",omitempty"`
}

func newBoxValue(details *model.BoxDetails) *BoxValue {
	return &BoxValue{
		Name:          details.Info.Name,
		Healthy:       details.Info.Healthy,
		Provider:      details.Provider.String(),
		Size:          details.Size.String(),
		LocalTemplate: details.LocalTemplate,
		GitTemplate:   details.GitTemplate,
		Env:           details.Env,
		Ports:         details.Ports,
	}
}
