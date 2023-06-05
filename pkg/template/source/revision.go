package source

import (
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/pkg/errors"
)

// TODO allow offline refresh if the repository already exists
// TODO add lock to support concurrency requests

type RevisionOpts struct {
	SourceCacheDir string
	SourceUrl      string
	SourceRevision string // default branch
	Revision       string
}

func refreshRevision(opts *RevisionOpts) error {

	// first time clone repo always with default revision
	// assume that path doesn't exist, or it's empty
	if _, err := git.PlainClone(opts.SourceCacheDir, false, &git.CloneOptions{
		URL:               opts.SourceUrl,
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
		ReferenceName:     plumbing.NewBranchReferenceName(opts.SourceRevision),
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

	// update previous revision and fetch latest changes
	// set automatically default revision in case override is invalid
	if err := workTree.Pull(&git.PullOptions{}); err != nil && err != git.NoErrAlreadyUpToDate {
		return errors.Wrap(err, "unable to update previous revision")
	}

	// resolve revision (branch|tag|sha) to hash
	hash, err := repository.ResolveRevision(plumbing.Revision(opts.Revision))
	if err != nil {
		return errors.Wrap(err, "unable to resolve hash revision")
	}

	// checkout latest revision
	if err := workTree.Checkout(&git.CheckoutOptions{Hash: *hash}); err != nil {
		return errors.Wrap(err, "unable to checkout revision")
	}

	return nil
}
