package template

import (
	"github.com/bmatcuk/doublestar/v4"
	"github.com/pkg/errors"

	box "github.com/hckops/hckctl/pkg/box/model"
	lab "github.com/hckops/hckctl/pkg/lab/model"
	"github.com/hckops/hckctl/pkg/schema"
	"github.com/hckops/hckctl/pkg/util"
)

func readTemplate(path string) (*TemplateValue, error) {
	data, err := util.ReadFile(path)
	if err != nil {
		return nil, errors.Wrapf(err, "template not found %s", path)
	}

	kind, err := schema.ValidateAll(data)
	if err != nil {
		return nil, errors.Wrapf(err, "invalid schema %s", data)
	}

	return &TemplateValue{kind, data}, nil
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
		if value, err := readTemplate(path); err != nil {
			results = append(results, (&TemplateValue{}).toValidated(path, false))
		} else {
			results = append(results, value.toValidated(path, true))
		}
	}
	return results, nil
}

func readBoxTemplate(path string) (*box.BoxV1, error) {
	if value, err := readTemplate(path); err != nil {
		return nil, err
	} else {
		return decodeBoxFromYaml(value.Data)
	}
}

func readLabTemplate(path string) (*lab.LabV1, error) {
	if value, err := readTemplate(path); err != nil {
		return nil, err
	} else {
		return decodeLabFromYaml(value.Data)
	}
}
