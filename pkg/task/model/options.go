package model

import (
	"fmt"
	"path"
	"strconv"
	"time"

	commonModel "github.com/hckops/hckctl/pkg/common/model"
	"github.com/hckops/hckctl/pkg/event"
)

type TaskClientOptions struct {
	Provider   TaskProvider
	DockerOpts *commonModel.DockerOptions
	KubeOpts   *commonModel.KubeOptions
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
	Template    *TaskV1
	Arguments   []string
	Labels      commonModel.Labels
	NetworkInfo commonModel.NetworkInfo
	StreamOpts  *commonModel.StreamOptions
	ShareDir    string
	LogDir      string
}

func (opts *RunOptions) GenerateLogFileName(provider TaskProvider, containerName string) string {
	timestamp := strconv.FormatInt(time.Now().UTC().UnixNano(), 10)
	return path.Join(opts.LogDir, fmt.Sprintf("%s-%s-%s", provider.String(), timestamp, containerName))
}
