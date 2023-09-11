package model

import (
	"fmt"

	commonModel "github.com/hckops/hckctl/pkg/common/model"
	"github.com/hckops/hckctl/pkg/util"
)

const (
	TagPrefixName = "task-"
)

type TaskV1 struct {
	Kind  string
	Name  string
	Tags  []string
	Image commonModel.Image
}

func (task *TaskV1) GenerateName() string {
	return fmt.Sprintf("%s%s-%s", TagPrefixName, util.ToLowerKebabCase(task.Name), util.RandomAlphanumeric(5))
}
