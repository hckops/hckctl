package old

import (
	"fmt"
	"github.com/hckops/hckctl/pkg/template"
	"github.com/pkg/errors"
)

type SourceTemplate interface {
	ReadTemplate() (*RawTemplate, error)
	ReadTemplates() ([]*TemplateValidated, error)
}

type SourceReader[T TemplateType] interface {
	Read() (*TemplateValue[T], error)
}

type CacheSourceOpts struct {
	cacheDir  string
	cacheName string
}

type LocalSource struct {
	cacheOpts *CacheSourceOpts
	path      string
}

func NewLocalSource(path string) *LocalSource {
	return &LocalSource{path: path}
}

func NewLocalCachedSource(path string, cacheDir string) *LocalSource {
	return &LocalSource{
		path: path,
		cacheOpts: &CacheSourceOpts{
			cacheDir:  cacheDir,
			cacheName: Local.String(),
		},
	}
}

func (src *LocalSource) ReadTemplate() (*RawTemplate, error) {
	return readRawTemplate(src.path)
}

func (src *LocalSource) ReadTemplates() ([]*TemplateValidated, error) {
	return readTemplates(src.path)
}

func (src *LocalSource) ReadBox() (*BoxInfo, error) {
	return readLocalBoxTemplate(src.cacheOpts, src.path)
}

func (src *LocalSource) ReadLab() (*LabInfo, error) {
	return readLocalLabTemplate(src.cacheOpts, src.path)
}

type RemoteSource struct {
	cacheOpts *CacheSourceOpts
	url       string
}

func NewRemoteSource(url string) *RemoteSource {
	return &RemoteSource{url: url}
}

func NewRemoteCachedSource(url string, cacheDir string) *RemoteSource {
	return &RemoteSource{
		url: url,
		cacheOpts: &CacheSourceOpts{
			cacheDir:  cacheDir,
			cacheName: Remote.String(),
		},
	}
}

func (src *RemoteSource) ReadTemplate() (*RawTemplate, error) {
	return nil, errors.New("not implemented")
}

func (src *RemoteSource) ReadTemplates() ([]*TemplateValidated, error) {
	return nil, errors.New("not implemented")
}

func (src *RemoteSource) ReadBox() (*BoxInfo, error) {
	return nil, errors.New("not implemented")
}

func (src *RemoteSource) ReadLab() (*LabInfo, error) {
	return nil, errors.New("not implemented")
}

type GitSource struct {
	opts *template.GitSourceOptions
	name string
}

func NewGitSource(opts *template.GitSourceOptions, name string) *GitSource {
	return &GitSource{opts, name}
}

func (src *GitSource) ReadTemplate() (*RawTemplate, error) {
	return readGitTemplate(src.opts, src.name)
}

func (src *GitSource) ReadTemplates() ([]*TemplateValidated, error) {
	wildcard := fmt.Sprintf("%s/**/*.{yml,yaml}", src.opts.CacheBaseDir)
	return readGitTemplates(src.opts, wildcard)
}
func (src *GitSource) ReadBox() (*BoxInfo, error) {
	var readGitTemplateInfo[box.BoxV1](src.opts, src.name)
	return readGitTemplateInfo[box.BoxV1](src.opts, src.name)
}
func (src *GitSource) ReadLab() (*LabInfo, error) {
	return readGitLabTemplate(src.opts, src.name)
}
