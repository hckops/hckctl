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
	"github.com/hckops/hckctl/pkg/template"
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

			# prints remote alpine template
			hckctl template alpine

			# prints specific version (branch|tag|sha) of alpine template
			hckctl template alpine --revision main

			# prints template in json format (default yaml)
			hckctl template alpine --format json

			# validates and prints local template
			hckctl template ../megalopolis/boxes/official/alpine.yml --local
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
	command.Flags().StringVarP(&opts.revision, revisionFlagName, "r", common.TemplateRevision, common.TemplateRevisionUsage)
	command.MarkFlagsMutuallyExclusive(localFlagName, revisionFlagName)

	command.AddCommand(NewTemplateListCmd())
	command.AddCommand(NewTemplateValidateCmd())

	return command
}

func (opts *templateCmdOptions) run(cmd *cobra.Command, args []string) error {
	format := opts.format.value()

	if opts.local {
		return printLocalTemplate(&template.LocalTemplateOpts{
			Path:   args[0],
			Format: format,
		})
	} else if len(args) == 1 {
		return printRemoteTemplate(&template.RemoteTemplateOpts{
			SourceCacheDir: opts.configRef.Config.Template.DirPath,
			SourceUrl:      common.TemplateSourceUrl,
			Revision:       opts.revision,
			Name:           args[0],
			Format:         format,
		})
	} else {
		cmd.HelpFunc()(cmd, args)
	}
	return nil
}

func printLocalTemplate(opts *template.LocalTemplateOpts) error {
	if templateValue, err := template.LoadLocalTemplate(opts); err != nil {
		log.Warn().Err(err).Msgf("error printing template: path=%s", opts.Path)
		return errors.New("invalid")
	} else {
		log.Debug().Msgf("print template: path=%s kind=%s\n%s", opts.Path, templateValue.Kind.String(), templateValue.Data)
		fmt.Print(templateValue.Data)
	}
	return nil
}

// TODO
func printRemoteTemplate(opts *template.RemoteTemplateOpts) error {
	return nil
}

//func printTemplate(loader template.TemplateLoader) error {
//	if templateValue, err := loader.Load(); err != nil {
//		log.Warn().Err(err).Msgf("error printing template: path=%s", loader.Path)
//		return errors.New("invalid")
//	} else {
//		log.Debug().Msgf("print template: path=%s kind=%s\n%s", opts.Path, templateValue.Kind.String(), templateValue.Data)
//		fmt.Print(templateValue.Data)
//	}
//	return nil
//}
