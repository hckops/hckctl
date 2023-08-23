package flag

import (
	"fmt"

	"github.com/hckops/hckctl/internal/command/common"
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
