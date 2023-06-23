package model

import (
	"github.com/hckops/hckctl/pkg/client/ssh"
	"io"
	"os"

	"github.com/hckops/hckctl/pkg/client/kubernetes"
	"github.com/hckops/hckctl/pkg/event"
)

type BoxProvider uint

const (
	Docker BoxProvider = iota
	Kubernetes
	Cloud
)

var providerValue = []string{"docker", "kube", "cloud"}

func (p BoxProvider) String() string {
	return providerValue[p]
}

type BoxInfo struct {
	Id   string
	Name string
}

type BoxOpts struct {
	Provider     BoxProvider
	KubeConfig   *kubernetes.KubeClientConfig
	SshConfig    *ssh.SshClientConfig
	InternalOpts *BoxInternalOpts
}

func NewBoxOpts(provider BoxProvider, kubeConfig *kubernetes.KubeClientConfig, sshConfig *ssh.SshClientConfig) *BoxOpts {
	return &BoxOpts{
		Provider:   provider,
		KubeConfig: kubeConfig,
		SshConfig:  sshConfig,
		InternalOpts: &BoxInternalOpts{
			Streams:  newDefaultStreams(true),
			EventBus: event.NewEventBus(),
		},
	}
}

type BoxInternalOpts struct {
	Streams  *BoxStreams
	EventBus *event.EventBus
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
