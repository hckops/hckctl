package flag

import (
	"fmt"
	"sort"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/thediveo/enumflag/v2"
)

type ProviderFlag enumflag.Flag

const (
	UnknownProviderFlag ProviderFlag = iota
	DockerProviderFlag
	KubeProviderFlag
	CloudProviderFlag
)

const (
	UnknownProvider = "unknown"
)

var allProviderIds = map[ProviderFlag][]string{
	DockerProviderFlag: {"docker"},
	KubeProviderFlag:   {"kube"}, // "k8s", "kubernetes"
	CloudProviderFlag:  {"cloud"},
}

func (p ProviderFlag) String() string {
	if p == UnknownProviderFlag {
		return UnknownProvider
	}
	return allProviderIds[p][0]
}

// ProviderIds builds a subset of all the available providers, required for the enum flag
func ProviderIds(providerFlags []ProviderFlag) map[ProviderFlag][]string {
	values := make(map[ProviderFlag][]string)
	for _, providerId := range providerFlags {
		if labels, ok := allProviderIds[providerId]; ok {
			values[providerId] = labels
		}
	}
	return values
}

// ProviderValues builds a list of labels to concatenate, required for the flag usage
func ProviderValues(providerIds map[ProviderFlag][]string) []string {
	var values []string
	for _, providerValues := range providerIds {
		for _, provider := range providerValues {
			values = append(values, provider)
		}
	}
	sort.Strings(values)
	return values
}

// ExistProvider verify if the given string is a valid provider
func ExistProvider(providerIds map[ProviderFlag][]string, value string) (ProviderFlag, error) {
	for providerId, providerValues := range providerIds {
		for _, providerValue := range providerValues {
			if value == providerValue {
				return providerId, nil
			}
		}
	}
	return UnknownProviderFlag, errors.New("invalid provider")
}

func AddProviderFlag(command *cobra.Command, providerIds map[ProviderFlag][]string) *ProviderFlag {
	const (
		flagName = "provider"
	)
	var providerFlag ProviderFlag
	providerValue := enumflag.NewWithoutDefault(&providerFlag, flagName, providerIds, enumflag.EnumCaseInsensitive)
	providerUsageValues := strings.Join(ProviderValues(providerIds), "|")
	providerUsage := fmt.Sprintf("switch provider, one of %s", providerUsageValues)
	command.Flags().Var(providerValue, flagName, providerUsage)
	return &providerFlag
}
