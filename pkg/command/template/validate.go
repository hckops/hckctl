package template

import (
	"fmt"

	"github.com/MakeNowJust/heredoc"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/hckops/hckctl/pkg/template/source"
)

// TODO see common list options: kind (compare expected), order, column
type templateValidateCmdOptions struct{}

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

	return command
}

func (opts *templateValidateCmdOptions) run(cmd *cobra.Command, args []string) error {
	if len(args) == 1 {
		return validateTemplate(args[0])
	} else {
		cmd.HelpFunc()(cmd, args)
	}
	return nil
}

// TODO color output
func validateTemplate(path string) error {
	src := source.NewLocalSource(path)

	// attempt single file validation
	if templateValue, err := src.Read(); err == nil {
		printValidTemplate(templateValue)

		// attempt wildcard validation
	} else if validations, err := src.ReadAll(); err == nil {
		log.Debug().Msgf("validated templates: %d", len(validations))
		fmt.Println(fmt.Sprintf("total: %d", len(validations)))

		for _, validation := range validations {
			if validation.IsValid {
				printValidTemplate(templateValue)
			} else {
				printInvalidTemplate(templateValue)
			}
		}

	} else {
		log.Warn().Err(err).Msgf("error validating template: path=%s", path)
		return errors.New("error")
	}

	return nil
}

func printValidTemplate(value *source.TemplateValue) {
	log.Debug().Msgf("valid template: kind=%s path=%s", value.Kind.String(), value.Path)
	fmt.Println(fmt.Sprintf("[OK] %s\t%s", value.Kind.String(), value.Path))
}

func printInvalidTemplate(value *source.TemplateValue) {
	log.Warn().Msgf("invalid template: kind=%s path=%s", value.Kind.String(), value.Path)
	fmt.Println(fmt.Sprintf("[KO] %s\t%s <<<<<<<<<<", value.Kind.String(), value.Path))
}
