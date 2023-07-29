package docker

import (
	"github.com/pkg/errors"

	"github.com/hckops/hckctl/pkg/box/model"
	"github.com/hckops/hckctl/pkg/client/docker"
	"github.com/hckops/hckctl/pkg/event"
)

type DockerBoxClient struct {
	client     *docker.DockerClient
	clientOpts *model.DockerBoxOptions
	eventBus   *event.EventBus
}

func NewDockerBoxClient(commonOpts *model.CommonBoxOptions, dockerOpts *model.DockerBoxOptions) (*DockerBoxClient, error) {
	return newDockerBoxClient(commonOpts, dockerOpts)
}

func (box *DockerBoxClient) Provider() model.BoxProvider {
	return model.Docker
}

func (box *DockerBoxClient) Events() *event.EventBus {
	return box.eventBus
}

func (box *DockerBoxClient) Create(opts *model.CreateOptions) (*model.BoxInfo, error) {
	defer box.close()
	return box.createBox(opts)
}

func (box *DockerBoxClient) Connect(opts *model.ConnectOptions) error {
	defer box.close()
	return box.connectBox(opts)
}

func (box *DockerBoxClient) Copy(string, string, string) error {
	defer box.close()
	return errors.New("not implemented")
}

func (box *DockerBoxClient) Describe(name string) (*model.BoxDetails, error) {
	defer box.close()
	return box.describe(name)
}

func (box *DockerBoxClient) List() ([]model.BoxInfo, error) {
	defer box.close()
	return box.listBoxes()
}

func (box *DockerBoxClient) Delete(names []string) ([]string, error) {
	defer box.close()
	return box.deleteBoxes(names)
}

func (box *DockerBoxClient) Clean() error {
	defer box.close()
	// TODO remove network and volumes
	return errors.New("not implemented")
}

func (box *DockerBoxClient) Version() (string, error) {
	return "", errors.New("not implemented")
}
