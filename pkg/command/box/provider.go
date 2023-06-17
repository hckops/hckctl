package box

import (
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/thediveo/enumflag/v2"

	"github.com/hckops/hckctl/pkg/box/model"
	"github.com/hckops/hckctl/pkg/command/common/flag"
)

func boxProviderIds() map[flag.ProviderFlag][]string {
	// whitelist
	var allowedBoxProviderIds = []flag.ProviderFlag{
		flag.DockerProviderFlag,
		flag.KubeProviderFlag,
	}

	return flag.ProviderIds(allowedBoxProviderIds)
}

func toBoxProvider(p flag.ProviderFlag) (model.BoxProvider, error) {
	switch p {
	case flag.DockerProviderFlag:
		return model.Docker, nil
	case flag.KubeProviderFlag:
		return model.Kubernetes, nil
	default:
		return model.Docker, errors.New("invalid provider")
	}
}

func validateBoxProvider(configValue string, providerId *flag.ProviderFlag) (model.BoxProvider, error) {
	if configProvider, err := flag.ExistProvider(boxProviderIds(), configValue); err != nil {
		// must return a valid iota
		return model.Docker, errors.New("invalid config provider")
	} else if providerId.String() == flag.UnknownProvider {
		// default config
		return toBoxProvider(configProvider)
	} else {
		// flag config
		return toBoxProvider(*providerId)
	}
}

func addBoxProviderFlag(command *cobra.Command) *flag.ProviderFlag {
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
