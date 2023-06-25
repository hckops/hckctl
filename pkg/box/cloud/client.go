package cloud

import (
	"github.com/pkg/errors"

	"github.com/hckops/hckctl/pkg/box/model"
	"github.com/hckops/hckctl/pkg/client/ssh"
	"github.com/hckops/hckctl/pkg/event"
)

type CloudBox struct {
	clientVersion string
	clientConfig  *ssh.SshClientConfig
	client        *ssh.SshClient
	streams       *model.BoxStreams
	eventBus      *event.EventBus
}

func NewCloudBox(internalOpts *model.BoxInternalOptions, clientConfig *ssh.SshClientConfig) (*CloudBox, error) {
	return newCloudBox(internalOpts, clientConfig)
}

func (box *CloudBox) Provider() model.BoxProvider {
	return model.Cloud
}

func (box *CloudBox) Events() *event.EventBus {
	return box.eventBus
}

func (box *CloudBox) Create(template *model.BoxV1) (*model.BoxInfo, error) {
	defer box.close()
	return box.createBox(template)
}

func (box *CloudBox) Connect(template *model.BoxV1, tunnelOpts *model.TunnelOptions, name string) error {
	defer box.close()
	// TODO tunnelOpts
	return box.execBox(template, name)
}

func (box *CloudBox) Open(template *model.BoxV1, tunnelOpts *model.TunnelOptions) error {
	defer box.close()
	// TODO tunnelOpts
	return errors.New("not implemented")
}

func (box *CloudBox) Copy(string, string, string) error {
	defer box.close()
	return errors.New("not implemented")
}

func (box *CloudBox) List() ([]model.BoxInfo, error) {
	defer box.close()
	return box.listBoxes()
}

func (box *CloudBox) Delete(name string) error {
	defer box.close()
	_, err := box.deleteBoxes([]string{name})
	return err
}

func (box *CloudBox) DeleteAll() ([]model.BoxInfo, error) {
	defer box.close()
	return box.deleteBoxes([]string{})
}
