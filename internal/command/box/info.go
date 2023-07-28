package box

import (
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/hckops/hckctl/internal/command/config"
	"github.com/hckops/hckctl/pkg/box/model"
	"github.com/hckops/hckctl/pkg/util"
)

// TODO output format yaml/json
type boxInfoCmdOptions struct {
	configRef *config.ConfigRef
}

func NewBoxInfoCmd(configRef *config.ConfigRef) *cobra.Command {

	opts := &boxInfoCmdOptions{
		configRef: configRef,
	}

	command := &cobra.Command{
		Use:   "info [name]",
		Short: "Describe a running box",
		RunE:  opts.run,
	}

	return command
}

func (opts *boxInfoCmdOptions) run(cmd *cobra.Command, args []string) error {

	if len(args) == 1 {
		boxName := args[0]
		log.Debug().Msgf("info box: boxName=%s", boxName)

		describeClient := func(invokeOpts *invokeOptions, boxDetails *model.BoxDetails) error {
			if value, err := util.EncodeYaml(newBoxValue(boxDetails)); err != nil {
				return err
			} else {
				invokeOpts.loader.Stop()
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
	Created       string
	Healthy       bool
	Size          string
	Provider      ProviderValue
	CacheTemplate *model.CachedTemplateInfo `yaml:"cache,omitempty"`
	GitTemplate   *model.GitTemplateInfo    `yaml:"git,omitempty"`
	Env           []model.BoxEnv            `yaml:",omitempty"` // TODO []string KEY=VALUE
	Ports         []model.BoxPort           `yaml:",omitempty"` // TODO []string tty/remote -> local (if not none)
}
type ProviderValue struct {
	Name           string
	DockerProvider *model.DockerProviderInfo `yaml:"docker,omitempty"`
	KubeProvider   *model.KubeProviderInfo   `yaml:"kubernetes,omitempty"`
}

// TODO match port alias and local
// TODO filter env from template
func newBoxValue(details *model.BoxDetails) *BoxValue {
	return &BoxValue{
		Name:    details.Info.Name,
		Created: details.Created.Format(time.RFC3339),
		Healthy: details.Info.Healthy,
		Size:    details.Size.String(),
		Provider: ProviderValue{
			Name:           details.ProviderInfo.Provider.String(),
			DockerProvider: details.ProviderInfo.DockerProvider,
			KubeProvider:   details.ProviderInfo.KubeProvider,
		},
		CacheTemplate: details.TemplateInfo.CachedTemplate,
		GitTemplate:   details.TemplateInfo.GitTemplate,
		Env:           details.Env,
		Ports:         details.Ports,
	}
}
