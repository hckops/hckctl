package template

import (
	"fmt"
)

// TODO add RemoteSource http

type SourceTemplate interface {
	ReadTemplate() (*TemplateValue, error)
	ReadTemplates() ([]*TemplateValidated, error)
	ReadBox() (*BoxTemplate, error)
	ReadLab() (*LabTemplate, error)
}

type LocalSource struct {
	path string
}

func NewLocalSource(path string) *LocalSource {
	return &LocalSource{path}
}

func (src *LocalSource) ReadTemplate() (*TemplateValue, error) {
	return readTemplate(src.path)
}

func (src *LocalSource) ReadTemplates() ([]*TemplateValidated, error) {
	return readTemplates(src.path)
}

func (src *LocalSource) ReadBox() (*BoxTemplate, error) {
	return readBoxTemplate(src.path, InvalidCommit)
}

func (src *LocalSource) ReadLab() (*LabTemplate, error) {
	return readLabTemplate(src.path, InvalidCommit)
}

type GitSource struct {
	opts *GitSourceOptions
	name string
}

func NewGitSource(opts *GitSourceOptions, name string) *GitSource {
	return &GitSource{opts, name}
}

func (src *GitSource) ReadTemplate() (*TemplateValue, error) {
	return readGitTemplate(src.opts, src.name)
}

func (src *GitSource) ReadTemplates() ([]*TemplateValidated, error) {
	wildcard := fmt.Sprintf("%s/**/*.{yml,yaml}", src.opts.CacheBaseDir)
	return readGitTemplates(src.opts, wildcard)
}
func (src *GitSource) ReadBox() (*BoxTemplate, error) {
	return readGitBoxTemplate(src.opts, src.name)
}
func (src *GitSource) ReadLab() (*LabTemplate, error) {
	return readGitLabTemplate(src.opts, src.name)
}
