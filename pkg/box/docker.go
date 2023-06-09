package box

import (
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

func (client *DockerClient) Setup() (*Connection, error) {
	return nil, nil
}

// TODO BoxInfo e.g. id, template

func (client *DockerClient) Create() (string, error) {
	return "", nil
}
