package box

import (
	"github.com/pkg/errors"

	"github.com/hckops/hckctl/pkg/box/docker"
	"github.com/hckops/hckctl/pkg/box/kubernetes"
	"github.com/hckops/hckctl/pkg/box/model"
	"github.com/hckops/hckctl/pkg/event"
)

type BoxClient interface {
	Events() *event.EventBus
	Create(template *model.BoxV1) (*model.BoxInfo, error)
	Exec(name string, command string) error
	Open(template *model.BoxV1) error
	List() ([]model.BoxInfo, error)
	Copy(name string, from string, to string) error
	Tunnel(name string) error
	Delete(name string) error
	DeleteAll() ([]model.BoxInfo, error)
}

func NewBoxClient(provider model.BoxProvider) (BoxClient, error) {
	opts := model.NewBoxOpts()
	switch provider {
	case model.Docker:
		return docker.NewDockerBox(opts)
	case model.Kubernetes:
		return kubernetes.NewKubeBox(opts)
	case model.ArgoCd:
		// TODO
		return nil, nil
	case model.Cloud:
		// TODO
		return nil, nil
	default:
		return nil, errors.New("invalid provider")
	}
}
