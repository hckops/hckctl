package old

import (
	"github.com/bmatcuk/doublestar/v4"
	box "github.com/hckops/hckctl/pkg/box/model"
	lab "github.com/hckops/hckctl/pkg/lab/model"
	"github.com/hckops/hckctl/pkg/schema"
	"github.com/hckops/hckctl/pkg/template"
	"github.com/hckops/hckctl/pkg/util"
	"github.com/pkg/errors"
)

func readRawTemplate(path string) (*RawTemplate, error) {
	data, err := util.ReadFile(path)
	if err != nil {
		return nil, errors.Wrapf(err, "template not found %s", path)
	}

	kind, err := schema.ValidateAll(data)
	if err != nil {
		return nil, errors.Wrapf(err, "invalid schema %s", data)
	}

	return &RawTemplate{kind, data}, nil
}

func readTemplates(wildcard string) ([]*TemplateValidated, error) {

	// https://github.com/golang/go/issues/11862
	paths, err := doublestar.FilepathGlob(wildcard,
		doublestar.WithFailOnPatternNotExist(), doublestar.WithFilesOnly(), doublestar.WithNoFollow())
	if err != nil {
		return nil, errors.Wrap(err, "invalid wildcard")
	}

	// validate all matching templates
	var results []*TemplateValidated
	for _, path := range paths {
		if value, err := readRawTemplate(path); err != nil {
			results = append(results, (&RawTemplate{}).toValidated(path, false))
		} else {
			results = append(results, value.toValidated(path, true))
		}
	}
	return results, nil
}

func decodeFromYaml(raw *RawTemplate) (any, error) {
	switch raw.Kind {
	case schema.KindBoxV1:
		return template.decodeBoxFromYaml(raw.Data)
	case schema.KindLabV1:
		return template.decodeLabFromYaml(raw.Data)
	}
}

// https://stackoverflow.com/questions/71047848/how-to-assign-or-return-generic-t-that-is-constrained-by-union
func readTemplate[T TemplateType](path string) (*TemplateValue[T], error) {
	raw, err := readRawTemplate(path)
	if err != nil {
		return nil, err
	}

	var templateType T
	switch typeRef := any(&templateType).(type) {
	case *box.BoxV1:
		template, err := template.decodeBoxFromYaml(raw.Data)
		if err != nil {
			return nil, err
		}
		*typeRef = *template
		return newBoxTemplate(template), nil

	case *lab.LabV1:
		template, err := template.decodeLabFromYaml(raw.Data)
		if err != nil {
			return nil, err
		}
		*typeRef = *template
	}

	return templateType, nil
}

func readBoxTemplate(path string) (*BoxTemplate, error) {
	raw, err := readRawTemplate(path)
	if err != nil {
		return nil, err
	}

	template, err := template.decodeBoxFromYaml(raw.Data)
	if err != nil {
		return nil, err
	}

	return newBoxTemplate(template), nil
}

func readLabTemplate(path string) (*LabTemplate, error) {
	value, err := readRawTemplate(path)
	if err != nil {
		return nil, err
	}

	template, err := template.decodeLabFromYaml(value.Data)
	if err != nil {
		return nil, err
	}

	return newLabTemplate(template), nil
}
