package config

import (
	"fmt"
	"github.com/thediveo/enumflag/v2"

	"github.com/spf13/cobra"

	"github.com/hckops/hckctl/pkg/command/common"
)

// TODO add command to "set" a field with dot notation and "reset" all to default
type configCmdOptions struct {
	global *common.GlobalCmdOptions
}

type ProviderFlag enumflag.Flag

const (
	DockerFlag ProviderFlag = iota
	KubernetesFlag
	CloudFlag
)

var ProviderIds = map[ProviderFlag][]string{
	DockerFlag:     {string(Docker)},
	KubernetesFlag: {string(Kubernetes)},
	CloudFlag:      {string(Cloud)},
}

func ProviderToId(provider ProviderFlag) string {
	return ProviderIds[provider][0]
}

func ProviderToFlag(value Provider) (ProviderFlag, error) {
	switch value {
	case Docker:
		return DockerFlag, nil
	case Kubernetes:
		return KubernetesFlag, nil
	case Cloud:
		return CloudFlag, nil
	default:
		return 999, fmt.Errorf("invalid provider")
	}
}

func NewConfigCmd(globalOpts *common.GlobalCmdOptions) *cobra.Command {

	opts := configCmdOptions{
		global: globalOpts,
	}

	command := &cobra.Command{
		Use:   "config",
		Short: "print current configurations",
		RunE:  opts.run,
	}

	return command
}

func (opts *configCmdOptions) run(cmd *cobra.Command, args []string) error {
	fmt.Println("not implemented")
	return nil
}
