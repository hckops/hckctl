package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToLowerKebabCase(t *testing.T) {
	assert.Equal(t, "hckops-my-test_value-example", ToLowerKebabCase("  hCKops/my-tEst_value$ExamplE\t"))
}
