package cloud

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCommandConstants(t *testing.T) {
	assert.Equal(t, CommandRequestType, "hck-v1")
	assert.Equal(t, CommandResponseError, "error")
	assert.Equal(t, CommandDelimiter, "::")
}

func TestCommandFromString(t *testing.T) {
	command, err := FromString("hck-box-create")
	assert.NoError(t, err)
	assert.Equal(t, CommandBoxCreate, command)
}

func TestCommandFromStringError(t *testing.T) {
	command, err := FromString("todo")
	assert.EqualError(t, err, "invalid command: todo")
	assert.Equal(t, Command(-1), command)
}

func TestNewCommandCreateBox(t *testing.T) {
	assert.Equal(t, "hck-box-create::my-group/my-name::main", NewCommandCreateBox("my-group/my-name", "main"))
}

func TestNewCommandExecBox(t *testing.T) {
	assert.Equal(t, "hck-box-exec::my-group/my-name::main::box-id", NewCommandExecBox("my-group/my-name", "main", "box-id"))
}

func TestNewCommandOpenBox(t *testing.T) {
	assert.Equal(t, "hck-box-open::my-group/my-name::main", NewCommandOpenBox("my-group/my-name", "main"))
}

func TestNewCommandListBox(t *testing.T) {
	assert.Equal(t, "hck-box-list", NewCommandListBox())
}

func TestNewCommandDeleteBox(t *testing.T) {
	assert.Equal(t, "hck-box-delete::my-group/my-name::main::box-id", NewCommandDeleteBox("my-group/my-name", "main", "box-id"))
}
