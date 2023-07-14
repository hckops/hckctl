package template

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"

	"github.com/hckops/hckctl/pkg/util"
)

func readGitTemplate(opts *SourceOptions, name string) (*TemplateValue, error) {
	if path, err := resolvePathWithRevision(opts, name); err != nil {
		return nil, err
	} else {
		return readTemplate(path)
	}
}

func resolvePathWithRevision(opts *SourceOptions, name string) (string, error) {
	if err := refreshSource(opts); err != nil {
		return "", errors.Wrap(err, "invalid template source")
	}

	path, err := resolvePath(opts.SourceCacheDir, name)
	if err != nil {
		return "", errors.Wrap(err, "invalid template name")
	}

	return path, nil
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

func readGitTemplates(opts *SourceOptions, wildcard string) ([]*TemplateValidated, error) {
	if err := refreshSource(opts); err != nil {
		return nil, errors.Wrap(err, "invalid template revision")
	}
	return readTemplates(wildcard)
}

func readGitBoxTemplate(opts *SourceOptions, name string) (*BoxTemplate, error) {
	if path, err := resolvePathWithRevision(opts, name); err != nil {
		return nil, err
	} else {
		return readBoxTemplate(path)
	}
}

func readGitLabTemplate(opts *SourceOptions, name string) (*LabTemplate, error) {
	if path, err := resolvePathWithRevision(opts, name); err != nil {
		return nil, err
	} else {
		return readLabTemplate(path)
	}
}
