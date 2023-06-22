package docker

import (
	"github.com/pkg/errors"

	"github.com/hckops/hckctl/pkg/box/model"
	"github.com/hckops/hckctl/pkg/client/docker"
	"github.com/hckops/hckctl/pkg/event"
)

type DockerBox struct {
	client   *docker.DockerClient
	streams  *model.BoxStreams
	eventBus *event.EventBus
}

func NewDockerBox(internalOpts *model.BoxInternalOpts) (*DockerBox, error) {
	return newDockerBox(internalOpts)
}

func (box *DockerBox) Provider() model.BoxProvider {
	return model.Docker
}

func (box *DockerBox) Events() *event.EventBus {
	return box.eventBus
}

func (box *DockerBox) Create(template *model.BoxV1) (*model.BoxInfo, error) {
	defer box.close()
	return box.createBox(template)
}

func (box *DockerBox) Exec(template *model.BoxV1, name string) error {
	defer box.close()
	return box.execBox(name, template.Shell)
}

func (box *DockerBox) Open(template *model.BoxV1) error {
	defer box.close()
	return box.openBox(template)
}

func (box *DockerBox) List() ([]model.BoxInfo, error) {
	defer box.close()
	return box.listBoxes()
}

func (box *DockerBox) Copy(string, string, string) error {
	defer box.close()
	return errors.New("not implemented")
}

func (box *DockerBox) Tunnel(string) error {
	defer box.close()
	return errors.New("not supported")
}

func (box *DockerBox) Delete(name string) error {
	defer box.close()
	return box.deleteBoxByName(name)
}

func (box *DockerBox) DeleteAll() ([]model.BoxInfo, error) {
	defer box.close()
	return box.deleteBoxes()
}
