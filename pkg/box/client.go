package box

import (
	"github.com/pkg/errors"

	"github.com/hckops/hckctl/pkg/template/model"
)

type BoxClient interface {
	Events() *EventBus
	Create() (string, error)
	//Copy()
	//Delete()
	//Exec()
	//List() ([]string, error)
	//Open()
	//Tunnel()
}

func NewBoxClient(provider BoxProvider, template *model.BoxV1) (BoxClient, error) {
	eventBus := NewEventBus()
	switch provider {
	case Docker:
		return NewDockerClient(template, eventBus)
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
