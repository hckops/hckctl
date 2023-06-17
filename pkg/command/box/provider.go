package box

import (
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/thediveo/enumflag/v2"

	"github.com/hckops/hckctl/pkg/box/model"
)

type providerFlag enumflag.Flag

const (
	unknownProviderFlag providerFlag = iota
	dockerProviderFlag
	kubeProviderFlag
	argoProviderFlag
	cloudProviderFlag
)

const (
	unknownProvider = "unknown"
)

var providerIds = map[providerFlag][]string{
	dockerProviderFlag: {"docker"},
	kubeProviderFlag:   {"kube", "k8s", "kubernetes"},
	argoProviderFlag:   {"argo", "argo-cd"},
	cloudProviderFlag:  {"cloud"},
}

func (p providerFlag) String() string {
	if p == unknownProviderFlag {
		return unknownProvider
	}
	return providerIds[p][0]
}

func (p providerFlag) toBoxProvider() (model.BoxProvider, error) {
	switch p {
	case dockerProviderFlag:
		return model.Docker, nil
	case kubeProviderFlag:
		return model.Kubernetes, nil
	case argoProviderFlag:
		return model.ArgoCd, nil
	case cloudProviderFlag:
		return model.Cloud, nil
	default:
		return model.Docker, errors.New("invalid flag provider")
	}
}

func providerValues() []string {
	var values []string
	for _, providerId := range providerIds {
		for _, provider := range providerId {
			values = append(values, provider)
		}
	}
	return values
}

func existProvider(value string) (model.BoxProvider, error) {
	for flag, providerId := range providerIds {
		for _, provider := range providerId {
			if value == provider {
				return flag.toBoxProvider()
			}
		}
	}
	return model.Docker, errors.New("invalid provider")
}

func validateProvider(configValue string, flagValue *providerFlag) (model.BoxProvider, error) {
	if configProvider, err := existProvider(configValue); err != nil {
		return configProvider, errors.New("invalid config provider")
	} else if flagValue.String() == unknownProvider {
		return configProvider, nil
	} else {
		return flagValue.toBoxProvider()
	}
}

func addBoxProviderFlag(command *cobra.Command) *providerFlag {
	const (
		flagName = "provider"
	)
	var enumFlagValue providerFlag
	providerValue := enumflag.NewWithoutDefault(&enumFlagValue, flagName, providerIds, enumflag.EnumCaseInsensitive)
	providerUsage := fmt.Sprintf("switch box provider, one of %s", strings.Join(providerValues(), "|"))
	command.Flags().Var(providerValue, flagName, providerUsage)

	return &enumFlagValue
}
