package flag

import (
	"github.com/spf13/cobra"

	commonFlag "github.com/hckops/hckctl/internal/command/common/flag"
	boxModel "github.com/hckops/hckctl/pkg/box/model"
	commonModel "github.com/hckops/hckctl/pkg/common/model"
)

const (
	noExecFlagName   = "no-exec"
	noTunnelFlagName = "no-tunnel"
)

type TunnelFlag struct {
	NoExec   bool
	NoTunnel bool
}

func (f *TunnelFlag) ToConnectOptions(template *boxModel.BoxV1, name string, temporary bool) *boxModel.ConnectOptions {
	return &boxModel.ConnectOptions{
		Template:      template,
		StreamOpts:    commonModel.NewStdStreamOpts(true),
		Name:          name,
		DisableExec:   f.NoExec,
		DisableTunnel: f.NoTunnel,
		DeleteOnExit:  temporary,
	}
}

func addNoExecFlag(command *cobra.Command, value *bool) string {
	const (
		flagUsage = "tunnel all ports without spawning a shell"
	)
	command.Flags().BoolVarP(value, noExecFlagName, commonFlag.NoneFlagShortHand, false, flagUsage)
	return noExecFlagName
}

func addNoTunnelFlag(command *cobra.Command, value *bool) string {
	const (
		flagUsage = "spawn a shell without tunneling the ports"
	)
	command.Flags().BoolVarP(value, noTunnelFlagName, commonFlag.NoneFlagShortHand, false, flagUsage)
	return noTunnelFlagName
}

func AddTunnelFlag(command *cobra.Command) *TunnelFlag {
	tunnelFlag := &TunnelFlag{}
	noExecFlag := addNoExecFlag(command, &tunnelFlag.NoExec)
	noTunnelFlag := addNoTunnelFlag(command, &tunnelFlag.NoTunnel)
	command.MarkFlagsMutuallyExclusive(noExecFlag, noTunnelFlag)
	return tunnelFlag
}
