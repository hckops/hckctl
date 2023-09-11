package model

import (
	"fmt"
	"strings"

	commonModel "github.com/hckops/hckctl/pkg/common/model"
	"github.com/hckops/hckctl/pkg/schema"
	"github.com/hckops/hckctl/pkg/util"
)

const (
	LabelBoxSize = "com.hckops.box.size"
)

func NewBoxLabels() commonModel.Labels {
	return map[string]string{
		commonModel.LabelSchemaKind: schema.KindBoxV1.String(),
	}
}

func AddBoxSize(labels commonModel.Labels, size ResourceSize) commonModel.Labels {
	return labels.AddLabel(LabelBoxSize, strings.ToLower(size.String()))
}

func ToBoxSize(labels commonModel.Labels) (ResourceSize, error) {
	if label, err := labels.Exist(LabelBoxSize); err != nil {
		return Small, err
	} else {
		return ExistResourceSize(label)
	}
}

func BoxLabelSelector() string {
	// value must be sanitized
	return fmt.Sprintf("%s=%s", commonModel.LabelSchemaKind, util.ToLowerKebabCase(schema.KindBoxV1.String()))
}
