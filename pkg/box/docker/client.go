package docker

import (
	"github.com/pkg/errors"

	"github.com/hckops/hckctl/pkg/box/model"
	"github.com/hckops/hckctl/pkg/client/docker"
	"github.com/hckops/hckctl/pkg/event"
)

type DockerBox struct {
	client *docker.DockerClient
	opts   *model.BoxOpts
}

func NewDockerBox(opts *model.BoxOpts) (*DockerBox, error) {
	return newDockerBox(opts)
}

func (box *DockerBox) Events() *event.EventBus {
	return box.opts.EventBus
}

func (box *DockerBox) Create(template *model.BoxV1) (*model.BoxInfo, error) {
	defer box.close()
	return box.createBox(template)
}

func (box *DockerBox) Exec(name string, command string) error {
	defer box.client.Close()
	return box.execBox(name, command)
}

func (box *DockerBox) Open(template *model.BoxV1) error {
	defer box.client.Close()
	return box.openBox(template)
}

func (box *DockerBox) List() ([]model.BoxInfo, error) {
	defer box.client.Close()
	return box.listBoxes()
}

func (box *DockerBox) Copy(string, string, string) error {
	defer box.client.Close()
	return errors.New("not implemented")
}

func (box *DockerBox) Tunnel(string) error {
	defer box.client.Close()
	return errors.New("not supported")
}

func (box *DockerBox) Delete(name string) error {
	defer box.client.Close()
	return box.deleteBoxByName(name)
}

func (box *DockerBox) DeleteAll() ([]model.BoxInfo, error) {
	defer box.client.Close()
	return box.deleteBoxes()
}
