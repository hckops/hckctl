package model

import (
	commonModel "github.com/hckops/hckctl/pkg/common/model"
	"github.com/hckops/hckctl/pkg/schema"
)

func NewTaskLabels() commonModel.Labels {
	return map[string]string{
		commonModel.LabelSchemaKind: schema.KindTaskV1.String(),
	}
}
