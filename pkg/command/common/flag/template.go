package flag

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/hckops/hckctl/pkg/command/common"
)

const (
	NoneFlagShortHand = ""
)

type SourceFlag struct {
	Revision string
	Local    bool
}

func AddRevisionFlag(command *cobra.Command, revision *string) string {
	const (
		flagName      = "revision"
		flagShortName = "r"
		flagUsage     = "megalopolis version, one of branch|tag|sha"
	)

	command.Flags().StringVarP(revision, flagName, flagShortName, common.TemplateSourceRevision, flagUsage)
	// overrides default template val
	_ = viper.BindPFlag(fmt.Sprintf("template.%s", flagName), command.Flags().Lookup(flagName))

	return flagName
}

func AddLocalFlag(command *cobra.Command, local *bool) string {
	const (
		flagName  = "local"
		flagUsage = "use a local template"
	)
	command.Flags().BoolVarP(local, flagName, NoneFlagShortHand, false, flagUsage)
	return flagName
}

func AddTemplateSourceFlag(command *cobra.Command) *SourceFlag {
	sourceFlag := &SourceFlag{}
	revisionFlagName := AddRevisionFlag(command, &sourceFlag.Revision)
	localFlagName := AddLocalFlag(command, &sourceFlag.Local)
	command.MarkFlagsMutuallyExclusive(revisionFlagName, localFlagName)
	return sourceFlag
}
