package template

import (
	"path/filepath"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/pkg/errors"

	box "github.com/hckops/hckctl/pkg/box/model"
	lab "github.com/hckops/hckctl/pkg/lab/model"
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
	return decodeFromYaml[T](raw)
}

func decodeFromYaml[T TemplateType](raw *RawTemplate) (*TemplateValue[T], error) {

	// https://stackoverflow.com/questions/71047848/how-to-assign-or-return-generic-t-that-is-constrained-by-union
	var templateType T
	switch typeRef := any(&templateType).(type) {
	case *box.BoxV1:
		template, err := decodeBoxFromYaml(raw.Data)
		if err != nil {
			return nil, err
		}
		*typeRef = *template

	case *lab.LabV1:
		template, err := decodeLabFromYaml(raw.Data)
		if err != nil {
			return nil, err
		}
		*typeRef = *template

	}

	return &TemplateValue[T]{Kind: raw.Kind, Data: templateType}, nil
}

// nil for generics
func none[T TemplateType]() T {
	return *new(T)
}

func readCachedTemplateInfo[T TemplateType](cacheOpts *CacheSourceOpts, path string, sourceType SourceType) (*TemplateInfo[T], error) {

	value, err := readTemplate[T](path)
	if err != nil {
		return nil, err
	}

	// TODO gen name
	// TODO save value.Data to cachedPath
	cachedPath := filepath.Join(cacheOpts.cacheDir, cacheOpts.cacheName, "TODO_GEN_NAME")

	return &TemplateInfo[T]{
		Value:      value,
		SourceType: sourceType,
		Cached:     true,
		Path:       cachedPath,
		Revision:   sourceType.String(),
	}, nil
}
