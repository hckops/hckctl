package box

import (
	"github.com/pkg/errors"

	"github.com/hckops/hckctl/pkg/box"
	"github.com/hckops/hckctl/pkg/box/docker"
	"github.com/hckops/hckctl/pkg/template/model"
)

func newBoxClient(provider box.BoxProvider, boxTemplate *model.BoxV1) (box.BoxClient, error) {
	switch provider {
	case box.Docker:
		return docker.NewDockerClient(boxTemplate), nil
	case box.Kubernetes:
		// TODO
		return nil, nil
	case box.Argo:
		// TODO
		return nil, nil
	case box.Cloud:
		// TODO
		return nil, nil
	default:
		return nil, errors.New("invalid provider")
	}
}
