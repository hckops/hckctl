package box

import (
	"github.com/pkg/errors"

	"github.com/hckops/hckctl/pkg/template/model"
)

// https://go.dev/blog/pipelines

type EventKind uint8

const (
	DebugEvent EventKind = iota
	InfoEvent
	SuccessEvent
	ErrorEvent
)

type BoxEvent struct {
	Kind    EventKind
	Source  string
	Message string
}

type BoxClient interface {
	Events() <-chan BoxEvent
	Create() (string, error)
	//Copy()
	//Delete()
	//Exec()
	//List() ([]string, error)
	//Open()
	//Tunnel()
}

func NewBoxClient(provider BoxProvider, template *model.BoxV1) (BoxClient, error) {
	switch provider {
	case Docker:
		return NewDockerClient(template)
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
