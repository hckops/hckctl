package model

import (
	"fmt"

	commonModel "github.com/hckops/hckctl/pkg/common/model"
	"github.com/hckops/hckctl/pkg/util"
)

const (
	TagPrefixName      = "task-"
	DefaultTaskCommand = "default"
)

// TODO add/review output, pages, license

type TaskV1 struct {
	Kind     string
	Name     string
	Tags     []string
	Image    commonModel.Image
	Commands []TaskCommand
}

type TaskCommand struct {
	Name string
	Args []string
}

func (task *TaskV1) GenerateName() string {
	return fmt.Sprintf("%s%s-%s", TagPrefixName, util.ToLowerKebabCase(task.Name), util.RandomAlphanumeric(5))
}

func (task *TaskV1) CommandMap() map[string]TaskCommand {
	commands := map[string]TaskCommand{}
	for _, command := range task.Commands {
		commands[command.Name] = command
	}
	return commands
}

func (task *TaskV1) DefaultCommandArgs() []string {
	if command, ok := task.CommandMap()[DefaultTaskCommand]; ok {
		return command.Args
	}
	return []string{}
}
