package flag

import (
	"fmt"

	"github.com/hckops/hckctl/pkg/box/model"
)

func ValidateTunnelFlag(provider model.BoxProvider, tunnelFlag *TunnelFlag) error {
	switch provider {
	// docker exposes automatically all ports
	case model.Docker:
		if tunnelFlag.NoExec || tunnelFlag.NoTunnel {
			return fmt.Errorf("flag not supported: provider=%s %s=%v %s=%v",
				model.Docker.String(), noExecFlagName, tunnelFlag.NoExec, noTunnelFlagName, tunnelFlag.NoTunnel)
		}
	}
	return nil
}
