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
			hckctl template validate ../megalopolis/boxes/official/alpine.yml

			# validates all templates in the given path (supports wildcard)
			hckctl template validate "../megalopolis/**/*.{yml,yaml}"
			hckctl template validate "../megalopolis/**/*alpine*"
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
	src := NewLocalSource(path)

	// attempt single file validation
	if templateValue, err := src.ReadTemplate(); err == nil {
		printValidTemplate(path, templateValue)

		// attempt wildcard validation
	} else if validations, err := src.ReadTemplates(); err == nil {
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

func printValidTemplate(path string, value *TemplateValue) {
	log.Debug().Msgf("valid template: kind=%s path=%s", value.Kind.String(), path)
	fmt.Println(fmt.Sprintf("[OK] %s\t%s", value.Kind.String(), path))
}

func printInvalidTemplate(path string) {
	log.Warn().Msgf("invalid template: path=%s", path)
	fmt.Println(fmt.Sprintf("[KO] >>>>>>>>>>\t%s", path))
}
