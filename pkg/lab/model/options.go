package model

import (
	boxModel "github.com/hckops/hckctl/pkg/box/model"
	commonModel "github.com/hckops/hckctl/pkg/common/model"
	"github.com/hckops/hckctl/pkg/event"
)

type LabClientOptions struct {
	Provider  LabProvider
	CloudOpts *commonModel.CloudOptions
}

type CommonLabOptions struct {
	EventBus *event.EventBus
}

func NewCommonLabOpts() *CommonLabOptions {
	return &CommonLabOptions{
		EventBus: event.NewEventBus(),
	}
}

type CreateOptions struct {
	LabTemplate   *LabV1
	BoxTemplates  map[string]*boxModel.BoxV1
	DumpTemplates map[string]*DumpV1
	Parameters    map[string]string
	Labels        map[string]string
}
