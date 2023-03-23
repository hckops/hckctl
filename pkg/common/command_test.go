package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewNewCommandCreateBox(t *testing.T) {
	assert.Equal(t, "my-group/my-name::main", NewCommandCreateBox("my-group/my-name", "main"))
}
