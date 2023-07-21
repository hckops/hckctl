package template

import (
	"fmt"

	box "github.com/hckops/hckctl/pkg/box/model"
	lab "github.com/hckops/hckctl/pkg/lab/model"
	"github.com/hckops/hckctl/pkg/schema"
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
	default:
		return "", fmt.Errorf("invalid kind: %v", kind)
	}
}
