package model

import (
	commonModel "github.com/hckops/hckctl/pkg/common/model"
	"github.com/hckops/hckctl/pkg/event"
	"io"
	"os"
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

// TODO generic ClientStreams?
type TaskStreams struct {
	In  io.ReadCloser
	Out io.Writer
	Err io.Writer
}

func NewTaskStreams() *TaskStreams {
	return &TaskStreams{
		In:  os.Stdin,
		Out: os.Stdout,
		Err: os.Stderr,
	}
}
