package template

import (
	"fmt"
	"strings"

	"github.com/MakeNowJust/heredoc"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/thediveo/enumflag/v2"

	"github.com/hckops/hckctl/internal/command/common"
	"github.com/hckops/hckctl/internal/command/common/flag"
	"github.com/hckops/hckctl/internal/command/config"
	. "github.com/hckops/hckctl/pkg/template"
)

type templateCmdOptions struct {
	configRef   *config.ConfigRef
	formatFlag  formatFlag
	sourceFlag  *flag.SourceFlag
	offlineFlag bool
}

func NewTemplateCmd(configRef *config.ConfigRef) *cobra.Command {

	opts := &templateCmdOptions{
		configRef: configRef,
	}

	command := &cobra.Command{
		Use:   "template [name]",
		Short: "Validate and print templates",
		Example: heredoc.Doc(`

			# prints a git template (supports multiple formats)
			hckctl template alpine
			hckctl template base/parrot
			hckctl template box/base/arch
			hckctl template box/vulnerable/dvwa.yml

			# prints a specific git version (branch|tag|sha) of the template
			hckctl template alpine --revision main

			# prints a template in json format (default yaml)
			hckctl template alpine --format json

			# prints the latest available template cached
			hckctl template alpine --offline

			# validates and prints a local template
			hckctl template ../megalopolis/box/base/alpine.yml --local
		`),
		Args: cobra.ExactArgs(1),
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
	opts.sourceFlag = flag.AddTemplateSourceFlag(command)

	// --offline or --local
	flag.AddOfflineFlag(command, &opts.offlineFlag)
	command.MarkFlagsMutuallyExclusive(flag.OfflineFlagName, flag.LocalFlagName)

	command.AddCommand(NewTemplateListCmd(configRef))
	command.AddCommand(NewTemplateValidateCmd())

	return command
}

func (opts *templateCmdOptions) run(cmd *cobra.Command, args []string) error {
	format := opts.formatFlag.String()

	// TODO add remote url
	if opts.sourceFlag.Local {
		path := args[0]
		log.Debug().Msgf("print local template: path=%s", path)

		return printTemplate(NewLocalValidator(path), format)
	} else {
		name := args[0]
		sourceOpts := &GitSourceOptions{
			CacheBaseDir:    opts.configRef.Config.Template.CacheDir,
			RepositoryUrl:   common.TemplateSourceUrl,
			DefaultRevision: common.TemplateSourceRevision,
			Revision:        opts.sourceFlag.Revision,
			AllowOffline:    opts.offlineFlag,
		}
		log.Debug().Msgf("print git template: name=%s revision=%s offline=%v", name, opts.sourceFlag.Revision, opts.offlineFlag)

		return printTemplate(NewGitValidator(sourceOpts, name), format)
	}
}

func printTemplate(src SourceValidator, format string) error {

	value, err := src.Parse()
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

func formatTemplate(value *RawTemplate, format string) (string, error) {
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
