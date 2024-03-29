package template

import (
	"fmt"
	"reflect"

	"gopkg.in/yaml.v3"

	box "github.com/hckops/hckctl/pkg/box/model"
	lab "github.com/hckops/hckctl/pkg/lab/model"
	"github.com/hckops/hckctl/pkg/schema"
	task "github.com/hckops/hckctl/pkg/task/model"
	"github.com/hckops/hckctl/pkg/util"
)

// formats, strips out comments, quotes, etc. and prints optional empty fields
func convertFromYamlToYaml(kind schema.SchemaKind, value string) (string, error) {
	switch kind {
	case schema.KindBoxV1:
		if model, err := decodeFromYaml[box.BoxV1](value); err != nil {
			return "", err
		} else {
			return util.EncodeYaml(model)
		}
	case schema.KindLabV1:
		if model, err := decodeFromYaml[lab.LabV1](value); err != nil {
			return "", err
		} else {
			return util.EncodeYaml(model)
		}
	case schema.KindTaskV1:
		if model, err := decodeFromYaml[task.TaskV1](value); err != nil {
			return "", err
		} else {
			return util.EncodeYaml(model)
		}
	case schema.KindDumpV1:
		if model, err := decodeFromYaml[lab.DumpV1](value); err != nil {
			return "", err
		} else {
			return util.EncodeYaml(model)
		}
	default:
		return "", fmt.Errorf("invalid kind: %v", kind)
	}
}

func convertFromYamlToJson(kind schema.SchemaKind, value string) (string, error) {
	switch kind {
	case schema.KindBoxV1:
		if model, err := decodeFromYaml[box.BoxV1](value); err != nil {
			return "", err
		} else {
			return util.EncodeJsonIndent(model)
		}
	case schema.KindLabV1:
		if model, err := decodeFromYaml[lab.LabV1](value); err != nil {
			return "", err
		} else {
			return util.EncodeJsonIndent(model)
		}
	case schema.KindTaskV1:
		if model, err := decodeFromYaml[task.TaskV1](value); err != nil {
			return "", err
		} else {
			return util.EncodeJsonIndent(model)
		}
	case schema.KindDumpV1:
		if model, err := decodeFromYaml[lab.DumpV1](value); err != nil {
			return "", err
		} else {
			return util.EncodeJsonIndent(model)
		}
	default:
		return "", fmt.Errorf("invalid kind: %v", kind)
	}
}

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

	case *task.TaskV1:
		var model task.TaskV1
		if err := yaml.Unmarshal([]byte(value), &model); err != nil {
			return none[T](), fmt.Errorf("task decoder error: %v", err)
		}
		*typeRef = model

	case *lab.DumpV1:
		var model lab.DumpV1
		if err := yaml.Unmarshal([]byte(value), &model); err != nil {
			return none[T](), fmt.Errorf("dump decoder error: %v", err)
		}
		*typeRef = model

	default:
		return templateType, fmt.Errorf("unable to decode yaml invalid schema %v", reflect.TypeOf(templateType))

	}
	return templateType, nil
}

// nil for generics
func none[T TemplateType]() T {
	return *new(T)
}
