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
