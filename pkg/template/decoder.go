package template

import (
	"fmt"

	"gopkg.in/yaml.v3"

	"github.com/hckops/hckctl/pkg/template/schema"
	"github.com/hckops/hckctl/pkg/util"
)

type Format string

const (
	YamlFormat Format = "yaml"
	JsonFormat Format = "json"
)

func (f Format) String() string {
	return string(f)
}

func ConvertFromYamlToJson(kind schema.SchemaKind, value string) (string, error) {
	switch kind {
	case schema.KindBoxV1:
		if box, err := DecodeBoxFromYaml(value); err != nil {
			return "", err
		} else {
			return util.EncodeJsonIndent(box)
		}
	case schema.KindLabV1:
		if lab, err := DecodeLabFromYaml(value); err != nil {
			return "", err
		} else {
			return util.EncodeJsonIndent(lab)
		}
	default:
		return "", fmt.Errorf("invalid kind: %v", kind)
	}
}

func DecodeBoxFromYaml(value string) (*BoxV1, error) {
	var box BoxV1
	if err := yaml.Unmarshal([]byte(value), &box); err != nil {
		return nil, fmt.Errorf("box decoder error: %v", err)
	}
	return &box, nil
}

func DecodeLabFromYaml(value string) (*LabV1, error) {
	var lab LabV1
	if err := yaml.Unmarshal([]byte(value), &lab); err != nil {
		return nil, fmt.Errorf("lab decoder error: %v", err)
	}
	return &lab, nil
}
