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

func ValidateNetworkVpnFlag(name string, networks map[string]model.VpnNetworkInfo) error {
	if strings.TrimSpace(name) == "" {
		return nil
	}
	if _, ok := networks[name]; ok {
		return nil
	}
	return fmt.Errorf("vpn network [%s] config not found", name)
}
