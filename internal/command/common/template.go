package common

import (
	"fmt"
	"strings"

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

func PrettyPath(cacheDir string, path string) string {
	// remove prefix and suffix
	return strings.NewReplacer(
		fmt.Sprintf("%s/", cacheDir), "",
		".yml", "",
		".yaml", "",
	).Replace(path)
}

func PrettyName[T template.TemplateType](info *template.TemplateInfo[T], cacheDir string, name string) string {
	if info.SourceType == template.Git {
		return fmt.Sprintf("%s@%s", PrettyPath(cacheDir, info.Path), info.Revision[:7])
	}
	return fmt.Sprintf("%s@%s", strings.ToLower(name), info.Revision)
}
