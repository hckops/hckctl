package loader

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/pkg/errors"

	"github.com/hckops/hckctl/pkg/command/common"
	"github.com/hckops/hckctl/pkg/util"
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

	if err := l.refreshRevision(); err != nil {
		return nil, errors.Wrap(err, "invalid template revision")
	}

	path, err := l.resolvePath()
	if err != nil {
		return nil, errors.Wrap(err, "invalid template name")
	}

	localOpts := &LocalTemplateOpts{Path: path, Format: l.opts.Format}
	return NewLocalTemplateLoader(localOpts).Load()
}

func (l *RemoteTemplateLoader) resolvePath() (string, error) {

	// list all base directories
	var directories []string
	if err := filepath.Walk(l.opts.SourceCacheDir, func(path string, info os.FileInfo, err error) error {
		// excludes "docker" and hidden directories i.e. ".git", ".github"
		if info.IsDir() &&
			!strings.HasPrefix(path, fmt.Sprintf("%s/.", l.opts.SourceCacheDir)) &&
			!strings.HasPrefix(path, fmt.Sprintf("%s/docker", l.opts.SourceCacheDir)) {
			directories = append(directories, path)
		}
		return nil
	}); err != nil {
		return "", errors.Wrap(err, "unable to resolve path directories")
	}

	// paths to attempts
	var paths []string
	for _, directory := range directories {
		paths = append(paths,
			fmt.Sprintf("%s/%s", directory, l.opts.Name),
			fmt.Sprintf("%s/%s.yml", directory, l.opts.Name),
			fmt.Sprintf("%s/%s.yaml", directory, l.opts.Name))
	}

	// match name with predefined paths
	var match []string
	for _, path := range paths {
		if _, err := util.ReadFile(path); err == nil {
			match = append(match, path)
		}
	}

	if len(match) == 1 {
		return match[0], nil
	} else if len(match) > 1 {
		return "", fmt.Errorf("unexpected match of multiple templates %v", match)
	}
	return "", errors.New("path not found")
}

// TODO allow offline refresh if the repository already exists
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

	// update previous revision and fetch latest changes
	// set automatically default revision in case override is invalid
	if err := workTree.Pull(&git.PullOptions{}); err != nil && err != git.NoErrAlreadyUpToDate {
		return errors.Wrap(err, "unable to update previous revision")
	}

	// resolve supported revision to hash
	hash, err := repository.ResolveRevision(plumbing.Revision(l.opts.Revision))
	if err != nil {
		return errors.Wrap(err, "unable to resolve hash revision")
	}

	// checkout latest revision
	if err := workTree.Checkout(&git.CheckoutOptions{Hash: *hash}); err != nil {
		return errors.Wrap(err, "unable to checkout revision")
	}

	return nil
}
