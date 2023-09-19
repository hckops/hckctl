package model

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateName(t *testing.T) {
	task := &TaskV1{
		Name: "my-name",
	}
	taskId := task.GenerateName()
	assert.True(t, strings.HasPrefix(taskId, "task-my-name-"))
	assert.Equal(t, 18, len(taskId))
}

func TestCommandMap(t *testing.T) {
	task := &TaskV1{
		Commands: []TaskCommand{
			{Name: "default", Args: []string{"arg1", "arg2"}},
			{Name: "example", Args: []string{"arg3", "arg4"}},
		},
	}
	expected := map[string]TaskCommand{
		"default": {Name: "default", Args: []string{"arg1", "arg2"}},
		"example": {Name: "example", Args: []string{"arg3", "arg4"}},
	}
	commands := task.CommandMap()
	assert.Equal(t, expected, commands)
}

func TestDefaultCommandArgs(t *testing.T) {
	task := &TaskV1{
		Commands: []TaskCommand{
			{Name: "default", Args: []string{"arg1", "arg2"}},
		},
	}
	expected := []string{"arg1", "arg2"}
	arguments := task.DefaultCommandArgs()
	assert.Equal(t, expected, arguments)

	emptyArgs := (&TaskV1{}).DefaultCommandArgs()
	assert.Equal(t, []string{}, emptyArgs)
}