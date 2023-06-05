package template

import (
	"fmt"

	"github.com/MakeNowJust/heredoc"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/hckops/hckctl/pkg/template/loader"
)

type templateValidateCmdOptions struct {
	kind string
}

func NewTemplateValidateCmd() *cobra.Command {

	opts := &templateValidateCmdOptions{}

	command := &cobra.Command{
		Use:   "validate [path]",
		Short: "validate templates",
		Example: heredoc.Doc(`

			# validates local template
			hckctl template validate ../megalopolis/boxes/official/alpine.yml

			# validates all templates in the given path (supports wildcard)
			hckctl template validate ../megalopolis/boxes/*
		`),
		RunE: opts.run,
	}

	// TODO implement flag: parse value and compare with result
	// command.Flags().StringVarP(&opts.kind, "kind", common.NoneFlagShortHand, "", "expected template kind")

	return command
}

// TODO check if file or use "filepath.Glob" - make sure to include (optional) yaml/yml
// https://gosamples.dev/list-files

// TODO use in gh-action: validate multiple templates in the given path (not only single file) + add regex filter
func (opts *templateValidateCmdOptions) run(cmd *cobra.Command, args []string) error {
	if len(args) == 1 {
		return validateLocalTemplate(args[0])
	} else {
		cmd.HelpFunc()(cmd, args)
	}
	return nil
}

func validateLocalTemplate(path string) error {
	if templateValue, err := loader.NewDefaultLocalTemplateLoader(path).Load(); err != nil {
		log.Warn().Err(err).Msgf("error validating local template: path=%s", path)
		return errors.New("KO")
	} else {
		log.Debug().Msgf("valid template: path=%s kind=%s", path, templateValue.Kind.String())
		fmt.Println(fmt.Sprintf("OK %s", templateValue.Kind.String()))
	}
	return nil
}
