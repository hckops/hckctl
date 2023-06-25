package flag

import (
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/thediveo/enumflag/v2"

	"github.com/hckops/hckctl/pkg/box/model"
	"github.com/hckops/hckctl/pkg/command/common/flag"
)

// whitelist of supported box providers
func BoxProviders() []flag.ProviderFlag {
	return []flag.ProviderFlag{
		flag.DockerProviderFlag,
		flag.KubeProviderFlag,
		flag.CloudProviderFlag,
	}
}

func ToBoxProvider(p flag.ProviderFlag) (model.BoxProvider, error) {
	switch p {
	case flag.DockerProviderFlag:
		return model.Docker, nil
	case flag.KubeProviderFlag:
		return model.Kubernetes, nil
	case flag.CloudProviderFlag:
		return model.Cloud, nil
	default:
		return model.Docker, errors.New("invalid provider")
	}
}

func boxProviderIds() map[flag.ProviderFlag][]string {
	return flag.ProviderIds(BoxProviders())
}

func ValidateBoxProvider(configValue string, providerId *flag.ProviderFlag) (model.BoxProvider, error) {
	if configProvider, err := flag.ExistProvider(boxProviderIds(), configValue); err != nil {
		// must return a valid iota
		return model.Docker, errors.New("invalid config provider")
	} else if providerId.String() == flag.UnknownProvider {
		// default config
		return ToBoxProvider(configProvider)
	} else {
		// flag config
		return ToBoxProvider(*providerId)
	}
}

func AddBoxProviderFlag(command *cobra.Command) *flag.ProviderFlag {
	const (
		flagName = "provider"
	)
	var boxProviderFlag flag.ProviderFlag
	providerIds := boxProviderIds()
	providerValue := enumflag.NewWithoutDefault(&boxProviderFlag, flagName, providerIds, enumflag.EnumCaseInsensitive)
	providerUsageValues := strings.Join(flag.ProviderValues(boxProviderIds()), "|")
	providerUsage := fmt.Sprintf("switch box provider, one of %s", providerUsageValues)
	command.Flags().Var(providerValue, flagName, providerUsage)

	return &boxProviderFlag
}
