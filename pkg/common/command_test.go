package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewCommandCreateBox(t *testing.T) {
	assert.Equal(t, "hck-box-create::my-group/my-name::main", NewCommandCreateBox("my-group/my-name", "main"))
}

func TestNewCommandExecBox(t *testing.T) {
	assert.Equal(t, "hck-box-exec::my-group/my-name::main::box-id", NewCommandExecBox("my-group/my-name", "main", "box-id"))
}

func TestNewCommandListBox(t *testing.T) {
	assert.Equal(t, "hck-box-list", NewCommandListBox())
}

func TestNewCommandDeleteBox(t *testing.T) {
	assert.Equal(t, "hck-box-delete::my-group/my-name::main::box-id", NewCommandDeleteBox("my-group/my-name", "main", "box-id"))
}
