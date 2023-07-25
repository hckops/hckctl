package flag

import (
	"github.com/spf13/cobra"

	commonFlag "github.com/hckops/hckctl/internal/command/common/flag"
	"github.com/hckops/hckctl/pkg/box/model"
)

const (
	tunnelOnlyFlagName = "tunnel-only"
	noTunnelFlagName   = "no-tunnel"
)

type TunnelFlag struct {
	TunnelOnly bool
	NoTunnel   bool
}

func (f *TunnelFlag) ToTunnelOptions() *model.TunnelOptions {
	return &model.TunnelOptions{
		Streams:    model.NewDefaultStreams(true),
		TunnelOnly: f.TunnelOnly,
		NoTunnel:   f.NoTunnel,
	}
}

func addTunnelOnlyFlag(command *cobra.Command, value *bool) string {
	const (
		flagUsage = "port-forward all ports without spawning a shell"
	)
	command.Flags().BoolVarP(value, tunnelOnlyFlagName, commonFlag.NoneFlagShortHand, false, flagUsage)
	return tunnelOnlyFlagName
}

func addNoTunnelFlag(command *cobra.Command, value *bool) string {
	const (
		flagUsage = "spawn a shell without port-forwarding the ports"
	)
	command.Flags().BoolVarP(value, noTunnelFlagName, commonFlag.NoneFlagShortHand, false, flagUsage)
	return noTunnelFlagName
}

func AddTunnelFlag(command *cobra.Command) *TunnelFlag {
	tunnelFlag := &TunnelFlag{}
	tunnelOnlyFlag := addTunnelOnlyFlag(command, &tunnelFlag.TunnelOnly)
	noTunnelFlag := addNoTunnelFlag(command, &tunnelFlag.NoTunnel)
	command.MarkFlagsMutuallyExclusive(tunnelOnlyFlag, noTunnelFlag)
	return tunnelFlag
}
