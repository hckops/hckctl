package model

import (
	"io"
	"os"

	"github.com/hckops/hckctl/pkg/event"
)

// TODO refactor to iota and type map[][] + String
type BoxProvider string

const (
	Docker     BoxProvider = "docker"
	Kubernetes BoxProvider = "kube"
	Argo       BoxProvider = "argo"
	Cloud      BoxProvider = "cloud"
)

func BoxProviders() []BoxProvider {
	return []BoxProvider{Docker, Kubernetes, Argo, Cloud}
}

func BoxProviderValues() []string {
	var values []string
	for _, provider := range BoxProviders() {
		values = append(values, string(provider))
	}
	return values
}

// TODO add provider
type BoxInfo struct {
	Id   string
	Name string
}

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
