package template

import (
	"fmt"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/hckops/hckctl/pkg/command/common"
	"github.com/hckops/hckctl/pkg/command/config"
	"github.com/hckops/hckctl/pkg/template"
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
		Short: "list available templates",
		Example: heredoc.Doc(`

			# list all templates
			hckctl template list

			# list templates from a specific version (branch|tag|sha)
			hckctl template list --revision main
		`),
		RunE: opts.run,
	}

	// --revision
	common.AddRevisionFlag(command, &opts.revision)

	return command
}

func (opts *templateListCmdOptions) run(cmd *cobra.Command, args []string) error {

	_ = template.RevisionOpts{
		SourceCacheDir: opts.configRef.Config.Template.CacheDir,
		SourceUrl:      common.TemplateSourceUrl,
		SourceRevision: common.TemplateSourceRevision,
		Revision:       opts.revision,
	}

	fmt.Println("not implemented")
	return nil
}
