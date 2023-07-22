package model

import (
	"fmt"
	"path/filepath"
	"strings"

	"golang.org/x/exp/maps"

	"github.com/hckops/hckctl/pkg/schema"
)

type BoxLabels map[string]string

const (
	LabelSchemaKind          = "com.hckops.schema.kind"
	LabelTemplateLocal       = "com.hckops.template.local"
	LabelTemplateGit         = "com.hckops.template.git"
	LabelTemplateGitUrl      = "com.hckops.template.git.url"
	LabelTemplateGitRevision = "com.hckops.template.git.revision"
	LabelTemplateGitCommit   = "com.hckops.template.git.commit"
	LabelTemplateGitDir      = "com.hckops.template.git.dir"
	LabelTemplateGitName     = "com.hckops.template.git.name"
	LabelTemplateCachePath   = "com.hckops.template.cache.path"
	LabelBoxSize             = "com.hckops.box.size"
)

func NewLocalLabels() BoxLabels {
	return map[string]string{
		LabelSchemaKind:    schema.KindBoxV1.String(),
		LabelTemplateLocal: "true",
	}
}

func NewGitLabels(url, revision, dir string) BoxLabels {
	return map[string]string{
		LabelSchemaKind:          schema.KindBoxV1.String(),
		LabelTemplateGit:         "true",
		LabelTemplateGitUrl:      url,
		LabelTemplateGitRevision: revision,
		LabelTemplateGitDir:      dir,
	}
}
func (l BoxLabels) AddLocalLabels(size ResourceSize, path string) BoxLabels {
	if _, err := l.exist(LabelTemplateLocal); err != nil {
		return l
	}
	return mergeWithCommonLabels(l, size, path)
}

func (l BoxLabels) AddGitLabels(size ResourceSize, path string, commit string) BoxLabels {
	if _, err := l.exist(LabelTemplateGit); err != nil {
		return l
	}

	l[LabelTemplateGitCommit] = commit

	templatePath := strings.SplitAfter(path, l[LabelTemplateGitDir])
	// unsafe assume valid path
	name := strings.TrimSuffix(strings.TrimPrefix(templatePath[1], "/"), filepath.Ext(path))
	l[LabelTemplateGitName] = name

	return mergeWithCommonLabels(l, size, path)
}

func mergeWithCommonLabels(labels BoxLabels, size ResourceSize, path string) BoxLabels {
	l := map[string]string{
		LabelTemplateCachePath: path, // absolute path
		LabelBoxSize:           strings.ToLower(size.String()),
	}

	// merge labels
	maps.Copy(labels, l)

	return labels
}

func (l BoxLabels) exist(name string) (string, error) {
	if label, ok := l[name]; !ok {
		return "", fmt.Errorf("label %s not found", name)
	} else {
		return label, nil
	}
}

func (l BoxLabels) ToSize() (ResourceSize, error) {
	if label, err := l.exist(LabelBoxSize); err != nil {
		return Small, err
	} else {
		return ExistResourceSize(label)
	}
}

func (l BoxLabels) ToCachedTemplateInfo() *CachedTemplateInfo {
	// TODO local and remote
	if _, err := l.exist(LabelTemplateLocal); err != nil {
		return nil
	}
	return &CachedTemplateInfo{Path: l[LabelTemplateCachePath]}
}

func (l BoxLabels) ToGitTemplateInfo() *GitTemplateInfo {
	if _, err := l.exist(LabelTemplateGit); err != nil {
		return nil
	}

	return &GitTemplateInfo{
		Url:      l[LabelTemplateGitUrl],
		Revision: l[LabelTemplateGitRevision],
		Commit:   l[LabelTemplateGitCommit],
		Name:     l[LabelTemplateGitName],
	}
}
