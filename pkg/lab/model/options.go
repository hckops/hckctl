package model

import (
	boxModel "github.com/hckops/hckctl/pkg/box/model"
	"github.com/hckops/hckctl/pkg/event"
	"github.com/hckops/hckctl/pkg/provider"
)

type LabClientOptions struct {
	Provider  LabProvider
	CloudOpts *provider.CloudOptions
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
	LabTemplate  *LabV1
	BoxTemplates map[string]*boxModel.BoxV1
	Parameters   map[string]string
	Labels       map[string]string
}
