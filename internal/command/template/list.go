package template

import (
	"fmt"
	"github.com/MakeNowJust/heredoc"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/hckops/hckctl/internal/command/common"
	"github.com/hckops/hckctl/internal/command/common/flag"
	"github.com/hckops/hckctl/internal/command/config"
	. "github.com/hckops/hckctl/pkg/template"
)

// TODO add "kind" filter (comma separated list) e.g. "box,lab,task"
// TODO add "order" to sort output
// TODO add "column" to output only specific fields
type templateListCmdOptions struct {
	configRef    *config.ConfigRef
	revisionFlag string
	offlineFlag  bool
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

			# list templates cached
			hckctl template list --offline
		`),
		Args: cobra.NoArgs,
		RunE: opts.run,
	}

	// --revision
	flag.AddTemplateRevisionFlag(command, &opts.revisionFlag)
	// --offline
	flag.AddTemplateOfflineFlag(command, &opts.offlineFlag)

	return command
}

func (opts *templateListCmdOptions) run(cmd *cobra.Command, args []string) error {
	return templateList(opts.configRef.Config.Template.CacheDir, opts.revisionFlag, opts.offlineFlag)
}

func templateList(cacheDir string, revision string, offline bool) error {
	sourceOpts := &GitSourceOptions{
		CacheBaseDir:    cacheDir,
		RepositoryUrl:   common.TemplateSourceUrl,
		DefaultRevision: common.TemplateSourceRevision,
		Revision:        revision,
		AllowOffline:    offline,
	}
	log.Debug().Msgf("list git templates: revision=%s offline=%v", revision, offline)

	// name is overridden with custom wildcard
	if validations, err := NewGitValidator(sourceOpts, "").Validate(); err != nil {
		log.Warn().Err(err).Msg("error listing templates")
		return errors.New("error")

	} else {
		var total int
		for _, validation := range validations {
			if validation.IsValid {
				total = total + 1
				prettyPath := common.PrettyPath(sourceOpts.CachePath(), validation.Path)

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
