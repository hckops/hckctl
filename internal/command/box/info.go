package box

import (
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/hckops/hckctl/internal/command/config"
	boxModel "github.com/hckops/hckctl/pkg/box/model"
	commonModel "github.com/hckops/hckctl/pkg/common/model"
	"github.com/hckops/hckctl/pkg/util"
)

// TODO add output format yaml/json
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
		Args:  cobra.ExactArgs(1),
		RunE:  opts.run,
	}

	return command
}

func (opts *boxInfoCmdOptions) run(cmd *cobra.Command, args []string) error {
	boxName := args[0]
	log.Debug().Msgf("info box: boxName=%s", boxName)

	describeClient := func(invokeOpts *invokeOptions, boxDetails *boxModel.BoxDetails) error {

		if value, err := util.EncodeYaml(newBoxValue(&invokeOpts.template.Value.Data, boxDetails)); err != nil {
			return err
		} else {
			invokeOpts.loader.Stop()
			fmt.Print(value)
		}
		return nil
	}
	return attemptRunBoxClients(opts.configRef, boxName, describeClient)
}

type BoxValue struct {
	Name          string
	Created       string
	Healthy       bool
	Size          string
	Provider      ProviderValue
	CacheTemplate *commonModel.CachedTemplateInfo `yaml:"cache,omitempty"`
	GitTemplate   *commonModel.GitTemplateInfo    `yaml:"git,omitempty"`
	Env           []string                        `yaml:",omitempty"`
	Ports         []string                        `yaml:",omitempty"`
}
type ProviderValue struct {
	Name           string
	DockerProvider *commonModel.DockerProviderInfo `yaml:"docker,omitempty"`
	KubeProvider   *commonModel.KubeProviderInfo   `yaml:"kubernetes,omitempty"`
}

func newBoxValue(template *boxModel.BoxV1, details *boxModel.BoxDetails) *BoxValue {

	var envs []string
	for _, e := range details.Env {
		if _, exists := template.EnvironmentVariables()[e.Key]; exists {
			envs = append(envs, fmt.Sprintf("%s=%s", e.Key, e.Value))
		}
	}
	// TODO not used: match runtime port alias and local port (bound) from template
	// TODO change return type of details.Ports to map[string]BoxPort
	var ports []string
	for _, p := range details.Ports {
		ports = append(ports, fmt.Sprintf("%s/%s -> %s", p.Alias, p.Remote, p.Local))
	}

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
		Env:           envs,
		Ports:         template.Network.Ports,
	}
}
