package terminal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultShell(t *testing.T) {
	expected := []string{"/bin/bash"}
	assert.Equal(t, expected, DefaultShellCommand(""))
	assert.Equal(t, expected, DefaultShellCommand("   "))
	assert.Equal(t, expected, DefaultShellCommand("\n\r\t"))

	assert.Equal(t, []string{"foo"}, DefaultShellCommand("foo"))
}
