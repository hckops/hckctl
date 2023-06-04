package loader

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"

	"github.com/hckops/hckctl/pkg/command/common"
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

// TODO list files https://gosamples.dev/list-files

func (l *RemoteTemplateLoader) Load() (*TemplateValue, error) {

	if err := l.refreshRevision(); err != nil {
		return nil, errors.Wrap(err, "invalid revision")
	}

	// check if git source exists
	// > if not download --> (if fail exit)
	// > otherwise update --> (if fail WARN offline but continue)
	// BUILD PATH -> load local template
	// list all templates

	localOpts := &LocalTemplateOpts{Path: "../megalopolis/boxes/official/alpine.yml", Format: l.opts.Format}
	return NewLocalTemplateLoader(localOpts).Load()
}

func (l *RemoteTemplateLoader) refreshRevision() error {

	// first time clone repo always with default revision
	// assume that path doesn't exist, or it's empty
	if _, err := git.PlainClone(l.opts.SourceCacheDir, false, &git.CloneOptions{
		URL:               l.opts.SourceUrl,
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
		ReferenceName:     plumbing.NewBranchReferenceName(common.TemplateRevision),
	}); err != nil && err != git.ErrRepositoryAlreadyExists {
		return errors.Wrap(err, "unable to clone repository")
	}

	// access repository
	repository, err := git.PlainOpen(l.opts.SourceCacheDir)
	if err != nil {
		return errors.Wrap(err, "unable to open repository")
	}
	workTree, err := repository.Worktree()
	if err != nil {
		return errors.Wrap(err, "unable to access repository")
	}

	// update current revision
	if err := workTree.Pull(&git.PullOptions{}); err != nil && err != git.NoErrAlreadyUpToDate {
		return errors.Wrap(err, "unable to update revision")
	}

	// attempts checkout for the latest revision
	checkoutStrategy := []*git.CheckoutOptions{
		{Branch: plumbing.NewBranchReferenceName(l.opts.Revision)},
		{Hash: plumbing.NewHash(l.opts.Revision)},
		{Branch: plumbing.NewTagReferenceName(l.opts.Revision)},
	}
	var checkoutSuccess bool
	for index, strategy := range checkoutStrategy {
		if err := workTree.Checkout(strategy); err == nil {
			log.Debug().Msgf("use checkout strategy %d", index)
			checkoutSuccess = true
		}
	}
	if !checkoutSuccess {
		return errors.Wrap(err, "unable to checkout revision")
	}

	if head, err := repository.Head(); err != nil {
		return errors.Wrap(err, "unable to verify revision")
	} else {
		log.Debug().Msgf("use revision revision=%s hash=%s", l.opts.Revision, head.Hash())
	}

	return nil
}
