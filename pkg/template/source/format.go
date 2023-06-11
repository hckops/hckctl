package source

import (
	"fmt"

	"gopkg.in/yaml.v3"

	"github.com/hckops/hckctl/pkg/template/model"
	"github.com/hckops/hckctl/pkg/template/schema"
	"github.com/hckops/hckctl/pkg/util"
)

// get rid of comments and print optional empty fields
func convertFromYamlToYaml(kind schema.SchemaKind, value string) (string, error) {
	switch kind {
	case schema.KindBoxV1:
		if box, err := decodeBoxFromYaml(value); err != nil {
			return "", err
		} else {
			return util.EncodeYaml(box)
		}
	case schema.KindLabV1:
		if lab, err := decodeLabFromYaml(value); err != nil {
			return "", err
		} else {
			return util.EncodeYaml(lab)
		}
	default:
		return "", fmt.Errorf("invalid kind: %v", kind)
	}
}

func convertFromYamlToJson(kind schema.SchemaKind, value string) (string, error) {
	switch kind {
	case schema.KindBoxV1:
		if box, err := decodeBoxFromYaml(value); err != nil {
			return "", err
		} else {
			return util.EncodeJsonIndent(box)
		}
	case schema.KindLabV1:
		if lab, err := decodeLabFromYaml(value); err != nil {
			return "", err
		} else {
			return util.EncodeJsonIndent(lab)
		}
	default:
		return "", fmt.Errorf("invalid kind: %v", kind)
	}
}

func decodeBoxFromYaml(value string) (*model.BoxV1, error) {
	var box model.BoxV1
	if err := yaml.Unmarshal([]byte(value), &box); err != nil {
		return nil, fmt.Errorf("box decoder error: %v", err)
	}
	return &box, nil
}

func decodeLabFromYaml(value string) (*model.LabV1, error) {
	var lab model.LabV1
	if err := yaml.Unmarshal([]byte(value), &lab); err != nil {
		return nil, fmt.Errorf("lab decoder error: %v", err)
	}
	return &lab, nil
}
