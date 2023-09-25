package flag

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/hckops/hckctl/pkg/common/model"
)

func TestValidateParametersFlag(t *testing.T) {

	emptyParams, emptyErr := ValidateParametersFlag([]string{})
	assert.Equal(t, model.Parameters{}, emptyParams)
	assert.Nil(t, emptyErr)

	inputs := []string{"key1=value1", " key2 \t=\nvalue2 "}
	expected := model.Parameters{"key1": "value1", "key2": "value2"}
	validParams, validErr := ValidateParametersFlag(inputs)
	assert.Equal(t, expected, validParams)
	assert.Nil(t, validErr)

	invalidFormatErrString := "invalid parameter format [invalid], expected KEY=VALUE"
	invalidFormatParams, invalidFormatErr := ValidateParametersFlag([]string{"invalid"})
	assert.Nil(t, invalidFormatParams)
	assert.EqualError(t, invalidFormatErr, invalidFormatErrString)

	invalidKeyErrString := "invalid parameter key format [=value], expected KEY=VALUE"
	invalidKeyParams, invalidKeyErr := ValidateParametersFlag([]string{"=value"})
	assert.Nil(t, invalidKeyParams)
	assert.EqualError(t, invalidKeyErr, invalidKeyErrString)

	invalidValueErrString := "invalid parameter value format [key=], expected KEY=VALUE"
	invalidValueParams, invalidValueErr := ValidateParametersFlag([]string{"key="})
	assert.Nil(t, invalidValueParams)
	assert.EqualError(t, invalidValueErr, invalidValueErrString)
}
