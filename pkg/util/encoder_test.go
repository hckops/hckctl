package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToKebabCase(t *testing.T) {
	assert.Equal(t, "hckops-my-test_value-example", ToKebabCase("hckops/my-test_value$example"))
}
