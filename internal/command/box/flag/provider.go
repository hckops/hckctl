package flag

import (
	"errors"

	"github.com/spf13/cobra"

	commonFlag "github.com/hckops/hckctl/internal/command/common/flag"
	"github.com/hckops/hckctl/pkg/box/model"
)

func BoxProviders() []commonFlag.ProviderFlag {
	// whitelist of supported box providers
	return []commonFlag.ProviderFlag{
		commonFlag.DockerProviderFlag,
		commonFlag.KubeProviderFlag,
		commonFlag.CloudProviderFlag,
	}
}

func ToBoxProvider(p commonFlag.ProviderFlag) (model.BoxProvider, error) {
	switch p {
	case commonFlag.DockerProviderFlag:
		return model.Docker, nil
	case commonFlag.KubeProviderFlag:
		return model.Kubernetes, nil
	case commonFlag.CloudProviderFlag:
		return model.Cloud, nil
	default:
		return commonFlag.UnknownProvider, errors.New("invalid provider")
	}
}

func boxProviderIds() map[commonFlag.ProviderFlag][]string {
	return commonFlag.ProviderIds(BoxProviders())
}

func ValidateBoxProviderFlag(configValue string, providerId *commonFlag.ProviderFlag) (model.BoxProvider, error) {
	if configProvider, err := commonFlag.ExistProvider(boxProviderIds(), configValue); err != nil {
		return commonFlag.UnknownProvider, errors.New("invalid config provider")
	} else if providerId.String() == commonFlag.UnknownProvider {
		// default config
		return ToBoxProvider(configProvider)
	} else {
		// flag
		return ToBoxProvider(*providerId)
	}
}

func AddBoxProviderFlag(command *cobra.Command) *commonFlag.ProviderFlag {
	return commonFlag.AddProviderFlag(command, boxProviderIds())
}
