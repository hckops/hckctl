package model

import (
	"io"
	"os"

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

type BoxStreams struct {
	In    io.ReadCloser
	Out   io.Writer
	Err   io.Writer
	IsTty bool // tty is false only for ssh tunnel
}

func NewDefaultStreams(tty bool) *BoxStreams {
	return &BoxStreams{
		In:    os.Stdin,
		Out:   os.Stdout,
		Err:   os.Stderr,
		IsTty: tty,
	}
}

type CreateOptions struct {
	Template *BoxV1
	Size     ResourceSize
	Labels   commonModel.Labels
}

type ConnectOptions struct {
	Template      *BoxV1
	Streams       *BoxStreams
	Name          string
	DisableExec   bool
	DisableTunnel bool
	DeleteOnExit  bool
}
