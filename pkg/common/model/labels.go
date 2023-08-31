package model

import (
	"fmt"
	"path/filepath"
	"strings"

	"golang.org/x/exp/maps"

	boxModel "github.com/hckops/hckctl/pkg/box/model"
	"github.com/hckops/hckctl/pkg/schema"
)

type Labels map[string]string

const (
	LabelSchemaKind          = "com.hckops.schema.kind"
	LabelTemplateLocal       = "com.hckops.template.local"
	LabelTemplateRemote      = "com.hckops.template.remote" // TODO not used
	LabelTemplateGit         = "com.hckops.template.git"
	LabelTemplateGitUrl      = "com.hckops.template.git.url"
	LabelTemplateGitRevision = "com.hckops.template.git.revision"
	LabelTemplateGitCommit   = "com.hckops.template.git.commit"
	LabelTemplateGitDir      = "com.hckops.template.git.dir"
	LabelTemplateGitName     = "com.hckops.template.git.name"
	LabelTemplateCachePath   = "com.hckops.template.cache.path"
	LabelBoxSize             = "com.hckops.box.size"
	LabelTaskSize            = "com.hckops.task.size" // TODO not used
)

func NewBoxLabels() Labels {
	return map[string]string{
		LabelSchemaKind: schema.KindBoxV1.String(),
	}
}

func (l Labels) addLabel(key string, value string) Labels {
	return l.addLabels(Labels{key: value})
}

func (l Labels) addLabels(labels Labels) Labels {
	// merge labels
	maps.Copy(labels, l)
	return labels
}

func (l Labels) exist(name string) (string, error) {
	if label, ok := l[name]; !ok {
		return "", fmt.Errorf("label %s not found", name)
	} else {
		return label, nil
	}
}

func (l Labels) AddDefaultLocal() Labels {
	return l.addLabel(LabelTemplateLocal, "true")
}

func (l Labels) AddDefaultGit(url, revision, dir string) Labels {
	return l.addLabels(Labels{
		LabelTemplateGit:         "true",
		LabelTemplateGitUrl:      url,
		LabelTemplateGitRevision: revision,
		LabelTemplateGitDir:      dir,
	})
}

func (l Labels) AddLocal(path string) Labels {
	if _, err := l.exist(LabelTemplateLocal); err != nil {
		return l
	}
	return l.addLabel(LabelTemplateCachePath, path)
}

func (l Labels) AddGit(path string, commit string) Labels {
	if _, err := l.exist(LabelTemplateGit); err != nil {
		return l
	}

	templatePath := strings.SplitAfter(path, l[LabelTemplateGitDir])
	// unsafe assume valid path
	name := strings.TrimSuffix(strings.TrimPrefix(templatePath[1], "/"), filepath.Ext(path))

	return l.addLabels(Labels{
		LabelTemplateGitCommit: commit,
		LabelTemplateGitName:   name,
		LabelTemplateCachePath: path,
	})
}

func (l Labels) ToCachedTemplateInfo() *CachedTemplateInfo {
	if _, err := l.exist(LabelTemplateLocal); err != nil {
		return nil
	}
	return &CachedTemplateInfo{Path: l[LabelTemplateCachePath]}
}

func (l Labels) ToGitTemplateInfo() *GitTemplateInfo {
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

func (l Labels) AddBoxSize(size boxModel.ResourceSize) Labels {
	return l.addLabel(LabelBoxSize, strings.ToLower(size.String()))
}

func (l Labels) ToBoxSize() (boxModel.ResourceSize, error) {
	if label, err := l.exist(LabelBoxSize); err != nil {
		return boxModel.Small, err
	} else {
		return boxModel.ExistResourceSize(label)
	}
}
