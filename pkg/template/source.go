package template

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/pkg/errors"
)

// TODO allow offline refresh if the repository already exists
// TODO add lock to support concurrent requests
// TODO rename remote to git
// TODO add remote http

type SourceOptions struct {
	SourceCacheDir string
	SourceUrl      string
	SourceRevision string // default branch
	Revision       string
	AllowOffline   bool
}

func refreshSource(opts *SourceOptions) error {

	// first time clone repo always with default revision
	// assume that path doesn't exist, or it's empty
	if _, err := git.PlainClone(opts.SourceCacheDir, false, &git.CloneOptions{
		URL:           opts.SourceUrl,
		ReferenceName: plumbing.NewBranchReferenceName(opts.SourceRevision),
	}); err != nil && err != git.ErrRepositoryAlreadyExists {
		return errors.Wrap(err, "unable to clone repository")
	}

	// access repository
	repository, err := git.PlainOpen(opts.SourceCacheDir)
	if err != nil {
		return errors.Wrap(err, "unable to open repository")
	}
	workTree, err := repository.Worktree()
	if err != nil {
		return errors.Wrap(err, "unable to access repository")
	}

	// fetch latest changes
	// https://git-scm.com/book/en/v2/Git-Internals-The-Refspec
	if err := repository.Fetch(&git.FetchOptions{
		RefSpecs: []config.RefSpec{"+refs/*:refs/*"},
	}); err != nil && err != git.NoErrAlreadyUpToDate {
		return errors.Wrap(err, "unable to fetch repository")
	}

	// resolve revision (branch|tag|sha) to hash
	hash, err := repository.ResolveRevision(plumbing.Revision(opts.Revision))
	if err != nil {
		return errors.Wrap(err, "unable to resolve revision")
	}

	// update latest revision
	if err := workTree.Checkout(&git.CheckoutOptions{Hash: *hash, Force: true}); err != nil {
		return errors.Wrap(err, "unable to checkout revision")
	}

	return nil
}
