package flag

import (
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/thediveo/enumflag/v2"

	"github.com/hckops/hckctl/pkg/box/model"
	commonFlag "github.com/hckops/hckctl/pkg/command/common/flag"
)

// whitelist of supported box providers
func BoxProviders() []commonFlag.ProviderFlag {
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
		return model.Docker, errors.New("invalid provider")
	}
}

func boxProviderIds() map[commonFlag.ProviderFlag][]string {
	return commonFlag.ProviderIds(BoxProviders())
}

func ValidateBoxProvider(configValue string, providerId *commonFlag.ProviderFlag) (model.BoxProvider, error) {
	if configProvider, err := commonFlag.ExistProvider(boxProviderIds(), configValue); err != nil {
		// must return a valid iota
		return model.Docker, errors.New("invalid config provider")
	} else if providerId.String() == commonFlag.UnknownProvider {
		// default config
		return ToBoxProvider(configProvider)
	} else {
		// flag config
		return ToBoxProvider(*providerId)
	}
}

func AddBoxProviderFlag(command *cobra.Command) *commonFlag.ProviderFlag {
	const (
		flagName = "provider"
	)
	var boxProviderFlag commonFlag.ProviderFlag
	providerIds := boxProviderIds()
	providerValue := enumflag.NewWithoutDefault(&boxProviderFlag, flagName, providerIds, enumflag.EnumCaseInsensitive)
	providerUsageValues := strings.Join(commonFlag.ProviderValues(boxProviderIds()), "|")
	providerUsage := fmt.Sprintf("switch box provider, one of %s", providerUsageValues)
	command.Flags().Var(providerValue, flagName, providerUsage)

	return &boxProviderFlag
}
