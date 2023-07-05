package template

import (
	"fmt"
	"strings"

	"github.com/MakeNowJust/heredoc"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/hckops/hckctl/pkg/command/common"
	"github.com/hckops/hckctl/pkg/command/common/flag"
	"github.com/hckops/hckctl/pkg/command/config"
	. "github.com/hckops/hckctl/pkg/template"
)

type templateListCmdOptions struct {
	configRef *config.ConfigRef
	revision  string
	kind      string // TODO filter comma separated list e.g. "box,lab"
	order     string // TODO sort output
	column    string // TODO output only specific fields
}

func NewTemplateListCmd(configRef *config.ConfigRef) *cobra.Command {

	opts := &templateListCmdOptions{
		configRef: configRef,
	}

	command := &cobra.Command{
		Use:   "list",
		Short: "List public templates",
		Example: heredoc.Doc(`

			# list all templates
			hckctl template list

			# list templates using a specific git version (branch|tag|sha)
			hckctl template list --revision main
		`),
		RunE: opts.run,
	}

	// --revision
	flag.AddRevisionFlag(command, &opts.revision)

	return command
}

func (opts *templateListCmdOptions) run(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return templateList(opts.configRef.Config.Template.CacheDir, opts.revision)
	} else {
		cmd.HelpFunc()(cmd, args)
	}
	return nil
}

func templateList(sourceDir string, revision string) error {
	sourceOpts := &SourceOptions{
		SourceCacheDir: sourceDir,
		SourceUrl:      common.TemplateSourceUrl,
		SourceRevision: common.TemplateSourceRevision,
		Revision:       revision,
	}
	// name is overridden with custom wildcard
	if validations, err := NewRemoteSource(sourceOpts, "").ReadTemplates(); err != nil {
		log.Warn().Err(err).Msg("error listing templates")
		return errors.New("error")

	} else {
		var total int
		for _, validation := range validations {
			if validation.IsValid {
				total = total + 1
				// remove prefix and suffix
				prettyPath := strings.NewReplacer(
					fmt.Sprintf("%s/", sourceDir), "",
					".yml", "",
					".yaml", "",
				).Replace(validation.Path)

				log.Debug().Msgf("found template: kind=%s pretty=%s path=%s", validation.Value.Kind.String(), prettyPath, validation.Path)
				fmt.Println(fmt.Sprintf("%s\t%s", validation.Value.Kind.String(), prettyPath))
			} else {
				log.Warn().Msgf("skipping invalid template: path=%s", validation.Path)
			}
		}
		log.Debug().Msgf("total templates: %d", total)
		fmt.Println(fmt.Sprintf("total: %d", total))
	}
	return nil
}
