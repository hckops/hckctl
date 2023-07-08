package model

import (
	"io"
	"os"

	"github.com/hckops/hckctl/pkg/client/docker"
	"github.com/hckops/hckctl/pkg/client/kubernetes"
	"github.com/hckops/hckctl/pkg/client/ssh"
	"github.com/hckops/hckctl/pkg/event"
)

type BoxInfo struct {
	Id   string
	Name string
}

type BoxClientOptions struct {
	Provider     BoxProvider
	CommonOpts   *BoxCommonOptions
	DockerConfig *docker.DockerClientConfig
	KubeConfig   *kubernetes.KubeClientConfig
	SshConfig    *ssh.SshClientConfig
}

type BoxCommonOptions struct {
	Version      string // TODO only cloud
	AllowOffline bool   // TODO only docker
	Streams      *BoxStreams
	EventBus     *event.EventBus
}

func NewBoxCommonOpts(version string) *BoxCommonOptions {
	return &BoxCommonOptions{
		Version:      version,
		AllowOffline: true, // always allow to start offline/obsolete images
		Streams:      newDefaultStreams(true),
		EventBus:     event.NewEventBus(),
	}
}

type BoxStreams struct {
	In    io.ReadCloser
	Out   io.Writer
	Err   io.Writer
	IsTty bool // tty is false only for ssh tunnel
}

func newDefaultStreams(tty bool) *BoxStreams {
	return &BoxStreams{
		In:    os.Stdin,
		Out:   os.Stdout,
		Err:   os.Stderr,
		IsTty: tty,
	}
}

type TunnelOptions struct {
	TunnelOnly bool
	NoTunnel   bool
}
