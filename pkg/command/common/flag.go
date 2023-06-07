package common

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	NoneFlagShortHand = ""
)

func AddRevisionFlag(command *cobra.Command, revision *string) string {
	const (
		flagName      = "revision"
		flagShortName = "r"
		flagUsage     = "megalopolis version, one of branch|tag|sha"
	)

	command.Flags().StringVarP(revision, flagName, flagShortName, TemplateSourceRevision, flagUsage)
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
