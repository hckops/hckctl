package template

import (
	"fmt"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/pkg/errors"
	"path/filepath"

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

func readBoxTemplate(path, hash string) (*BoxTemplate, error) {
	value, err := readTemplate(path)
	if err != nil {
		return nil, err
	}

	template, err := decodeBoxFromYaml(value.Data)
	if err != nil {
		return nil, err
	}

	absolutePath, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("unable to resolve absolute box path %s", absolutePath)
	}

	return &BoxTemplate{
		Template: template,
		Path:     absolutePath,
		Commit:   hash,
	}, nil
}

func readLabTemplate(path, hash string) (*LabTemplate, error) {
	value, err := readTemplate(path)
	if err != nil {
		return nil, err
	}

	template, err := decodeLabFromYaml(value.Data)
	if err != nil {
		return nil, err
	}

	absolutePath, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("unable to resolve absolute lab path %s", absolutePath)
	}

	return &LabTemplate{
		Template: template,
		Path:     absolutePath,
		Commit:   hash,
	}, nil
}
