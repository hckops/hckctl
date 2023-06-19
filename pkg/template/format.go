package template

import (
	"fmt"

	"gopkg.in/yaml.v3"

	box "github.com/hckops/hckctl/pkg/box/model"
	lab "github.com/hckops/hckctl/pkg/lab/model"
	"github.com/hckops/hckctl/pkg/template/schema"
	"github.com/hckops/hckctl/pkg/util"
)

// formats, strips out comments, quotes, etc. and prints optional empty fields
func convertFromYamlToYaml(kind schema.SchemaKind, value string) (string, error) {
	switch kind {
	case schema.KindBoxV1:
		if model, err := decodeBoxFromYaml(value); err != nil {
			return "", err
		} else {
			return util.EncodeYaml(model)
		}
	case schema.KindLabV1:
		if model, err := decodeLabFromYaml(value); err != nil {
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
		if model, err := decodeBoxFromYaml(value); err != nil {
			return "", err
		} else {
			return util.EncodeJsonIndent(model)
		}
	case schema.KindLabV1:
		if model, err := decodeLabFromYaml(value); err != nil {
			return "", err
		} else {
			return util.EncodeJsonIndent(model)
		}
	default:
		return "", fmt.Errorf("invalid kind: %v", kind)
	}
}

func decodeBoxFromYaml(value string) (*box.BoxV1, error) {
	var model box.BoxV1
	if err := yaml.Unmarshal([]byte(value), &model); err != nil {
		return nil, fmt.Errorf("box decoder error: %v", err)
	}
	return &model, nil
}

func decodeLabFromYaml(value string) (*lab.LabV1, error) {
	var model lab.LabV1
	if err := yaml.Unmarshal([]byte(value), &model); err != nil {
		return nil, fmt.Errorf("lab decoder error: %v", err)
	}
	return &model, nil
}
