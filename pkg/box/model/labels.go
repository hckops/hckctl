package model

import (
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

func NewGitLabels(url, revision, name string) BoxLabels {
	return map[string]string{
		LabelSchemaKind:          schema.KindBoxV1.String(),
		LabelTemplateGit:         "true",
		LabelTemplateGitUrl:      url,
		LabelTemplateGitRevision: revision,
		LabelTemplateGitName:     name,
	}
}

func (l BoxLabels) AddLabels(path string, commit string, size ResourceSize) BoxLabels {
	labels := map[string]string{
		LabelTemplateCommonPath: path, // absolute path
		LabelBoxSize:            strings.ToLower(size.String()),
	}

	// add commit label only to git template
	if _, exist := l[LabelTemplateGitRevision]; exist {
		l[LabelTemplateGitCommit] = commit
	}

	// merge labels
	maps.Copy(labels, l)

	return labels
}
