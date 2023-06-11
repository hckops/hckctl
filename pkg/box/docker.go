package box

import (
	"github.com/pkg/errors"

	"github.com/hckops/hckctl/pkg/client"
	"github.com/hckops/hckctl/pkg/client/docker"
)

type DockerBox struct {
	client *docker.DockerClient
	opts   *boxOpts
}

func NewDockerBox(opts *boxOpts) (*DockerBox, error) {

	dockerClient, err := docker.NewDockerClient(opts.eventBus)
	if err != nil {
		return nil, errors.Wrap(err, "error docker box")
	}

	return &DockerBox{
		client: dockerClient,
		opts:   opts,
	}, nil
}

func (b *DockerBox) Events() *client.EventBus {
	return b.opts.eventBus
}

func (b *DockerBox) Create() (*BoxInfo, error) {
	return nil, nil
}

func (b *DockerBox) Exec(boxId string) error {
	return nil
}

func (b *DockerBox) Copy(boxId string, from string, to string) error {
	return nil
}

func (b *DockerBox) List() ([]string, error) {
	return nil, nil
}

func (b *DockerBox) Open() error {
	return nil
}

func (b *DockerBox) Tunnel() error {
	return nil
}

func (b *DockerBox) Delete(boxId string) error {
	return nil
}
