package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultShell(t *testing.T) {
	assert.Equal(t, "/bin/bash", DefaultShell(""))
	assert.Equal(t, "/bin/bash", DefaultShell("   "))
	assert.Equal(t, "/bin/bash", DefaultShell("\n\r\t"))

	assert.Equal(t, "foo", DefaultShell("foo"))
}

func TestToKebabCase(t *testing.T) {
	assert.Equal(t, "hckops-my-test_value-example", ToKebabCase("hckops/my-test_value$example"))
}
