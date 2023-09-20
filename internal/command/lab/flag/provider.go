package flag

import (
	"errors"

	"github.com/spf13/cobra"

	commonFlag "github.com/hckops/hckctl/internal/command/common/flag"
	"github.com/hckops/hckctl/pkg/lab/model"
)

func labProviders() []commonFlag.ProviderFlag {
	return []commonFlag.ProviderFlag{
		commonFlag.CloudProviderFlag,
	}
}

func toLabProvider(p commonFlag.ProviderFlag) (model.LabProvider, error) {
	switch p {
	case commonFlag.CloudProviderFlag:
		return model.Cloud, nil
	default:
		return commonFlag.UnknownProvider, errors.New("invalid provider")
	}
}

func labProviderIds() map[commonFlag.ProviderFlag][]string {
	return commonFlag.ProviderIds(labProviders())
}

func ValidateLabProviderFlag(configValue string, providerId *commonFlag.ProviderFlag) (model.LabProvider, error) {
	if configProvider, err := commonFlag.ExistProvider(labProviderIds(), configValue); err != nil {
		return commonFlag.UnknownProvider, errors.New("invalid config provider")
	} else if providerId.String() == commonFlag.UnknownProvider {
		// default config
		return toLabProvider(configProvider)
	} else {
		// flag
		return toLabProvider(*providerId)
	}
}

func AddLabProviderFlag(command *cobra.Command) *commonFlag.ProviderFlag {
	return commonFlag.AddProviderFlag(command, labProviderIds())
}
