package old

import (
	"fmt"
	box "github.com/hckops/hckctl/pkg/box/model"
	lab "github.com/hckops/hckctl/pkg/lab/model"
	"github.com/hckops/hckctl/pkg/template"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"

	"github.com/hckops/hckctl/pkg/util"
)

func readGitTemplate(opts *template.GitSourceOptions, name string) (*RawTemplate, error) {
	if path, _, err := resolvePathWithRevision(opts, name); err != nil {
		return nil, err
	} else {
		return readRawTemplate(path)
	}
}

func resolvePathWithRevision(opts *template.GitSourceOptions, name string) (string, string, error) {
	hash, err := template.refreshSource(opts)
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

func readGitTemplates(opts *template.GitSourceOptions, wildcard string) ([]*TemplateValidated, error) {
	if _, err := template.refreshSource(opts); err != nil {
		return nil, errors.Wrap(err, "invalid template revision")
	}
	return readTemplates(wildcard)
}

func Zero[T any]() T {
	return *new(T)
}

func readGitTemplateInfo[T TemplateType](opts *template.GitSourceOptions, name string) (*TemplateInfo[T], error) {
	path, hash, err := resolvePathWithRevision(opts, name)
	if err != nil {
		return nil, err
	}

	template, err := readTemplate[T](path)
	if err != nil {
		return nil, err
	}

	info, err := newGitTemplateInfo[T](template, path, hash)

	return info, nil
}

func readGitBoxTemplate(opts *template.GitSourceOptions, name string) (*BoxInfo, error) {
	path, hash, err := resolvePathWithRevision(opts, name)
	if err != nil {
		return nil, err
	}

	template, err := readTemplate[box.BoxV1](path)
	if err != nil {
		return nil, err
	}

	info, err := newGitTemplateInfo[box.BoxV1](template, path, hash)

	return info, nil
}

func readGitLabTemplate(opts *template.GitSourceOptions, name string) (*LabInfo, error) {
	path, hash, err := resolvePathWithRevision(opts, name)
	if err != nil {
		return nil, err
	}

	template, err := readLabTemplate(path)
	if err != nil {
		return nil, err
	}

	info, err := newGitTemplateInfo[lab.LabV1](template., path, hash)
	if err != nil {
		return nil, err
	}

	return newGitTemplateInfo(template, path, hash)
}
