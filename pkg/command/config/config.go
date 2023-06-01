package config

import (
	"fmt"
	"github.com/hckops/hckctl/pkg/util"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/thediveo/enumflag/v2"

	"github.com/spf13/cobra"

	"github.com/hckops/hckctl/pkg/command/common"
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

type configCmdOptions struct {
	common *common.CommonCmdOptions
}

func NewConfigCmd(commonOpts *common.CommonCmdOptions) *cobra.Command {

	opts := configCmdOptions{
		common: commonOpts,
	}

	command := &cobra.Command{
		Use:   "config",
		Short: "print current configurations",
		RunE:  opts.run,
	}

	return command
}

func (opts *configCmdOptions) run(cmd *cobra.Command, args []string) error {

	log.Info().Msgf("NewConfigCmd.run > %v", *opts.common)

	value, err := util.ToYaml(opts.common.Config)
	if err != nil {
		return errors.Wrap(err, "error encoding config")
	}
	fmt.Print(value)

	return nil
}
