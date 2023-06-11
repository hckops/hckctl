package box

import (
	"github.com/pkg/errors"

	"github.com/hckops/hckctl/pkg/client"
	"github.com/hckops/hckctl/pkg/template/model"
)

type BoxClient interface {
	Events() *client.EventBus
	Create() (*BoxInfo, error)
	Exec(info BoxInfo) error
	Copy(info BoxInfo, from string, to string) error
	List() ([]BoxInfo, error)
	Open() error
	Tunnel(info BoxInfo) error
	Delete(info BoxInfo) error
}

func NewBoxClient(provider BoxProvider, template *model.BoxV1) (BoxClient, error) {
	opts := newBoxOpts(template)
	switch provider {
	case Docker:
		return NewDockerBox(opts)
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
