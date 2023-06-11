package box

import (
	"io"
	"os"

	"github.com/hckops/hckctl/pkg/client"
	"github.com/hckops/hckctl/pkg/template/model"
)

type BoxProvider string

const (
	Docker     BoxProvider = "docker"
	Kubernetes BoxProvider = "kube"
	Argo       BoxProvider = "argo"
	Cloud      BoxProvider = "cloud"
)

func BoxProviderValues() []string {
	return []string{string(Docker), string(Kubernetes), string(Argo), string(Cloud)}
}

func BoxProviderFromEventSource(source client.EventSource) BoxProvider {
	switch source {
	case client.DockerSource:
		return Docker
	case client.KubeSource:
		return Kubernetes
	case client.ArgoSource:
		return Argo
	case client.CloudSource:
		return Cloud
	default:
		return "INVALID_SOURCE"
	}
}

type boxOpts struct {
	template *model.BoxV1
	streams  *boxStreams
	eventBus *client.EventBus
}

func newBoxOpts(template *model.BoxV1) *boxOpts {
	return &boxOpts{
		template: template,
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
