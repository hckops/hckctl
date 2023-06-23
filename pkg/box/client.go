package box

import (
	"github.com/hckops/hckctl/pkg/box/cloud"
	"github.com/pkg/errors"

	"github.com/hckops/hckctl/pkg/box/docker"
	"github.com/hckops/hckctl/pkg/box/kubernetes"
	"github.com/hckops/hckctl/pkg/box/model"
	"github.com/hckops/hckctl/pkg/event"
)

type BoxClient interface {
	Provider() model.BoxProvider
	Events() *event.EventBus
	Create(template *model.BoxV1) (*model.BoxInfo, error)
	Exec(template *model.BoxV1, name string) error // TODO exec --tunnel (docker do nothing)
	Open(template *model.BoxV1) error              // TODO open --tunnel (docker do nothing)
	List() ([]model.BoxInfo, error)
	Copy(name string, from string, to string) error
	Tunnel(name string) error
	Delete(name string) error
	DeleteAll() ([]model.BoxInfo, error)
}

// TODO https://stackoverflow.com/questions/30261032/how-to-implement-an-abstract-class-in-go
// TODO https://golangbyexample.com/go-abstract-class

func NewBoxClient(opts *model.BoxOpts) (BoxClient, error) {
	switch opts.Provider {
	case model.Docker:
		return docker.NewDockerBox(opts.InternalOpts)
	case model.Kubernetes:
		return kubernetes.NewKubeBox(opts.InternalOpts, opts.KubeConfig)
	case model.Cloud:
		return cloud.NewCloudBox(opts.InternalOpts, opts.CloudConfig)
	default:
		return nil, errors.New("invalid provider")
	}
}
