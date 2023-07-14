package template

import (
	"fmt"

	"github.com/pkg/errors"

	box "github.com/hckops/hckctl/pkg/box/model"
	lab "github.com/hckops/hckctl/pkg/lab/model"
	"github.com/hckops/hckctl/pkg/schema"
)

type TemplateValue struct {
	Kind schema.SchemaKind
	Data string
}

type TemplateValidated struct {
	Value   *TemplateValue
	Path    string
	IsValid bool
}

type BoxTemplate struct {
	Template *box.BoxV1
	Path     string
}

type LabTemplate struct {
	Template *lab.LabV1
	Path     string
}

func (t *TemplateValue) ToYaml() (*TemplateValue, error) {
	if yamlValue, err := convertFromYamlToYaml(t.Kind, t.Data); err != nil {
		return nil, errors.Wrap(err, "conversion to yaml failed")
	} else {
		t.Data = yamlValue
		return t, nil
	}
}

func (t *TemplateValue) ToJson() (*TemplateValue, error) {
	if jsonValue, err := convertFromYamlToJson(t.Kind, t.Data); err != nil {
		return nil, errors.Wrap(err, "conversion to json failed")
	} else {
		t.Data = jsonValue
		return t, nil
	}
}

func (t *TemplateValue) toValidated(path string, isValid bool) *TemplateValidated {
	return &TemplateValidated{t, path, isValid}
}

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
	return readBoxTemplate(src.path)
}
func (src *LocalSource) ReadLab() (*LabTemplate, error) {
	return readLabTemplate(src.path)
}

type GitSource struct {
	opts *SourceOptions
	name string
}

func NewGitSource(opts *SourceOptions, name string) *GitSource {
	return &GitSource{opts, name}
}

func (src *GitSource) ReadTemplate() (*TemplateValue, error) {
	return readGitTemplate(src.opts, src.name)
}

func (src *GitSource) ReadTemplates() ([]*TemplateValidated, error) {
	wildcard := fmt.Sprintf("%s/**/*.{yml,yaml}", src.opts.SourceCacheDir)
	return readGitTemplates(src.opts, wildcard)
}
func (src *GitSource) ReadBox() (*BoxTemplate, error) {
	return readGitBoxTemplate(src.opts, src.name)
}
func (src *GitSource) ReadLab() (*LabTemplate, error) {
	return readGitLabTemplate(src.opts, src.name)
}
