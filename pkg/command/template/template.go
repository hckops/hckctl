package template

import (
	"fmt"
	"strings"

	"github.com/MakeNowJust/heredoc"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/thediveo/enumflag/v2"

	"github.com/hckops/hckctl/pkg/command/common"
	"github.com/hckops/hckctl/pkg/command/config"
	"github.com/hckops/hckctl/pkg/template/source"
)

type templateCmdOptions struct {
	configRef *config.ConfigRef
	format    formatFlag
	local     bool
	revision  string
}

func NewTemplateCmd(configRef *config.ConfigRef) *cobra.Command {

	opts := &templateCmdOptions{
		configRef: configRef,
	}

	command := &cobra.Command{
		Use:   "template [name]",
		Short: "validate and print templates",
		Example: heredoc.Doc(`

			# prints remote template (supports multiple formats)
			hckctl template alpine
			hckctl template official/parrot
			hckctl template boxes/vulnerable/dvwa
			hckctl template labs/official/htb-alpine.yml

			# prints specific version (branch|tag|sha)
			hckctl template alpine --revision main

			# prints template in json format (default yaml)
			hckctl template alpine --format json

			# validates and prints local template
			hckctl template ../megalopolis/boxes/official/alpine.yml --local
		`),
		RunE: opts.run,
	}

	const (
		formatFlagName = "format"
		localFlagName  = "local"
	)

	// --format (enum)
	formatValue := enumflag.New(&opts.format, formatFlagName, formatIds, enumflag.EnumCaseInsensitive)
	formatUsage := fmt.Sprintf("output format, one of %s", strings.Join(formatValues(), "|"))
	command.Flags().Var(formatValue, formatFlagName, formatUsage)

	// --local
	command.Flags().BoolVarP(&opts.local, localFlagName, common.NoneFlagShortHand, false, "use local template")
	// --revision
	revisionFlagName := common.AddRevisionFlag(command, &opts.revision)
	command.MarkFlagsMutuallyExclusive(localFlagName, revisionFlagName)

	command.AddCommand(NewTemplateListCmd(configRef))
	command.AddCommand(NewTemplateValidateCmd())

	return command
}

func (opts *templateCmdOptions) run(cmd *cobra.Command, args []string) error {
	format := opts.format.value()

	if len(args) == 1 && opts.local {
		path := args[0]
		log.Debug().Msgf("print local template: %s", path)

		return printTemplate(source.NewLocalSource(path), format)

	} else if len(args) == 1 {
		name := args[0]
		revisionOpts := &source.RevisionOpts{
			SourceCacheDir: opts.configRef.Config.Template.CacheDir,
			SourceUrl:      common.TemplateSourceUrl,
			SourceRevision: common.TemplateSourceRevision,
			Revision:       opts.revision,
		}
		log.Debug().Msgf("print remote template: %s", name)

		return printTemplate(source.NewRemoteSource(revisionOpts, name), format)

	} else {
		cmd.HelpFunc()(cmd, args)
	}
	return nil
}

func printTemplate(src source.TemplateSource, format string) error {

	templateValue, err := src.Read()
	if err != nil {
		log.Warn().Err(err).Msg("error printing template")
		return errors.New("invalid template")
	}

	if format == source.JsonFormat.String() {
		if jsonTemplateValue, err := templateValue.ToJson(); err != nil {
			return errors.New("format error")
		} else {
			log.Debug().Msgf("print template: kind=%s format=%s\n%s", jsonTemplateValue.Kind.String(), format, jsonTemplateValue.Data)
			// add newline
			fmt.Println(jsonTemplateValue.Data)
		}
	} else {
		log.Debug().Msgf("print template: kind=%s format=%s\n%s", templateValue.Kind.String(), format, templateValue.Data)
		fmt.Print(templateValue.Data)
	}

	return nil
}
