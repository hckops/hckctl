package model

import (
	"io"
	"os"

	"github.com/hckops/hckctl/pkg/event"
)

type BoxClientOptions struct {
	Provider   BoxProvider
	DockerOpts *DockerBoxOptions
	KubeOpts   *KubeBoxOptions
	CloudOpts  *CloudBoxOptions
}

type CommonBoxOptions struct {
	EventBus *event.EventBus
}

func NewCommonBoxOpts() *CommonBoxOptions {
	return &CommonBoxOptions{
		EventBus: event.NewEventBus(),
	}
}

type DockerBoxOptions struct {
	NetworkName          string
	IgnoreImagePullError bool
}

type KubeBoxOptions struct {
	InCluster  bool
	ConfigPath string
	Namespace  string
}

type CloudBoxOptions struct {
	Version  string
	Address  string
	Username string
	Token    string
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
	Labels   BoxLabels
}

type ConnectOptions struct {
	Template      *BoxV1
	Streams       *BoxStreams
	Name          string
	DisableExec   bool
	DisableTunnel bool
	DeleteOnExit  bool
}
