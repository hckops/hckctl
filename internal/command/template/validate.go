package template

import (
	"fmt"

	"github.com/MakeNowJust/heredoc"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	. "github.com/hckops/hckctl/pkg/template"
)

// TODO see common list options: kind (compare expected), order, column
type templateValidateCmdOptions struct{}

func NewTemplateValidateCmd() *cobra.Command {

	opts := &templateValidateCmdOptions{}

	command := &cobra.Command{
		Use:   "validate [path]",
		Short: "Validate one or more templates",
		Example: heredoc.Doc(`

			# validates a local template
			hckctl template validate ../megalopolis/box/base/alpine.yml

			# validates all templates in the given path (supports wildcards)
			hckctl template validate "../megalopolis/**/*.{yml,yaml}"
			hckctl template validate "../megalopolis/**/*alpine*"
		`),
		Args: cobra.ExactArgs(1),
		RunE: opts.run,
	}

	return command
}

func (opts *templateValidateCmdOptions) run(cmd *cobra.Command, args []string) error {
	return templateValidate(args[0])
}

// TODO color output
func templateValidate(path string) error {
	src := NewLocalValidator(path)

	// attempt single file validation
	if templateValue, err := src.Parse(); err == nil {
		printValidTemplate(path, templateValue)

		// attempt wildcard validation
	} else if validations, err := src.Validate(); err == nil {
		for _, validation := range validations {
			if validation.IsValid {
				printValidTemplate(validation.Path, validation.Value)
			} else {
				printInvalidTemplate(validation.Path)
			}
		}
		log.Debug().Msgf("validated templates: %d", len(validations))
		fmt.Println(fmt.Sprintf("total: %d", len(validations)))

	} else {
		log.Warn().Err(err).Msgf("error validating template: path=%s", path)
		return errors.New("error")
	}

	return nil
}

func printValidTemplate(path string, value *RawTemplate) {
	log.Debug().Msgf("valid template: kind=%s path=%s", value.Kind.String(), path)
	fmt.Println(fmt.Sprintf("[OK] %s\t%s", value.Kind.String(), path))
}

func printInvalidTemplate(path string) {
	log.Warn().Msgf("invalid template: path=%s", path)
	fmt.Println(fmt.Sprintf("[KO] >>>>>>>>>>\t%s", path))
}
