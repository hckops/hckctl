package loader

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"

	"github.com/hckops/hckctl/pkg/template"
	"github.com/hckops/hckctl/pkg/util"
)

type RemoteTemplateOpts struct {
	RevisionOpts *template.RevisionOpts
	Name         string
	Format       string
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

	if err := template.RefreshRevision(l.opts.RevisionOpts); err != nil {
		return nil, errors.Wrap(err, "invalid template revision")
	}

	path, err := resolvePath(l.opts.RevisionOpts.SourceCacheDir, l.opts.Name)
	if err != nil {
		return nil, errors.Wrap(err, "invalid template name")
	}

	localOpts := &LocalTemplateOpts{Path: path, Format: l.opts.Format}
	return NewLocalTemplateLoader(localOpts).Load()
}

func resolvePath(sourceCacheDir, name string) (string, error) {

	// list all base directories
	var directories []string
	if err := filepath.Walk(sourceCacheDir, func(path string, info os.FileInfo, err error) error {
		// TODO whitelist vs blacklist
		// excludes "docker" and hidden directories i.e. ".git", ".github"
		if info.IsDir() &&
			!strings.HasPrefix(path, fmt.Sprintf("%s/.", sourceCacheDir)) &&
			!strings.HasPrefix(path, fmt.Sprintf("%s/docker", sourceCacheDir)) {
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
			fmt.Sprintf("%s/%s", directory, name),
			fmt.Sprintf("%s/%s.yml", directory, name),
			fmt.Sprintf("%s/%s.yaml", directory, name))
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
