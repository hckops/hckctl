package docker

import (
	"github.com/pkg/errors"

	"github.com/hckops/hckctl/pkg/box/model"
	"github.com/hckops/hckctl/pkg/client/docker"
	"github.com/hckops/hckctl/pkg/event"
)

type DockerBox struct {
	client       *docker.DockerClient
	clientConfig *docker.DockerClientConfig
	eventBus     *event.EventBus
}

func NewDockerBox(commonOpts *model.BoxCommonOptions, clientConfig *docker.DockerClientConfig) (*DockerBox, error) {
	return newDockerBox(commonOpts, clientConfig)
}

func (box *DockerBox) Provider() model.BoxProvider {
	return model.Docker
}

func (box *DockerBox) Events() *event.EventBus {
	return box.eventBus
}

func (box *DockerBox) Create(templateOpts *model.TemplateOptions) (*model.BoxInfo, error) {
	defer box.close()
	return box.createBox(templateOpts)
}

func (box *DockerBox) Connect(template *model.BoxV1, tunnelOpts *model.TunnelOptions, name string) error {
	defer box.close()
	return box.connectBox(template, tunnelOpts, name)
}

func (box *DockerBox) Open(templateOpts *model.TemplateOptions, tunnelOpts *model.TunnelOptions) error {
	defer box.close()
	return box.openBox(templateOpts, tunnelOpts)
}

func (box *DockerBox) Copy(string, string, string) error {
	defer box.close()
	return errors.New("not implemented")
}

func (box *DockerBox) List() ([]model.BoxInfo, error) {
	defer box.close()
	return box.listBoxes()
}

func (box *DockerBox) Delete(name string) error {
	defer box.close()
	return box.deleteBoxByName(name)
}

func (box *DockerBox) DeleteAll() ([]model.BoxInfo, error) {
	defer box.close()
	return box.deleteBoxes()
}

func (box *DockerBox) Version() (string, error) {
	return "", errors.New("not implemented")
}
