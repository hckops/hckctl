package template

import (
	"fmt"

	"gopkg.in/yaml.v3"

	box "github.com/hckops/hckctl/pkg/box/model"
	lab "github.com/hckops/hckctl/pkg/lab/model"
)

func decodeFromYaml[T TemplateType](value string) (T, error) {

	// https://stackoverflow.com/questions/71047848/how-to-assign-or-return-generic-t-that-is-constrained-by-union
	var templateType T
	switch typeRef := any(&templateType).(type) {
	case *string:
		*typeRef = value

	case *box.BoxV1:
		var model box.BoxV1
		if err := yaml.Unmarshal([]byte(value), &model); err != nil {
			return none[T](), fmt.Errorf("box decoder error: %v", err)
		}
		*typeRef = model

	case *lab.LabV1:
		var model lab.LabV1
		if err := yaml.Unmarshal([]byte(value), &model); err != nil {
			return none[T](), fmt.Errorf("lab decoder error: %v", err)
		}
		*typeRef = model

	}
	return templateType, nil
}

// nil for generics
func none[T TemplateType]() T {
	return *new(T)
}
