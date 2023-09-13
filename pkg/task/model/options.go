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

type RunOptions struct {
	Template   *TaskV1
	Arguments  []string
	Labels     commonModel.Labels
	StreamOpts *commonModel.StreamOptions
}
