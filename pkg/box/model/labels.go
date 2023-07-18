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
	LabelTemplateCommonPath  = "com.hckops.template.common.path"
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

func (l BoxLabels) AddLabels(path string, commit string, size ResourceSize) BoxLabels {
	labels := map[string]string{
		LabelTemplateCommonPath: path, // absolute path
		LabelBoxSize:            strings.ToLower(size.String()),
	}

	// add labels only to git template
	if _, exist := l[LabelTemplateGit]; exist {
		l[LabelTemplateGitCommit] = commit

		templatePath := strings.SplitAfter(path, l[LabelTemplateGitDir])
		name := strings.TrimSuffix(strings.TrimPrefix(templatePath[1], "/"), filepath.Ext(path))
		l[LabelTemplateGitName] = name
	}

	// merge labels
	maps.Copy(labels, l)

	return labels
}

// TODO test
func (l BoxLabels) exists(name string) (string, error) {
	if label, ok := l[name]; !ok {
		return "", fmt.Errorf("label %s not found", name)
	} else {
		return label, nil
	}
}

// TODO test
func (l BoxLabels) ToSize() (ResourceSize, error) {
	if label, err := l.exists(LabelBoxSize); err != nil {
		return Small, err
	} else {
		return ExistResourceSize(label)
	}
}

// TODO test
func (l BoxLabels) ToBoxTemplateInfo() (*BoxTemplateInfo, error) {
	if _, err := l.exists(LabelTemplateGit); err != nil {
		return nil, err
	}

	return &BoxTemplateInfo{
		Url:      l[LabelTemplateGitUrl],
		Revision: l[LabelTemplateGitRevision],
		Commit:   l[LabelTemplateGitCommit],
		Name:     l[LabelTemplateGitName],
	}, nil
}
