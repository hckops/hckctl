package schema

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/santhosh-tekuri/jsonschema/v5"
	"gopkg.in/yaml.v3"
)

type schemaValidation struct {
	validationFunc func(string) error
	kind           SchemaKind
}

func ValidateAll(data string) (SchemaKind, error) {
	schemaValidations := []schemaValidation{
		{ValidateBoxV1, KindBoxV1},
		{ValidateLabV1, KindLabV1},
		{ValidateTaskV1, KindTaskV1},
		{ValidateDumpV1, KindDumpV1},
	}
	var validationErrors []error
	for _, sv := range schemaValidations {
		if err := sv.validationFunc(data); err == nil {
			// found valid schema
			return sv.kind, nil
		} else {
			validationErrors = append(validationErrors, errors.Wrapf(err, "failed to match schema %s", sv.kind))
		}
	}
	return -1, fmt.Errorf("unable to find matching schema %v", validationErrors)
}

func ValidateBoxV1(data string) error {
	return validateSchema("box-v1.json", boxV1Schema, data)
}

func ValidateLabV1(data string) error {
	return validateSchema("lab-v1.json", labV1Schema, data)
}

func ValidateTaskV1(data string) error {
	return validateSchema("task-v1.json", taskV1Schema, data)
}

func ValidateDumpV1(data string) error {
	return validateSchema("dump-v1.json", dumpV1Schema, data)
}

// returns nil if valid
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
