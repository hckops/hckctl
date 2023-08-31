package model

import (
	"fmt"
	"path/filepath"
	"strings"

	"golang.org/x/exp/maps"
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
)

func (l Labels) AddLabel(key string, value string) Labels {
	return l.AddLabels(Labels{key: value})
}

func (l Labels) AddLabels(labels Labels) Labels {
	// merge labels
	maps.Copy(labels, l)
	return labels
}

func (l Labels) Exist(name string) (string, error) {
	if label, ok := l[name]; !ok {
		return "", fmt.Errorf("label %s not found", name)
	} else {
		return label, nil
	}
}

func (l Labels) AddDefaultLocal() Labels {
	return l.AddLabel(LabelTemplateLocal, "true")
}

func (l Labels) AddDefaultGit(url, revision, dir string) Labels {
	return l.AddLabels(Labels{
		LabelTemplateGit:         "true",
		LabelTemplateGitUrl:      url,
		LabelTemplateGitRevision: revision,
		LabelTemplateGitDir:      dir,
	})
}

func (l Labels) AddLocal(path string) Labels {
	if _, err := l.Exist(LabelTemplateLocal); err != nil {
		return l
	}
	return l.AddLabel(LabelTemplateCachePath, path)
}

func (l Labels) AddGit(path string, commit string) Labels {
	if _, err := l.Exist(LabelTemplateGit); err != nil {
		return l
	}

	templatePath := strings.SplitAfter(path, l[LabelTemplateGitDir])
	// unsafe assume valid path
	name := strings.TrimSuffix(strings.TrimPrefix(templatePath[1], "/"), filepath.Ext(path))

	return l.AddLabels(Labels{
		LabelTemplateGitCommit: commit,
		LabelTemplateGitName:   name,
		LabelTemplateCachePath: path,
	})
}

func (l Labels) ToCachedTemplateInfo() *CachedTemplateInfo {
	if _, err := l.Exist(LabelTemplateLocal); err != nil {
		return nil
	}
	return &CachedTemplateInfo{Path: l[LabelTemplateCachePath]}
}

func (l Labels) ToGitTemplateInfo() *GitTemplateInfo {
	if _, err := l.Exist(LabelTemplateGit); err != nil {
		return nil
	}

	return &GitTemplateInfo{
		Url:      l[LabelTemplateGitUrl],
		Revision: l[LabelTemplateGitRevision],
		Commit:   l[LabelTemplateGitCommit],
		Name:     l[LabelTemplateGitName],
	}
}
