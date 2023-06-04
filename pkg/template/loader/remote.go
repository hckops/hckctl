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

// TODO test checkout/pull with a repository
// TODO do i need IsPathNotExist or i can just use clone error and revert "utils"

// TODO list files https://gosamples.dev/list-files

func (l *RemoteTemplateLoader) Load() (*TemplateValue, error) {

	path := "/tmp/full/example/nested"
	// c69656300c560af847e2e74da6d97a806cde2572
	// main
	// docker-0.4.9
	revision := "main"

	// first time clone repo always with default revision
	// assume that path doesn't exist, or it's empty
	if _, err := git.PlainClone(path, false, &git.CloneOptions{
		URL:               "https://github.com/hckops/actions", // TODO l.opts.SourceUrl
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
		ReferenceName:     plumbing.NewBranchReferenceName(common.TemplateRevision),
	}); err != nil && err != git.ErrRepositoryAlreadyExists {
		return nil, errors.Wrap(err, "unable to clone repository")
	}

	// access repository
	repository, err := git.PlainOpen(path)
	if err != nil {
		return nil, errors.Wrap(err, "unable to open repository")
	}
	workTree, err := repository.Worktree()
	if err != nil {
		return nil, errors.Wrap(err, "unable to access repository")
	}

	if err := workTree.Pull(&git.PullOptions{}); err != nil && err != git.NoErrAlreadyUpToDate {
		return nil, errors.Wrap(err, "unable to update revision")
	}

	// attempts checkout with the supported revision
	checkoutStrategy := []*git.CheckoutOptions{
		{Branch: plumbing.NewBranchReferenceName(revision)},
		{Hash: plumbing.NewHash(revision)},
		{Branch: plumbing.NewTagReferenceName(revision)},
	}
	var checkoutSuccess bool
	for index, strategy := range checkoutStrategy {
		if err := workTree.Checkout(strategy); err == nil {
			log.Debug().Msgf("use checkout strategy %d", index)
			checkoutSuccess = true
		}
	}
	if !checkoutSuccess {
		return nil, errors.Wrap(err, "unable to checkout revision")
	}

	if head, err := repository.Head(); err != nil {
		return nil, errors.Wrap(err, "unable to verify revision")
	} else {
		log.Debug().Msgf("use revision revision=%s hash=%s", revision, head.Hash())
	}

	// check if git source exists
	// > if not download --> (if fail exit)
	// > otherwise update --> (if fail WARN offline but continue)
	// BUILD PATH -> load local template
	// list all templates

	localOpts := &LocalTemplateOpts{Path: "../megalopolis/boxes/official/alpine.yml", Format: l.opts.Format}
	return NewLocalTemplateLoader(localOpts).Load()
}
