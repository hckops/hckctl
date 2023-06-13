package box

import (
	"io"
	"os"

	"github.com/hckops/hckctl/pkg/client"
)

type BoxProvider string

const (
	Docker     BoxProvider = "docker"
	Kubernetes BoxProvider = "kube"
	Argo       BoxProvider = "argo" // TODO remove, only labs
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

type boxOpts struct {
	streams  *boxStreams
	eventBus *client.EventBus
}

func newBoxOpts() *boxOpts {
	return &boxOpts{
		streams:  newDefaultStreams(true),
		eventBus: client.NewEventBus(),
	}
}

type boxStreams struct {
	in    io.Reader
	out   io.Writer
	err   io.Writer
	isTty bool // tty false for tunnel only
}

func newDefaultStreams(tty bool) *boxStreams {
	return &boxStreams{
		in:    os.Stdin,
		out:   os.Stdout,
		err:   os.Stderr,
		isTty: tty,
	}
}

type BoxInfo struct {
	Id   string
	Name string
}
