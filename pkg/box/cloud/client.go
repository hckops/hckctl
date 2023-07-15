package cloud

import (
	"github.com/pkg/errors"

	"github.com/hckops/hckctl/pkg/box/model"
	"github.com/hckops/hckctl/pkg/client/ssh"
	"github.com/hckops/hckctl/pkg/event"
)

type CloudBox struct {
	clientConfig *ssh.SshClientConfig
	client       *ssh.SshClient
	eventBus     *event.EventBus
}

func NewCloudBox(commonOpts *model.BoxCommonOptions, clientConfig *ssh.SshClientConfig) (*CloudBox, error) {
	return newCloudBox(commonOpts, clientConfig)
}

func (box *CloudBox) Provider() model.BoxProvider {
	return model.Cloud
}

func (box *CloudBox) Events() *event.EventBus {
	return box.eventBus
}

func (box *CloudBox) Create(templateOpts *model.TemplateOptions) (*model.BoxInfo, error) {
	defer box.close()
	return box.createBox(templateOpts)
}

func (box *CloudBox) Connect(template *model.BoxV1, tunnelOpts *model.TunnelOptions, name string) error {
	defer box.close()
	return box.execBox(template, tunnelOpts, name)
}

func (box *CloudBox) Open(templateOpts *model.TemplateOptions, tunnelOpts *model.TunnelOptions) error {
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

func (box *CloudBox) Delete(names []string) ([]model.BoxInfo, error) {
	defer box.close()
	return box.deleteBoxes(names)
}

func (box *CloudBox) Clean() error {
	defer box.close()
	return errors.New("not implemented")
}

func (box *CloudBox) Version() (string, error) {
	defer box.close()
	return box.version()
}
