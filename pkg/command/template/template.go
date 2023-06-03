package template

import (
	"fmt"
	"github.com/hckops/hckctl/pkg/old/util"
	"strings"

	"github.com/MakeNowJust/heredoc"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/thediveo/enumflag/v2"

	"github.com/hckops/hckctl/pkg/command/common"
)

type templateCmdOptions struct {
	local    bool
	revision string
	format   formatFlag
}

func NewTemplateCmd() *cobra.Command {

	opts := &templateCmdOptions{}

	command := &cobra.Command{
		Use:   "template [name]",
		Short: "validate and print template",
		Example: heredoc.Doc(`

			# prints remote alpine template
			hckctl template alpine

			# prints specific version (branch|tag|sha) of alpine template
			hckctl template alpine --revision main

			# prints template in json format (default yaml)
			hckctl template alpine --format json

			# validate and prints local template
			hckctl template boxes/official/alpine.yml --local
		`),
		RunE: opts.run,
	}

	const (
		formatFlagName   = "format"
		localFlagName    = "local"
		revisionFlagName = "revision"
	)

	// --format (enum)
	formatValue := enumflag.New(&opts.format, formatFlagName, formatIds, enumflag.EnumCaseInsensitive)
	formatUsage := fmt.Sprintf("output format, one of %s", strings.Join(formatValues(), "|"))
	command.Flags().Var(formatValue, formatFlagName, formatUsage)

	// --local
	command.Flags().BoolVarP(&opts.local, localFlagName, common.NoneFlagShortHand, false, "use local template")
	// --revision
	command.Flags().StringVarP(&opts.revision, revisionFlagName, "r", common.RevisionBranch, common.RevisionUsage)
	command.MarkFlagsMutuallyExclusive(localFlagName, revisionFlagName)

	command.AddCommand(NewTemplateListCmd())
	command.AddCommand(NewTemplateValidateCmd())

	return command
}

func (opts *templateCmdOptions) run(cmd *cobra.Command, args []string) error {
	format := opts.format.value()

	if opts.local {
		return printLocalTemplate(format, args[0])
	} else if len(args) == 1 {
		return printRemoteTemplate(format, args[0], opts.revision)
	} else {
		cmd.HelpFunc()(cmd, args)
	}
	return nil
}

func printLocalTemplate(format, path string) error {
	log.Debug().Msgf("print local template: format=%v path=%s", format, path)

	// TODO refactor and move all in RequestLocalTemplate
	localTemplate, err := util.ReadFile(path)
	if err != nil {
		return errors.Wrapf(err, "local template not found %s", localTemplate)
	}

	// TODO validation

	fmt.Print(localTemplate)
	return nil
}

func printRemoteTemplate(format, name, revision string) error {
	log.Debug().Msgf("print remote template: format=%v name=%s revision=%s", format, name, revision)
	return nil
}
