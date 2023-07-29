package box

import (
	"github.com/pkg/errors"

	"github.com/hckops/hckctl/pkg/box/cloud"
	"github.com/hckops/hckctl/pkg/box/docker"
	"github.com/hckops/hckctl/pkg/box/kubernetes"
	"github.com/hckops/hckctl/pkg/box/model"
	"github.com/hckops/hckctl/pkg/event"
)

type BoxClient interface {
	Provider() model.BoxProvider
	Events() *event.EventBus
	Create(opts *model.CreateOptions) (*model.BoxInfo, error)
	Connect(opts *model.ConnectOptions) error
	Copy(name string, from string, to string) error // TODO not implemented
	Describe(name string) (*model.BoxDetails, error)
	List() ([]model.BoxInfo, error)
	Delete(names []string) ([]string, error) // empty "names" means all boxes
	Clean() error                            // TODO delete source in params: remove local and git cache
	Version() (string, error)                // TODO replace string with BoxVersion interface, return both client and server version
}

func NewBoxClient(opts *model.BoxClientOptions) (BoxClient, error) {
	switch opts.Provider {
	case model.Docker:
		return docker.NewDockerBoxClient(opts.CommonOpts, opts.DockerOpts)
	case model.Kubernetes:
		return kubernetes.NewKubeBoxClient(opts.CommonOpts, opts.KubeOpts)
	case model.Cloud:
		return cloud.NewCloudBoxClient(opts.CommonOpts, opts.CloudOpts)
	default:
		return nil, errors.New("invalid provider")
	}
}
