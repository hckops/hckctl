package model

import (
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
	Template   *LabV1
	Parameters map[string]string // TODO
	Labels     map[string]string // TODO refactor BoxLabels
}
