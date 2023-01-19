package template

import (
	_ "embed"
	"fmt"

	"github.com/santhosh-tekuri/jsonschema/v5"
	"gopkg.in/yaml.v3"

	"github.com/hckops/hckctl/internal/model"
)

//go:embed schema/box-v1.json
var schemaBoxV1 string

// TODO currently supports "box/v1" only: iterate over all validators/versions
// see objectForKind e.g. yaml -> json https://github.com/argoproj/argo-workflows/blob/master/workflow/common/parse.go
func ValidateAllSchema(data string) error {
	return ValidateBoxV1(data)
}

func ParseValidBoxV1(data string) (*model.BoxV1, error) {
	if err := ValidateBoxV1(data); err != nil {
		return nil, err
	}

	box, err := ParseBoxV1(data)
	if err != nil {
		return nil, err
	}

	return box, nil
}

// returns nil if valid
func ValidateBoxV1(data string) error {
	if err := validateSchema("box-v1.json", schemaBoxV1, data); err != nil {
		return err
	}
	return nil
}

func ParseBoxV1(data string) (*model.BoxV1, error) {
	// TODO generics ?!
	var box model.BoxV1
	if err := yaml.Unmarshal([]byte(data), &box); err != nil {
		return nil, fmt.Errorf("decode error: %v", err)
	}
	return &box, nil
}

func validateSchema(schemaName string, schemaValue string, data string) error {
	schema, err := jsonschema.CompileString(schemaName, schemaValue)
	if err != nil {
		return fmt.Errorf("schema error: %v", err)
	}

	var model interface{}
	if err := yaml.Unmarshal([]byte(data), &model); err != nil {
		return fmt.Errorf("yaml error: %v", err)
	}

	if err := schema.Validate(model); err != nil {
		return fmt.Errorf("validation error: %v", err)
	}

	return nil
}
