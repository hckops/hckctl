package flag

import (
	"github.com/spf13/cobra"

	commonFlag "github.com/hckops/hckctl/internal/command/common/flag"
)

type CommandFlag struct {
	Inline bool
	Preset string
	Inputs []string
}

func addInlineFlag(command *cobra.Command, value *bool) string {
	const (
		flagName  = "inline"
		flagUsage = "use inline arguments"
	)
	command.Flags().BoolVarP(value, flagName, commonFlag.NoneFlagShortHand, false, flagUsage)
	return flagName
}

func addPresetFlag(command *cobra.Command, value *string) string {
	const (
		flagName  = "command"
		flagUsage = "use preset arguments"
	)
	command.Flags().StringVarP(value, flagName, commonFlag.NoneFlagShortHand, "", flagUsage)
	return flagName
}

func addInputsFlag(command *cobra.Command, value *[]string) string {
	const (
		flagName  = "input"
		flagUsage = "override command arguments"
	)
	command.Flags().StringArrayVarP(value, flagName, commonFlag.NoneFlagShortHand, []string{}, flagUsage)
	return flagName
}

func AddCommandFlag(command *cobra.Command) *CommandFlag {
	commandFlag := &CommandFlag{}
	inlineFlag := addInlineFlag(command, &commandFlag.Inline)
	presetFlag := addPresetFlag(command, &commandFlag.Preset)
	inputsFlag := addInputsFlag(command, &commandFlag.Inputs)
	command.MarkFlagsMutuallyExclusive(inlineFlag, presetFlag)
	command.MarkFlagsMutuallyExclusive(inlineFlag, inputsFlag)
	return commandFlag
}
