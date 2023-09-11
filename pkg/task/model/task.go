package model

import (
	commonModel "github.com/hckops/hckctl/pkg/common/model"
)

type TaskV1 struct {
	Kind  string
	Name  string
	Tags  []string
	Image commonModel.Image
}
