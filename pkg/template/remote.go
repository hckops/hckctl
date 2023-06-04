package template

import (
	"github.com/go-git/go-git/v5"
	"github.com/hckops/hckctl/pkg/util"
	"github.com/rs/zerolog/log"
)

type RemoteTemplateOpts struct {
	SourceCacheDir string
	SourceUrl      string
	Revision       string
	Name           string
	Format         string
}

type RemoteTemplateLoader struct {
	opts *RemoteTemplateOpts
}

func NewRemoteTemplateLoader(opts *RemoteTemplateOpts) *RemoteTemplateLoader {
	return &RemoteTemplateLoader{
		opts: opts,
	}
}

func (l *RemoteTemplateLoader) Load() (*TemplateValue, error) {
	return nil, nil
}

func LoadRemoteTemplate(opts *RemoteTemplateOpts) (*TemplateValue, error) {

	if util.IsPathNotExist(opts.SourceCacheDir) {
		repository, err := git.PlainClone(opts.SourceCacheDir, false, &git.CloneOptions{
			URL:               opts.SourceUrl,
			RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
		})
		log.Error().Err(err).Msg("clone error")
		ref, err := repository.Head()
		log.Error().Err(err).Msg("head error")
		log.Info().Msgf("head=%v", ref.Hash())
	}

	// check if git source exists
	// > if not download --> (if fail exit)
	// > otherwise update --> (if fail WARN offline but continue)
	// load local template
	// list all templates

	return nil, nil
}
