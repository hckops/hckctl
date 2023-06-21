package box

import (
	"github.com/pkg/errors"

	"github.com/hckops/hckctl/pkg/box/docker"
	"github.com/hckops/hckctl/pkg/box/kubernetes"
	"github.com/hckops/hckctl/pkg/box/model"
	"github.com/hckops/hckctl/pkg/event"
)

// TODO move util.terminal in client package

type BoxClient interface {
	Provider() model.BoxProvider
	Events() *event.EventBus
	Create(template *model.BoxV1) (*model.BoxInfo, error)
	Exec(name string, command string) error
	Open(template *model.BoxV1) error // TODO exec --tunnel (docker do nothing)
	List() ([]model.BoxInfo, error)
	Copy(name string, from string, to string) error
	Tunnel(name string) error // TODO exec --tunnel (docker do nothing)
	Delete(name string) error
	DeleteAll() ([]model.BoxInfo, error)
}

func NewBoxClient(opts *model.BoxOpts) (BoxClient, error) {
	switch opts.Provider {
	case model.Docker:
		return docker.NewDockerBox(opts.InternalOpts)
	case model.Kubernetes:
		return kubernetes.NewKubeBox(opts.InternalOpts, opts.KubeConfig)
	case model.Cloud:
		// TODO
		return nil, nil
	default:
		return nil, errors.New("invalid provider")
	}
}
