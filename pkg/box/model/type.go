package model

import (
	"io"
	"os"

	"github.com/hckops/hckctl/pkg/client/cloud"
	"github.com/hckops/hckctl/pkg/client/kubernetes"
	"github.com/hckops/hckctl/pkg/event"
)

type BoxProvider uint

const (
	Docker BoxProvider = iota
	Kubernetes
	ArgoCd
	Cloud
)

func BoxProviders() []BoxProvider {
	// only supported: Kubernetes, ArgoCd, Cloud
	return []BoxProvider{Docker}
}

var providerValue = []string{"docker", "kube", "argo-cd", "cloud"}

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
	CloudConfig  *cloud.CloudClientConfig
	InternalOpts *BoxInternalOpts
}

func NewBoxOpts(provider BoxProvider, kubeConfig *kubernetes.KubeClientConfig, cloudConfig *cloud.CloudClientConfig) *BoxOpts {
	return &BoxOpts{
		Provider:    provider,
		KubeConfig:  kubeConfig,
		CloudConfig: cloudConfig,
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
	IsTty bool // tty false for tunnel only
}

func newDefaultStreams(tty bool) *BoxStreams {
	return &BoxStreams{
		In:    os.Stdin,
		Out:   os.Stdout,
		Err:   os.Stderr,
		IsTty: tty,
	}
}
