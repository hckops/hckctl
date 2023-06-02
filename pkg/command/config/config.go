package config

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/thediveo/enumflag/v2"

	"github.com/spf13/cobra"

	"github.com/hckops/hckctl/pkg/command/common"
	"github.com/hckops/hckctl/pkg/util"
)

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

// TODO add command to "set" a field with dot notation and "reset" all to default

type configCmdOptions struct{}

func NewConfigCmd(configRef *common.ConfigRef) *cobra.Command {

	opts := &configCmdOptions{}

	command := &cobra.Command{
		Use:   "config",
		Short: "print current configurations",
		RunE: func(cmd *cobra.Command, args []string) error {
			return opts.print(configRef.Config)
		},
	}

	return command
}

func (opts *configCmdOptions) print(configValue *common.Config) error {

	if value, err := util.ToYaml(configValue); err != nil {
		return errors.Wrap(err, "error encoding config")
	} else {
		fmt.Print(value)
	}

	return nil
}
