package box

import (
	"fmt"
	"strings"

	"github.com/MakeNowJust/heredoc"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/hckops/hckctl/pkg/command/common"
	"github.com/hckops/hckctl/pkg/command/config"
	"github.com/hckops/hckctl/pkg/template/schema"
	"github.com/hckops/hckctl/pkg/template/source"
)

type boxListCmdOptions struct {
	configRef *config.ConfigRef
	revision  string
}

func NewBoxListCmd(configRef *config.ConfigRef) *cobra.Command {

	opts := &boxListCmdOptions{
		configRef: configRef,
	}

	command := &cobra.Command{
		Use:   "list",
		Short: "list available templates",
		Example: heredoc.Doc(`

			# list all templates
			hckctl box list

			# list templates from a specific version (branch|tag|sha)
			hckctl box list --revision main
		`),
		RunE: opts.run,
	}

	// --revision
	common.AddRevisionFlag(command, &opts.revision)

	return command
}

// TODO list running boxes !!! remove revision

func (opts *boxListCmdOptions) run(cmd *cobra.Command, args []string) error {

	revisionOpts := &source.RevisionOpts{
		SourceCacheDir: opts.configRef.Config.Template.CacheDir,
		SourceUrl:      common.TemplateSourceUrl,
		SourceRevision: common.TemplateSourceRevision,
		Revision:       opts.revision,
	}
	if templates, err := source.NewRemoteSource(revisionOpts, "").ReadTemplates(); err != nil {
		log.Warn().Err(err).Msg("error listing boxes")
		return errors.New("error")

	} else {
		var total int
		for _, template := range templates {
			if template.IsValid && template.Value.Kind == schema.KindBoxV1 {
				total = total + 1
				// remove prefix and suffix
				prettyPath := strings.NewReplacer(
					fmt.Sprintf("%s/boxes/", opts.configRef.Config.Template.CacheDir), "",
					".yml", "",
					".yaml", "",
				).Replace(template.Path)

				log.Debug().Msgf("found template: pretty=%s path=%s", prettyPath, template.Path)
				fmt.Println(prettyPath)
			} else {
				log.Warn().Msgf("skipping invalid template: path=%s", template.Path)
			}
		}
		log.Debug().Msgf("total templates: %d", total)
		fmt.Println(fmt.Sprintf("\ntotal: %d", total))
	}
	return nil
}
