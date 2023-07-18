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
	Create(templateOpts *model.TemplateOptions) (*model.BoxInfo, error)
	Connect(template *model.BoxV1, tunnelOpts *model.TunnelOptions, name string) error
	Open(templateOpts *model.TemplateOptions, tunnelOpts *model.TunnelOptions) error
	Copy(name string, from string, to string) error
	Describe(name string) (*model.BoxDetails, error)
	List() ([]model.BoxInfo, error)
	Delete(names []string) ([]string, error) // empty "names" means all boxes
	Clean() error
	Version() (string, error) // TODO replace string with BoxVersion interface, return both client and server version
}

// TODO https://stackoverflow.com/questions/30261032/how-to-implement-an-abstract-class-in-go
// TODO https://golangbyexample.com/go-abstract-class

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
