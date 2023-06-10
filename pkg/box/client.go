package box

import (
	"github.com/pkg/errors"
	"io"
	"os"

	"github.com/hckops/hckctl/pkg/template/model"
)

type BoxStreams struct {
	In    io.Reader
	Out   io.Writer
	Err   io.Writer
	IsTty bool // tty false for tunnel only
}

func NewDefaultStreams(tty bool) *BoxStreams {
	return &BoxStreams{
		In:    os.Stdin,
		Out:   os.Stdout,
		Err:   os.Stderr,
		IsTty: tty,
	}
}

type BoxInfo struct {
	Id   string
	Name string
}

type BoxClient interface {
	Events() *EventBus
	Create() (*BoxInfo, error)
	//Copy()
	//Delete()
	Exec(boxId string) error
	//List() ([]string, error)
	Open() error
	//Tunnel()
}

func NewBoxClient(provider BoxProvider, template *model.BoxV1) (BoxClient, error) {
	eventBus := NewEventBus()
	streams := NewDefaultStreams(true)
	switch provider {
	case Docker:
		return NewDockerClient(template, streams, eventBus)
	case Kubernetes:
		// TODO
		return nil, nil
	case Argo:
		// TODO
		return nil, nil
	case Cloud:
		// TODO
		return nil, nil
	default:
		return nil, errors.New("invalid provider")
	}
}
