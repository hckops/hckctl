package model

import (
	"io"
	"os"

	"github.com/hckops/hckctl/pkg/client/docker"
	"github.com/hckops/hckctl/pkg/client/kubernetes"
	"github.com/hckops/hckctl/pkg/client/ssh"
	"github.com/hckops/hckctl/pkg/event"
)

// TODO add actual bound ports
// TODO add labels info
type BoxInfo struct {
	Id      string
	Name    string
	Healthy bool // TODO BoxStatus healthy, offline, unknown, error, etc
}

type BoxClientOptions struct {
	Provider     BoxProvider
	CommonOpts   *BoxCommonOptions
	DockerConfig *docker.DockerClientConfig
	KubeConfig   *kubernetes.KubeClientConfig
	SshConfig    *ssh.SshClientConfig
}

type BoxCommonOptions struct {
	EventBus *event.EventBus
}

func NewBoxCommonOpts() *BoxCommonOptions {
	return &BoxCommonOptions{
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

type TemplateOptions struct {
	Template *BoxV1
	Size     ResourceSize
	Labels   BoxLabels
}

type TunnelOptions struct {
	Streams    *BoxStreams
	TunnelOnly bool
	NoTunnel   bool
}
