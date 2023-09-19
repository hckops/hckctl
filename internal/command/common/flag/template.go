package flag

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/hckops/hckctl/internal/command/common"
)

const (
	revisionFlagName = "revision"
	LocalFlagName    = "local"
	OfflineFlagName  = "offline"
)

type TemplateSourceFlag struct {
	Revision string
	Local    bool
	// TODO add remote
}

func AddTemplateRevisionFlag(command *cobra.Command, revision *string) string {
	const (
		flagShortName = "r"
		flagUsage     = "megalopolis version, one of branch|tag|sha"
	)

	command.Flags().StringVarP(revision, revisionFlagName, flagShortName, common.TemplateSourceRevision, flagUsage)
	// overrides default template val
	_ = viper.BindPFlag(fmt.Sprintf("template.%s", revisionFlagName), command.Flags().Lookup(revisionFlagName))

	return revisionFlagName
}

func AddTemplateLocalFlag(command *cobra.Command, local *bool) string {
	const (
		flagUsage = "use a local template"
	)
	command.Flags().BoolVarP(local, LocalFlagName, NoneFlagShortHand, false, flagUsage)
	return LocalFlagName
}

func AddTemplateOfflineFlag(command *cobra.Command, offline *bool) string {
	const (
		flagUsage = "ignore latest git templates"
	)
	command.Flags().BoolVarP(offline, OfflineFlagName, NoneFlagShortHand, false, flagUsage)
	return OfflineFlagName
}

func AddTemplateSourceFlag(command *cobra.Command) *TemplateSourceFlag {
	sourceFlag := &TemplateSourceFlag{}
	revisionFlag := AddTemplateRevisionFlag(command, &sourceFlag.Revision)
	localFlag := AddTemplateLocalFlag(command, &sourceFlag.Local)
	command.MarkFlagsMutuallyExclusive(revisionFlag, localFlag)
	return sourceFlag
}

func ValidateTemplateSourceFlag(provider *ProviderFlag, sourceFlag *TemplateSourceFlag) error {
	switch *provider {
	// clients can't decide revision or deploy custom templates
	case CloudProviderFlag:
		if sourceFlag.Revision != common.TemplateSourceRevision {
			return fmt.Errorf("flag not supported: provider=%s %s=%s",
				provider.String(), revisionFlagName, sourceFlag.Revision)
		}
		if sourceFlag.Local {
			return fmt.Errorf("flag not supported: provider=%s %s=%v",
				provider.String(), LocalFlagName, sourceFlag.Local)
		}
	}
	return nil
}
