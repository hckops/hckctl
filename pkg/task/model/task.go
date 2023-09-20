package model

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"

	commonModel "github.com/hckops/hckctl/pkg/common/model"
	"github.com/hckops/hckctl/pkg/util"
)

const (
	tagPrefixName      = "task-"
	defaultTaskCommand = "default"
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
	Name      string
	Arguments []string
}

func (command *TaskCommand) ExpandCommandArguments(parameters commonModel.Parameters) ([]string, error) {
	var expandedArguments []string

	for _, argument := range command.Arguments {
		// image flags and commands must be separated
		for _, raw := range strings.Fields(argument) {
			expanded, err := util.Expand(raw, parameters)
			if err != nil {
				return nil, errors.Wrapf(err, "unable to expand argument %s", argument)
			}
			expandedArguments = append(expandedArguments, expanded)
		}
	}
	return expandedArguments, nil
}

func (task *TaskV1) GenerateName() string {
	return fmt.Sprintf("%s%s-%s", tagPrefixName, util.ToLowerKebabCase(task.Name), util.RandomAlphanumeric(5))
}

func (task *TaskV1) CommandMap() map[string]TaskCommand {
	commands := map[string]TaskCommand{}
	for _, command := range task.Commands {
		commands[command.Name] = command
	}
	return commands
}

func (task *TaskV1) DefaultCommand(name string) (TaskCommand, error) {
	if strings.TrimSpace(name) == "" {
		if command, ok := task.CommandMap()[defaultTaskCommand]; ok {
			return command, nil
		} else {
			return TaskCommand{}, fmt.Errorf("%s command not found", defaultTaskCommand)
		}
	}

	if command, ok := task.CommandMap()[name]; ok {
		return command, nil
	}
	return TaskCommand{}, fmt.Errorf("%s command not found", name)
}

func (task *TaskV1) Pretty() string {
	value, _ := util.EncodeJsonIndent(task)
	return value
}
