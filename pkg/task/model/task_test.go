package model

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	commonModel "github.com/hckops/hckctl/pkg/common/model"
)

func TestGenerateName(t *testing.T) {
	task := &TaskV1{
		Name: "my-name",
	}
	taskId := task.GenerateName()
	assert.True(t, strings.HasPrefix(taskId, "task-my-name-"))
	assert.Equal(t, 18, len(taskId))
}

func TestMainContainerName(t *testing.T) {
	testBox := &TaskV1{
		Image: commonModel.Image{
			Repository: "hckops/my-image",
		},
	}
	assert.Equal(t, "hckops-my-image", testBox.MainContainerName())
}

func TestCommandMap(t *testing.T) {
	task := &TaskV1{
		Commands: []TaskCommand{
			{Name: "default", Arguments: []string{"arg1", "arg2"}},
			{Name: "example", Arguments: []string{"arg3", "arg4"}},
		},
	}
	expected := map[string]TaskCommand{
		"default": {Name: "default", Arguments: []string{"arg1", "arg2"}},
		"example": {Name: "example", Arguments: []string{"arg3", "arg4"}},
	}
	commands := task.CommandMap()
	assert.Equal(t, expected, commands)
}

func TestLoadCommand(t *testing.T) {
	task := &TaskV1{
		Commands: []TaskCommand{
			{Name: "default", Arguments: []string{"arg1", "arg2"}},
		},
	}
	expected := TaskCommand{Name: "default", Arguments: []string{"arg1", "arg2"}}

	commandEmpty, errEmpty := task.LoadCommand("")
	assert.Equal(t, expected, commandEmpty)
	assert.Nil(t, errEmpty)

	commandDefault, errDefault := task.LoadCommand("default")
	assert.Equal(t, expected, commandDefault)
	assert.Nil(t, errDefault)

	commandInvalid, errInvalid := (&TaskV1{}).LoadCommand("foo")
	assert.Equal(t, TaskCommand{}, commandInvalid)
	assert.EqualError(t, errInvalid, "foo command not found")
}

func TestPretty(t *testing.T) {
	task := &TaskV1{
		Kind: "task/v1",
		Name: "whalesay",
		Tags: []string{"test"},
		Image: commonModel.Image{
			Repository: "docker/whalesay",
		},
		Commands: []TaskCommand{
			{Name: "default", Arguments: []string{"cowsay", "${hello:hckops}"}},
		},
	}
	json := `{
  "Kind": "task/v1",
  "Name": "whalesay",
  "Tags": [
    "test"
  ],
  "Image": {
    "Repository": "docker/whalesay",
    "Version": ""
  },
  "Commands": [
    {
      "Name": "default",
      "Arguments": [
        "cowsay",
        "${hello:hckops}"
      ]
    }
  ]
}`
	assert.Equal(t, json, task.Pretty())
}

func TestExpandCommandArguments(t *testing.T) {
	command := TaskCommand{Arguments: []string{
		" -a ",
		"-b bbb",
		"-c ${ccc:CCC}",
		"-d ${ddd:DDD}",
		"-e f --g HHH",
		"-l ${lll:LLL:MMM:NNN}",
	}}
	parameters := commonModel.Parameters{
		"bbb": "BBB",
		"ddd": "AAA",
	}
	expected := []string{
		"-a", "-b", "bbb", "-c", "CCC", "-d", "AAA", "-e", "f", "--g", "HHH", "-l", "LLL:MMM:NNN",
	}
	expanded, err := command.ExpandCommandArguments(parameters)

	assert.Len(t, expanded, 13)
	assert.Equal(t, expected, expanded)
	assert.Nil(t, err)
}
