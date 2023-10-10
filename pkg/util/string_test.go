package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToLowerKebabCase(t *testing.T) {
	assert.Equal(t, "hckops-my-test_value-example", ToLowerKebabCase("  hCKops/my-tEst_value$ExamplE\t"))
}

func TestBase64(t *testing.T) {
	value := "hello world"
	encoded := Base64Encode(value)
	assert.Equal(t, "aGVsbG8gd29ybGQ=", encoded)

	decoded, ok := Base64Decode(encoded)
	assert.True(t, ok)
	assert.Equal(t, value, decoded)
}
