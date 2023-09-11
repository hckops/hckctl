package model

import (
	commonModel "github.com/hckops/hckctl/pkg/common/model"
	"github.com/hckops/hckctl/pkg/schema"
)

func NewLabLabels() commonModel.Labels {
	return map[string]string{
		commonModel.LabelSchemaKind: schema.KindLabV1.String(),
	}
}
