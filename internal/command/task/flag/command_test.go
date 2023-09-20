package flag

import (
	"testing"

	"github.com/stretchr/testify/assert"

	commonModel "github.com/hckops/hckctl/pkg/common/model"
)

func TestValidateCommandInputsFlag(t *testing.T) {

	emptyParams, emptyErr := ValidateCommandInputsFlag([]string{})
	assert.Equal(t, commonModel.Parameters{}, emptyParams)
	assert.Nil(t, emptyErr)

	inputs := []string{"key1=value1", " key2 \t=\nvalue2 "}
	expected := commonModel.Parameters{"key1": "value1", "key2": "value2"}
	validParams, validErr := ValidateCommandInputsFlag(inputs)
	assert.Equal(t, expected, validParams)
	assert.Nil(t, validErr)

	invalidFormatErrString := "invalid input format [invalid], expected KEY=VALUE"
	invalidFormatParams, invalidFormatErr := ValidateCommandInputsFlag([]string{"invalid"})
	assert.Nil(t, invalidFormatParams)
	assert.EqualError(t, invalidFormatErr, invalidFormatErrString)

	invalidKeyErrString := "invalid input key format [=value], expected KEY=VALUE"
	invalidKeyParams, invalidKeyErr := ValidateCommandInputsFlag([]string{"=value"})
	assert.Nil(t, invalidKeyParams)
	assert.EqualError(t, invalidKeyErr, invalidKeyErrString)

	invalidValueErrString := "invalid input value format [key=], expected KEY=VALUE"
	invalidValueParams, invalidValueErr := ValidateCommandInputsFlag([]string{"key="})
	assert.Nil(t, invalidValueParams)
	assert.EqualError(t, invalidValueErr, invalidValueErrString)
}
