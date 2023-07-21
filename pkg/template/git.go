package template

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"

	"github.com/hckops/hckctl/pkg/util"
)

type GitSource[T TemplateType] struct {
	opts *GitSourceOptions
	name string
}

func (src *GitSource[T]) Parse() (*RawTemplate, error) {
	return readGitTemplate(src.opts, src.name)
}

func (src *GitSource[T]) Validate() ([]*TemplateValidated, error) {
	wildcard := fmt.Sprintf("%s/**/*.{yml,yaml}", src.opts.CacheBaseDir)
	return readGitTemplates(src.opts, wildcard)
}

func (src *GitSource[T]) Read() (*TemplateInfo[T], error) {
	return readGitTemplateInfo[T](src.opts, src.name)
}

func readGitTemplate(opts *GitSourceOptions, name string) (*RawTemplate, error) {
	if path, _, err := resolvePathWithRevision(opts, name); err != nil {
		return nil, err
	} else {
		return readRawTemplate(path)
	}
}

func resolvePathWithRevision(opts *GitSourceOptions, name string) (string, string, error) {
	hash, err := refreshSource(opts)
	if err != nil {
		return "", "", errors.Wrap(err, "invalid template source")
	}

	path, err := resolvePath(opts.CacheBaseDir, name)
	if err != nil {
		return "", "", errors.Wrap(err, "invalid template name")
	}
	return path, hash, nil
}

func resolvePath(sourceCacheDir, name string) (string, error) {

	// list all base directories
	var directories []string
	if err := filepath.Walk(sourceCacheDir, func(path string, info os.FileInfo, err error) error {
		// excludes "docker" and hidden directories e.g. ".git", ".github"
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

func readGitTemplates(opts *GitSourceOptions, wildcard string) ([]*TemplateValidated, error) {
	if _, err := refreshSource(opts); err != nil {
		return nil, errors.Wrap(err, "invalid template revision")
	}
	return readTemplates(wildcard)
}

func readGitTemplateInfo[T TemplateType](opts *GitSourceOptions, name string) (*TemplateInfo[T], error) {
	path, hash, err := resolvePathWithRevision(opts, name)
	if err != nil {
		return nil, err
	}

	return readTemplateInfo[T](Git, path, hash)
}
