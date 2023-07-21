package template

import (
	"github.com/pkg/errors"

	box "github.com/hckops/hckctl/pkg/box/model"
	lab "github.com/hckops/hckctl/pkg/lab/model"
	"github.com/hckops/hckctl/pkg/schema"
)

type SourceType uint

const (
	Local SourceType = iota
	Remote
	Git
)

var sources = []string{"local", "remote", "git"}

func (s SourceType) String() string {
	return sources[s]
}

type TemplateType interface {
	string | box.BoxV1 | lab.LabV1
}

type TemplateValue[T TemplateType] struct {
	Kind schema.SchemaKind
	Data T // string or actual model
}

type TemplateValidated struct {
	Value   *RawTemplate
	Path    string
	IsValid bool
}

type TemplateInfo[T TemplateType] struct {
	Value      *TemplateValue[T]
	SourceType SourceType
	Cached     bool
	Path       string // absolute path cached or resolved git path
	Revision   string // local/remote always none, git commit
}

// fixes receiver types
type RawTemplate TemplateValue[string]

func (t *RawTemplate) ToYaml() (*RawTemplate, error) {
	if yamlValue, err := convertFromYamlToYaml(t.Kind, t.Data); err != nil {
		return nil, errors.Wrap(err, "conversion to yaml failed")
	} else {
		t.Data = yamlValue
		return t, nil
	}
}

func (t *RawTemplate) ToJson() (*RawTemplate, error) {
	if jsonValue, err := convertFromYamlToJson(t.Kind, t.Data); err != nil {
		return nil, errors.Wrap(err, "conversion to json failed")
	} else {
		t.Data = jsonValue
		return t, nil
	}
}
