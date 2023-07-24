package template

import (
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/pkg/errors"
)

type GitSourceOptions struct {
	CacheBaseDir    string
	RepositoryUrl   string
	DefaultRevision string
	Revision        string
	AllowOffline    bool
}

func (s *GitSourceOptions) CacheDirName() string {
	// extracts repository name
	index := strings.LastIndex(s.RepositoryUrl, "/")
	return strings.TrimSuffix(strings.TrimPrefix(s.RepositoryUrl[index:], "/"), filepath.Ext(s.RepositoryUrl))
}

func (s *GitSourceOptions) CachePath() string {
	return filepath.Join(s.CacheBaseDir, s.CacheDirName())
}

// returns the resolved commit sha
func refreshSource(opts *GitSourceOptions) (string, error) {

	// first time clone repo always with default revision
	// assume that path doesn't exist, or it's empty
	if _, err := git.PlainClone(opts.CachePath(), false, &git.CloneOptions{
		URL:           opts.RepositoryUrl,
		ReferenceName: plumbing.NewBranchReferenceName(opts.DefaultRevision),
	}); err != nil && err != git.ErrRepositoryAlreadyExists {
		return "", errors.Wrap(err, "unable to clone repository")
	}

	// access repository
	repository, err := git.PlainOpen(opts.CachePath())
	if err != nil {
		return "", errors.Wrap(err, "unable to open repository")
	}
	workTree, err := repository.Worktree()
	if err != nil {
		return "", errors.Wrap(err, "unable to access repository")
	}

	// fetch latest changes, ignore error if is offline
	// https://git-scm.com/book/en/v2/Git-Internals-The-Refspec
	if err := repository.Fetch(&git.FetchOptions{
		RefSpecs: []config.RefSpec{"+refs/*:refs/*"},
	}); err != nil && err != git.NoErrAlreadyUpToDate && !opts.AllowOffline {
		return "", errors.Wrap(err, "unable to fetch repository")
	}

	// resolve revision (branch|tag|sha) to hash
	hash, err := repository.ResolveRevision(plumbing.Revision(opts.Revision))
	if err != nil {
		return "", errors.Wrap(err, "unable to resolve revision")
	}

	// update latest revision
	if err := workTree.Checkout(&git.CheckoutOptions{Hash: *hash, Force: true}); err != nil {
		return "", errors.Wrap(err, "unable to checkout revision")
	}

	return hash.String(), nil
}
