package flag

import (
	"fmt"
	"strings"

	"github.com/hckops/hckctl/internal/command/common"
	"github.com/hckops/hckctl/pkg/common/model"
)

func ValidateSourceFlag(provider *ProviderFlag, sourceFlag *SourceFlag) error {
	switch *provider {
	// clients can't decide revision or deploy custom templates
	case CloudProviderFlag:
		if sourceFlag.Revision != common.TemplateSourceRevision {
			return fmt.Errorf("flag not supported: provider=%s %s=%s",
				provider.String(), RevisionFlagName, sourceFlag.Revision)
		}
		if sourceFlag.Local {
			return fmt.Errorf("flag not supported: provider=%s %s=%v",
				provider.String(), LocalFlagName, sourceFlag.Local)
		}
	}
	return nil
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
