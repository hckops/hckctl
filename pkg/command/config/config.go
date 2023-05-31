package config

import (
	"fmt"
	"github.com/hckops/hckctl/pkg/util"
	"github.com/pkg/errors"
	"github.com/thediveo/enumflag/v2"

	"github.com/spf13/cobra"

	"github.com/hckops/hckctl/pkg/command/common"
)

// TODO add command to "set" a field with dot notation and "reset" all to default
type configCmdOptions struct {
	global *common.GlobalCmdOptions
	config *common.ConfigV1
}

type ProviderFlag enumflag.Flag

const (
	DockerFlag ProviderFlag = iota
	KubernetesFlag
	CloudFlag
)

var ProviderIds = map[ProviderFlag][]string{
	DockerFlag:     {string(common.Docker)},
	KubernetesFlag: {string(common.Kubernetes)},
	CloudFlag:      {string(common.Cloud)},
}

func ProviderToId(provider ProviderFlag) string {
	return ProviderIds[provider][0]
}

func ProviderToFlag(value common.Provider) (ProviderFlag, error) {
	switch value {
	case common.Docker:
		return DockerFlag, nil
	case common.Kubernetes:
		return KubernetesFlag, nil
	case common.Cloud:
		return CloudFlag, nil
	default:
		return 999, fmt.Errorf("invalid provider")
	}
}

// TODO commands reset, edit

func NewConfigCmd(globalOpts *common.GlobalCmdOptions, config *common.ConfigV1) *cobra.Command {

	opts := configCmdOptions{
		global: globalOpts,
		config: config,
	}

	command := &cobra.Command{
		Use:   "config",
		Short: "print current configurations",
		RunE:  opts.run,
	}

	return command
}

func (opts *configCmdOptions) run(cmd *cobra.Command, args []string) error {
	value, err := util.ToYaml(opts.config)
	if err != nil {
		return errors.Wrap(err, "error encoding config")
	}
	fmt.Print(value)
	return nil
}
