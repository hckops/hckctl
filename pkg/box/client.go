package box

import (
	"context"

	"github.com/pkg/errors"

	"github.com/hckops/hckctl/pkg/template/model"
)

type Connection struct {
	ctx context.Context
	Out chan string // TODO
}

type BoxClient interface {
	Open() (*Connection, error)
}

func NewBoxClient(provider BoxProvider, template *model.BoxV1) (BoxClient, error) {
	switch provider {
	case Docker:
		return NewDockerClient(template), nil
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
