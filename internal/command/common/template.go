package common

import (
	"github.com/hckops/hckctl/pkg/template"
)

func NewGitSourceOptions(cacheDir string, revision string) *template.GitSourceOptions {
	return &template.GitSourceOptions{
		CacheBaseDir:    cacheDir,
		RepositoryUrl:   TemplateSourceUrl,
		DefaultRevision: TemplateSourceRevision,
		Revision:        revision,
		AllowOffline:    true,
	}
}
