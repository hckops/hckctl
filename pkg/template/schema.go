package template

import (
	_ "embed"
	"encoding/json"
	"fmt"

	"github.com/santhosh-tekuri/jsonschema/v5"
	"gopkg.in/yaml.v3"
)

//go:embed schema/box-v1.json
var schemaBoxV1 string

func ParseBoxV1(data string) (*BoxV1, error) {
	if err := ValidateBoxV1(data); err != nil {
		return nil, fmt.Errorf("validation error")
	}

	var box *BoxV1
	if err := yaml.Unmarshal([]byte(data), box); err != nil {
		return nil, fmt.Errorf("unmarshal error")
	}

	return box, nil
}

func ValidateBoxV1(data string) error {
	if err := validateSchema("box-v1.json", schemaBoxV1, data); err != nil {
		return fmt.Errorf("validation error")
	}

	return nil
}

func validateSchema(schemaName string, schemaValue string, data string) error {
	schema, err := jsonschema.CompileString(schemaName, schemaValue)
	if err != nil {
		return fmt.Errorf("schema error: %s", schemaName)
	}

	var model interface{}
	if err := yaml.Unmarshal([]byte(data), &model); err != nil {
		return fmt.Errorf("yaml error: %s", schemaName)
	}

	if validationError := schema.Validate(model).(*jsonschema.ValidationError); validationError != nil {
		errorDetails, _ := json.MarshalIndent(validationError.DetailedOutput(), "", "  ")
		return fmt.Errorf(string(errorDetails))
	}

	return nil
}
