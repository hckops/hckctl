package source

import (
	"path/filepath"

	"github.com/pkg/errors"

	"github.com/hckops/hckctl/pkg/template/model"
	"github.com/hckops/hckctl/pkg/template/schema"
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

	return &TemplateValue{kind, path, data}, nil
}

func readTemplates(wildcard string) ([]*TemplateValidated, error) {

	paths, err := filepath.Glob(wildcard)
	if err != nil {
		return nil, errors.Wrap(err, "invalid wildcard")
	}

	// validate all matching templates
	var results []*TemplateValidated
	for _, path := range paths {
		if value, err := readTemplate(path); err == nil {
			results = append(results, value.toValidated(true))
		} else {
			results = append(results, value.toValidated(false))
		}
	}
	return results, nil
}

func readBoxTemplate(path string) (*model.BoxV1, error) {
	if value, err := readTemplate(path); err != nil {
		return nil, err
	} else {
		return decodeBoxFromYaml(value.Data)
	}
}

func readLabTemplate(path string) (*model.LabV1, error) {
	if value, err := readTemplate(path); err != nil {
		return nil, err
	} else {
		return decodeLabFromYaml(value.Data)
	}
}
