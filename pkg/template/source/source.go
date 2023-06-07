package source

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/hckops/hckctl/pkg/template/model"
	"github.com/hckops/hckctl/pkg/template/schema"
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

type TemplateSource interface {
	Read() (*TemplateValue, error)
	ReadAll() ([]*TemplateValidated, error)
	ReadBox() (*model.BoxV1, error)
	ReadLab() (*model.LabV1, error)
}

type LocalSource struct {
	path string
}

func NewLocalSource(path string) *LocalSource {
	return &LocalSource{path}
}

func (src *LocalSource) Read() (*TemplateValue, error) {
	return readTemplate(src.path)
}

func (src *LocalSource) ReadAll() ([]*TemplateValidated, error) {
	return readTemplates(src.path)
}
func (src *LocalSource) ReadBox() (*model.BoxV1, error) {
	return readBoxTemplate(src.path)
}
func (src *LocalSource) ReadLab() (*model.LabV1, error) {
	return readLabTemplate(src.path)
}

type RemoteSource struct {
	opts *RevisionOpts
	name string
}

func NewRemoteSource(opts *RevisionOpts, name string) *RemoteSource {
	return &RemoteSource{opts, name}
}

func (src *RemoteSource) Read() (*TemplateValue, error) {
	return readRemoteTemplate(src.opts, src.name)
}

func (src *RemoteSource) ReadAll() ([]*TemplateValidated, error) {
	wildcard := fmt.Sprintf("%s/**/*.{yml,yaml}", src.opts.SourceCacheDir)
	return readRemoteTemplates(src.opts, wildcard)
}
func (src *RemoteSource) ReadBox() (*model.BoxV1, error) {
	return readRemoteBoxTemplate(src.opts, src.name)
}
func (src *RemoteSource) ReadLab() (*model.LabV1, error) {
	return readRemoteLabTemplate(src.opts, src.name)
}
