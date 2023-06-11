package box

import (
	"github.com/pkg/errors"

	"github.com/hckops/hckctl/pkg/client"
	"github.com/hckops/hckctl/pkg/template/model"
)

type BoxClient interface {
	Events() *client.EventBus
	Create() (*BoxInfo, error)
	Exec(boxId string) error // boxId == containerName
	Copy(boxId string, from string, to string) error
	List() ([]string, error)
	Open() error
	Tunnel() error
	Delete(boxId string) error
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
