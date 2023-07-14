package model

import (
	"fmt"

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

func BoxLabel() string {
	return fmt.Sprintf("%s=%s", LabelSchemaKind, schema.KindBoxV1.String())
}

func NewLocalLabels() map[string]string {
	return map[string]string{
		LabelSchemaKind:    schema.KindBoxV1.String(),
		LabelTemplateLocal: "true",
	}
}

func NewGitLabels(name, url, revision string) map[string]string {
	return map[string]string{
		LabelSchemaKind:          schema.KindBoxV1.String(),
		LabelTemplateGit:         "true",
		LabelTemplateGitName:     name,
		LabelTemplateGitUrl:      url,
		LabelTemplateGitRevision: revision,
	}
}
