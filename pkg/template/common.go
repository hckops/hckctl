package template

import (
	"github.com/bmatcuk/doublestar/v4"
	"github.com/pkg/errors"

	"github.com/hckops/hckctl/pkg/schema"
	"github.com/hckops/hckctl/pkg/util"
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

func readTemplate[T TemplateType](path string) (*TemplateValue[T], error) {
	raw, err := readRawTemplate(path)
	if err != nil {
		return nil, err
	}

	value, err := decodeFromYaml[T](raw.Data)
	if err != nil {
		return nil, err
	}

	return &TemplateValue[T]{
		Kind: raw.Kind,
		Data: value,
	}, nil
}

func readTemplateInfo[T TemplateType](sourceType SourceType, path string, revision string) (*TemplateInfo[T], error) {

	value, err := readTemplate[T](path)
	if err != nil {
		return nil, err
	}

	return &TemplateInfo[T]{
		Value:      value,
		SourceType: sourceType,
		Path:       path,
		Revision:   revision,
	}, nil
}
