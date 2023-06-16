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
	. "github.com/hckops/hckctl/pkg/template"
)

type templateCmdOptions struct {
	configRef  *config.ConfigRef
	formatFlag formatFlag
	sourceFlag *common.SourceFlag
}

func NewTemplateCmd(configRef *config.ConfigRef) *cobra.Command {

	opts := &templateCmdOptions{
		configRef: configRef,
	}

	command := &cobra.Command{
		Use:   "template [name]",
		Short: "Validate and print templates",
		Example: heredoc.Doc(`

			# prints a remote template (supports multiple formats)
			hckctl template alpine
			hckctl template base/parrot
			hckctl template boxes/vulnerable/dvwa
			hckctl template labs/ctf/htb-alpine.yml

			# prints a specific version (branch|tag|sha)
			hckctl template alpine --revision main

			# prints a template in json format (default yaml)
			hckctl template alpine --format json

			# validates and prints a local template
			hckctl template ../megalopolis/boxes/base/alpine.yml --local
		`),
		RunE: opts.run,
	}

	const (
		formatFlagName = "format"
	)
	// --format (enum)
	formatValue := enumflag.New(&opts.formatFlag, formatFlagName, formatIds, enumflag.EnumCaseInsensitive)
	formatUsage := fmt.Sprintf("output format, one of %s", strings.Join(formatValues(), "|"))
	command.Flags().Var(formatValue, formatFlagName, formatUsage)

	// --revision or --local
	opts.sourceFlag = common.AddTemplateSourceFlag(command)

	command.AddCommand(NewTemplateListCmd(configRef))
	command.AddCommand(NewTemplateValidateCmd())

	return command
}

func (opts *templateCmdOptions) run(cmd *cobra.Command, args []string) error {
	format := opts.formatFlag.String()

	if len(args) == 1 && opts.sourceFlag.Local {
		path := args[0]
		log.Debug().Msgf("print local template: path=%s", path)

		return printTemplate(NewLocalSource(path), format)

	} else if len(args) == 1 {
		name := args[0]
		revisionOpts := &RevisionOpts{
			SourceCacheDir: opts.configRef.Config.Template.CacheDir,
			SourceUrl:      common.TemplateSourceUrl,
			SourceRevision: common.TemplateSourceRevision,
			Revision:       opts.sourceFlag.Revision,
		}
		log.Debug().Msgf("print remote template: name=%s revision=%s", name, opts.sourceFlag.Revision)

		return printTemplate(NewRemoteSource(revisionOpts, name), format)

	} else {
		cmd.HelpFunc()(cmd, args)
	}
	return nil
}

func printTemplate(src TemplateSource, format string) error {

	value, err := src.ReadTemplate()
	if err != nil {
		log.Warn().Err(err).Msg("error reading template")
		return errors.New("invalid template")
	}

	if formatted, err := formatTemplate(value, format); err != nil {
		log.Warn().Err(err).Msg("error printing template")
		return errors.New("format error")
	} else {
		log.Debug().Msgf("print template: kind=%s format=%s\n%s", value.Kind.String(), format, formatted)
		fmt.Print(formatted)
	}
	return nil
}

func formatTemplate(value *TemplateValue, format string) (string, error) {
	switch format {
	case jsonFlag.String():
		if jsonValue, err := value.ToJson(); err != nil {
			return "", err
		} else {
			// add newline
			return fmt.Sprintf("%s\n", jsonValue.Data), nil
		}
	default:
		if yamlValue, err := value.ToYaml(); err != nil {
			return "", err
		} else {
			return yamlValue.Data, nil
		}
	}
}
