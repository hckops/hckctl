package cloud

import (
	"github.com/pkg/errors"

	"github.com/hckops/hckctl/pkg/box/model"
	"github.com/hckops/hckctl/pkg/client/ssh"
	"github.com/hckops/hckctl/pkg/event"
)

type CloudBoxClient struct {
	client     *ssh.SshClient
	clientOpts *model.CloudBoxOptions
	eventBus   *event.EventBus
}

func NewCloudBoxClient(commonOpts *model.CommonBoxOptions, cloudOpts *model.CloudBoxOptions) (*CloudBoxClient, error) {
	return newCloudBoxClient(commonOpts, cloudOpts)
}

func (box *CloudBoxClient) Provider() model.BoxProvider {
	return model.Cloud
}

func (box *CloudBoxClient) Events() *event.EventBus {
	return box.eventBus
}

func (box *CloudBoxClient) Create(templateOpts *model.TemplateOptions) (*model.BoxInfo, error) {
	defer box.close()
	return box.createBox(templateOpts)
}

func (box *CloudBoxClient) Connect(template *model.BoxV1, tunnelOpts *model.TunnelOptions, name string) error {
	defer box.close()
	return box.execBox(template, tunnelOpts, name)
}

func (box *CloudBoxClient) Open(templateOpts *model.TemplateOptions, tunnelOpts *model.TunnelOptions) error {
	defer box.close()
	// TODO tunnelOpts
	return errors.New("not implemented")
}

func (box *CloudBoxClient) Copy(string, string, string) error {
	defer box.close()
	return errors.New("not implemented")
}

func (box *CloudBoxClient) List() ([]model.BoxInfo, error) {
	defer box.close()
	return box.listBoxes()
}

func (box *CloudBoxClient) Delete(names []string) ([]model.BoxInfo, error) {
	defer box.close()
	return box.deleteBoxes(names)
}

func (box *CloudBoxClient) Clean() error {
	defer box.close()
	return errors.New("not implemented")
}

func (box *CloudBoxClient) Version() (string, error) {
	defer box.close()
	return box.version()
}
