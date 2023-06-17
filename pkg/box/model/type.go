package model

import (
	"io"
	"os"

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
	return []BoxProvider{Docker, Kubernetes, ArgoCd, Cloud}
}

var providerValue = []string{"docker", "kube", "argo-cd", "cloud"}

func (p BoxProvider) String() string {
	return providerValue[p]
}

type BoxInfo struct {
	Id   string
	Name string
}

// TODO change resource to flag in command i.e s/m/l and move ResourceOptions here

//Memory: "512Mi",
//Cpu:    "500m",

type BoxOpts struct {
	Streams  *BoxStreams
	EventBus *event.EventBus
}

func NewBoxOpts() *BoxOpts {
	return &BoxOpts{
		Streams:  newDefaultStreams(true),
		EventBus: event.NewEventBus(),
	}
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
