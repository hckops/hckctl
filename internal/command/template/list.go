package template

import (
	"fmt"
	"strings"

	"github.com/MakeNowJust/heredoc"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/hckops/hckctl/internal/command/common"
	"github.com/hckops/hckctl/internal/command/common/flag"
	"github.com/hckops/hckctl/internal/command/config"
	. "github.com/hckops/hckctl/pkg/template"
)

type templateListCmdOptions struct {
	configRef *config.ConfigRef
	revision  string
	offline   bool
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

			# list templates cached
			hckctl template list --offline
		`),
		Args: cobra.NoArgs,
		RunE: opts.run,
	}

	// --revision
	flag.AddRevisionFlag(command, &opts.revision)
	// --offline
	flag.AddOfflineFlag(command, &opts.offline)

	return command
}

func (opts *templateListCmdOptions) run(cmd *cobra.Command, args []string) error {
	return templateList(opts.configRef.Config.Template.CacheDir, opts.revision, opts.offline)
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
				// remove prefix and suffix
				prettyPath := strings.NewReplacer(
					fmt.Sprintf("%s/", cacheDir), "",
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
