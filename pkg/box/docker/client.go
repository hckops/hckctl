package docker

import (
	"github.com/pkg/errors"

	boxModel "github.com/hckops/hckctl/pkg/box/model"
	"github.com/hckops/hckctl/pkg/client/docker"
	commonDocker "github.com/hckops/hckctl/pkg/common/docker"
	commonModel "github.com/hckops/hckctl/pkg/common/model"
	"github.com/hckops/hckctl/pkg/event"
)

type DockerBoxClient struct {
	client       *docker.DockerClient
	clientOpts   *commonModel.DockerOptions
	dockerCommon *commonDocker.DockerCommonClient
	eventBus     *event.EventBus
}

func NewDockerBoxClient(commonOpts *boxModel.CommonBoxOptions, dockerOpts *commonModel.DockerOptions) (*DockerBoxClient, error) {
	return newDockerBoxClient(commonOpts, dockerOpts)
}

func (box *DockerBoxClient) Provider() boxModel.BoxProvider {
	return boxModel.Docker
}

func (box *DockerBoxClient) Events() *event.EventBus {
	return box.eventBus
}

func (box *DockerBoxClient) Create(opts *boxModel.CreateOptions) (*boxModel.BoxInfo, error) {
	defer box.close()
	return box.createBox(opts)
}

func (box *DockerBoxClient) Connect(opts *boxModel.ConnectOptions) error {
	defer box.close()
	return box.connectBox(opts)
}

func (box *DockerBoxClient) Copy(string, string, string) error {
	defer box.close()
	return errors.New("not implemented")
}

func (box *DockerBoxClient) Describe(name string) (*boxModel.BoxDetails, error) {
	defer box.close()
	return box.describeBox(name)
}

func (box *DockerBoxClient) List() ([]boxModel.BoxInfo, error) {
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
