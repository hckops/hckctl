package template

import (
	"github.com/pkg/errors"

	box "github.com/hckops/hckctl/pkg/box/model"
	lab "github.com/hckops/hckctl/pkg/lab/model"
	"github.com/hckops/hckctl/pkg/schema"
)

const (
	InvalidCommit = "INVALID_COMMIT"
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
	Commit   string
}

type LabTemplate struct {
	Template *lab.LabV1
	Path     string
	Commit   string
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
