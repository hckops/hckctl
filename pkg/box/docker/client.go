package docker

import (
	"github.com/hckops/hckctl/pkg/box"
	"github.com/hckops/hckctl/pkg/template/model"
)

type DockerClientOpts struct {
	template *model.BoxV1
}

func NewDockerClient(template *model.BoxV1) *DockerClient {
	return &DockerClient{
		opts: &DockerClientOpts{
			template: template,
		},
	}
}

type DockerClient struct {
	opts *DockerClientOpts
}

func (client *DockerClient) Open() (*box.Connection, error) {
	return nil, nil
}
