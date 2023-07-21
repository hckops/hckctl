package flag

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/hckops/hckctl/pkg/command/common"
)

const (
	RevisionFlagName = "revision"
	LocalFlagName    = "local"
	OfflineFlagName  = "offline"
)

// TODO add remote

type SourceFlag struct {
	Revision string
	Local    bool
}

func AddRevisionFlag(command *cobra.Command, revision *string) string {
	const (
		flagShortName = "r"
		flagUsage     = "megalopolis version, one of branch|tag|sha"
	)

	command.Flags().StringVarP(revision, RevisionFlagName, flagShortName, common.TemplateSourceRevision, flagUsage)
	// overrides default template val
	_ = viper.BindPFlag(fmt.Sprintf("template.%s", RevisionFlagName), command.Flags().Lookup(RevisionFlagName))

	return RevisionFlagName
}

func AddLocalFlag(command *cobra.Command, local *bool) string {
	const (
		flagUsage = "use a local template"
	)
	command.Flags().BoolVarP(local, LocalFlagName, NoneFlagShortHand, false, flagUsage)
	return LocalFlagName
}

func AddTemplateSourceFlag(command *cobra.Command) *SourceFlag {
	sourceFlag := &SourceFlag{}
	revisionFlag := AddRevisionFlag(command, &sourceFlag.Revision)
	localFlag := AddLocalFlag(command, &sourceFlag.Local)
	command.MarkFlagsMutuallyExclusive(revisionFlag, localFlag)
	return sourceFlag
}

func AddOfflineFlag(command *cobra.Command, offline *bool) string {
	const (
		flagUsage = "ignore latest git templates"
	)
	command.Flags().BoolVarP(offline, OfflineFlagName, NoneFlagShortHand, false, flagUsage)
	return OfflineFlagName
}
