package template

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/hckops/hckctl/pkg/template/schema"
	"github.com/hckops/hckctl/pkg/util"
)

type RequestLocalTemplate struct {
	Path   string
	Format string
}

type ResponseLocalTemplate struct {
	Kind  schema.SchemaKind
	Value string
}

func LoadLocalTemplate(request *RequestLocalTemplate) (*ResponseLocalTemplate, error) {
	value, err := util.ReadFile(request.Path)
	if err != nil {
		return nil, errors.Wrapf(err, "local template not found %s", request.Path)
	}

	kind, err := schema.ValidateAll(value)
	if err != nil {
		return nil, errors.Wrapf(err, "invalid schema %s", value)
	}

	switch request.Format {
	case YamlFormat.String():
		return &ResponseLocalTemplate{kind, value}, nil
	case JsonFormat.String():
		if jsonValue, err := ConvertFromYamlToJson(kind, value); err != nil {
			return nil, errors.Wrap(err, "conversion from yaml to json failed")
		} else {
			// adds newline only for json
			return &ResponseLocalTemplate{kind, fmt.Sprintf("%s\n", jsonValue)}, nil
		}
	default:
		return nil, fmt.Errorf("invalid Format %s", request.Format)
	}
}
