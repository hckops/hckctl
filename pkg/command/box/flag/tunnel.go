package flag

import (
	"github.com/spf13/cobra"

	commonFlag "github.com/hckops/hckctl/pkg/command/common/flag"
)

type TunnelFlag struct {
	TunnelOnly bool
	NoTunnel   bool
}

func addTunnelOnlyFlag(command *cobra.Command, value *bool) string {
	const (
		flagName  = "tunnel-only"
		flagUsage = "TODO"
	)
	command.Flags().BoolVarP(value, flagName, commonFlag.NoneFlagShortHand, false, flagUsage)
	return flagName
}

func addNoTunnelFlag(command *cobra.Command, value *bool) string {
	const (
		flagName  = "no-tunnel"
		flagUsage = "TODO"
	)
	command.Flags().BoolVarP(value, flagName, commonFlag.NoneFlagShortHand, false, flagUsage)
	return flagName
}

func AddTunnelFlag(command *cobra.Command) *TunnelFlag {
	tunnelFlag := &TunnelFlag{}
	tunnelOnlyFlagName := addTunnelOnlyFlag(command, &tunnelFlag.TunnelOnly)
	noTunnelFlagName := addNoTunnelFlag(command, &tunnelFlag.NoTunnel)
	command.MarkFlagsMutuallyExclusive(tunnelOnlyFlagName, noTunnelFlagName)
	return tunnelFlag
}
