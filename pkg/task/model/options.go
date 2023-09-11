package model

import (
	commonModel "github.com/hckops/hckctl/pkg/common/model"
	"github.com/hckops/hckctl/pkg/event"
)

type TaskClientOptions struct {
	Provider   TaskProvider
	DockerOpts *commonModel.DockerOptions
}

type CommonTaskOptions struct {
	EventBus *event.EventBus
}

func NewCommonTaskOpts() *CommonTaskOptions {
	return &CommonTaskOptions{
		EventBus: event.NewEventBus(),
	}
}

type CreateOptions struct {
	TaskTemplate *TaskV1
	Parameters   map[string]string
	Labels       commonModel.Labels
}
