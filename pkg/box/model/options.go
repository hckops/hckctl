package model

import (
	commonModel "github.com/hckops/hckctl/pkg/common/model"
	"github.com/hckops/hckctl/pkg/event"
)

type BoxClientOptions struct {
	Provider   BoxProvider
	DockerOpts *commonModel.DockerOptions
	KubeOpts   *commonModel.KubeOptions
	CloudOpts  *commonModel.CloudOptions
}

type CommonBoxOptions struct {
	EventBus *event.EventBus
}

func NewCommonBoxOpts() *CommonBoxOptions {
	return &CommonBoxOptions{
		EventBus: event.NewEventBus(),
	}
}

type CreateOptions struct {
	Template    *BoxV1
	Size        ResourceSize
	Labels      commonModel.Labels
	NetworkInfo commonModel.NetworkInfo
}

type ConnectOptions struct {
	Template      *BoxV1
	StreamOpts    *commonModel.StreamOptions
	Name          string
	DisableExec   bool
	DisableTunnel bool
	DeleteOnExit  bool
}
