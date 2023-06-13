package box

import (
	"github.com/pkg/errors"

	"github.com/hckops/hckctl/pkg/box/docker"
	"github.com/hckops/hckctl/pkg/template/model"
)

type BoxClient interface {
	Events() *EventBus
	Create(template *model.BoxV1) (*BoxInfo, error)
	Exec(name string, command string) error
	Copy(name string, from string, to string) error
	List() ([]BoxInfo, error)
	Open(template *model.BoxV1) error
	Tunnel(name string) error
	Delete(name string) error
}

func NewBoxClient(provider BoxProvider) (BoxClient, error) {
	opts := newBoxOpts()
	switch provider {
	case Docker:
		return docker.NewDockerBox(opts)
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
