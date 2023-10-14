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
	Provider() model.BoxProvider // TODO replace with generics
	Events() *event.EventBus
	Create(opts *model.CreateOptions) (*model.BoxInfo, error)
	Connect(opts *model.ConnectOptions) error
	Describe(name string) (*model.BoxDetails, error)
	List() ([]model.BoxInfo, error)
	Delete(names []string) ([]string, error) // empty "names" means all boxes
	Clean() error                            // TODO delete source in params: remove local and git cache
	Version() (string, error)                // TODO replace string with BoxVersion interface, return both client and server version
}

func NewBoxClient(opts *model.BoxClientOptions) (BoxClient, error) {
	commonOpts := model.NewCommonBoxOpts()
	switch opts.Provider {
	case model.Docker:
		return docker.NewDockerBoxClient(commonOpts, opts.DockerOpts)
	case model.Kubernetes:
		return kubernetes.NewKubeBoxClient(commonOpts, opts.KubeOpts)
	case model.Cloud:
		return cloud.NewCloudBoxClient(commonOpts, opts.CloudOpts)
	default:
		return nil, errors.New("invalid provider")
	}
}
