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
	LabelTemplateGitName     = "com.hckops.template.git.name"
	LabelTemplateGitUrl      = "com.hckops.template.git.url"
	LabelTemplateGitRevision = "com.hckops.template.git.revision"
	LabelTemplateCommonPath  = "com.hckops.template.common.path"
	LabelBoxSize             = "com.hckops.box.size"
)

func NewLocalLabels() BoxLabels {
	return map[string]string{
		LabelSchemaKind:    schema.KindBoxV1.String(),
		LabelTemplateLocal: "true",
	}
}

func NewGitLabels(name, url, revision string) BoxLabels {
	return map[string]string{
		LabelSchemaKind:          schema.KindBoxV1.String(),
		LabelTemplateGit:         "true",
		LabelTemplateGitName:     name,
		LabelTemplateGitUrl:      url,
		LabelTemplateGitRevision: revision,
	}
}

func (l BoxLabels) AddLabels(path string, size ResourceSize) BoxLabels {
	labels := map[string]string{
		LabelTemplateCommonPath: path,
		LabelBoxSize:            strings.ToLower(size.String()),
	}
	// merge labels
	maps.Copy(labels, l)

	return labels
}
