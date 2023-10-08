package flag

import (
	"errors"

	"github.com/spf13/cobra"

	commonFlag "github.com/hckops/hckctl/internal/command/common/flag"
	"github.com/hckops/hckctl/pkg/task/model"
)

func taskProviders() []commonFlag.ProviderFlag {
	return []commonFlag.ProviderFlag{
		commonFlag.DockerProviderFlag,
		commonFlag.KubeProviderFlag,
	}
}

func toTaskProvider(p commonFlag.ProviderFlag) (model.TaskProvider, error) {
	switch p {
	case commonFlag.DockerProviderFlag:
		return model.Docker, nil
	case commonFlag.KubeProviderFlag:
		return model.Kubernetes, nil
	default:
		return commonFlag.UnknownProvider, errors.New("invalid provider")
	}
}

func taskProviderIds() map[commonFlag.ProviderFlag][]string {
	return commonFlag.ProviderIds(taskProviders())
}

func ValidateTaskProviderFlag(configValue string, providerId *commonFlag.ProviderFlag) (model.TaskProvider, error) {
	if configProvider, err := commonFlag.ExistProvider(taskProviderIds(), configValue); err != nil {
		return commonFlag.UnknownProvider, errors.New("invalid config provider")
	} else if providerId.String() == commonFlag.UnknownProvider {
		// default config
		return toTaskProvider(configProvider)
	} else {
		// flag
		return toTaskProvider(*providerId)
	}
}

func AddTaskProviderFlag(command *cobra.Command) *commonFlag.ProviderFlag {
	return commonFlag.AddProviderFlag(command, taskProviderIds())
}
