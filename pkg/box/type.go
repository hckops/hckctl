package box

import (
	"io"
	"os"
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
	EventBus *EventBus
}

func newBoxOpts() *BoxOpts {
	return &BoxOpts{
		Streams:  newDefaultStreams(true),
		EventBus: newEventBus(),
	}
}

type BoxStreams struct {
	In    io.Reader
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
