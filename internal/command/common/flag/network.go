package flag

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/hckops/hckctl/pkg/common/model"
)

func AddNetworkVpnFlag(command *cobra.Command, networkVpn *string) string {
	const (
		flagName  = "network-vpn"
		flagUsage = "connect to a vpn network"
	)
	command.Flags().StringVarP(networkVpn, flagName, NoneFlagShortHand, "", flagUsage)
	return flagName
}

func ValidateNetworkVpnFlag(name string, networks map[string]model.VpnNetworkInfo) (*model.VpnNetworkInfo, error) {
	if strings.TrimSpace(name) == "" {
		return nil, nil
	}
	if vpnNetworkInfo, ok := networks[name]; ok {
		return &vpnNetworkInfo, nil
	}
	return nil, fmt.Errorf("vpn network [%s] config not found", name)
}
