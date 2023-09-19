package flag

import (
	"github.com/spf13/cobra"

	"github.com/hckops/hckctl/internal/command/common"
)

func AddNetworkVpnFlag(command *cobra.Command, networkVpn *string) string {
	const (
		flagName  = "network-vpn"
		flagUsage = "connect to a vpn network"
	)
	command.Flags().StringVarP(networkVpn, flagName, NoneFlagShortHand, common.DefaultVpnName, flagUsage)
	return flagName
}
