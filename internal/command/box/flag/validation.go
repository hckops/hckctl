package flag

import (
	"fmt"

	"github.com/hckops/hckctl/internal/command/common"
	commonFlag "github.com/hckops/hckctl/internal/command/common/flag"
	"github.com/hckops/hckctl/pkg/box/model"
)

func ValidateSourceFlag(provider model.BoxProvider, sourceFlag *commonFlag.SourceFlag) error {
	switch provider {
	// clients can't decide revision or deploy custom templates
	case model.Cloud:
		if sourceFlag.Revision != common.TemplateSourceRevision {
			return fmt.Errorf("flag not supported: provider=%s %s=%s",
				model.Cloud.String(), commonFlag.RevisionFlagName, sourceFlag.Revision)
		}
		if sourceFlag.Local {
			return fmt.Errorf("flag not supported: provider=%s %s=%v",
				model.Cloud.String(), commonFlag.LocalFlagName, sourceFlag.Local)
		}
	}
	return nil
}

func ValidateTunnelFlag(provider model.BoxProvider, tunnelFlag *TunnelFlag) error {
	switch provider {
	// docker exposes automatically all ports
	case model.Docker:
		if tunnelFlag.TunnelOnly || tunnelFlag.NoTunnel {
			return fmt.Errorf("flag not supported: provider=%s %s=%v %s=%v",
				model.Docker.String(), tunnelOnlyFlagName, tunnelFlag.TunnelOnly, noTunnelFlagName, tunnelFlag.NoTunnel)
		}
	}
	return nil
}
