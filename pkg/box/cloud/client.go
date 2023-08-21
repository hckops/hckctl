package cloud

import (
	"github.com/pkg/errors"

	"github.com/hckops/hckctl/pkg/box/model"
	"github.com/hckops/hckctl/pkg/client/ssh"
	"github.com/hckops/hckctl/pkg/event"
	"github.com/hckops/hckctl/pkg/provider"
)

type CloudBoxClient struct {
	client     *ssh.SshClient
	clientOpts *provider.CloudOptions
	eventBus   *event.EventBus
}

func NewCloudBoxClient(commonOpts *model.CommonBoxOptions, cloudOpts *provider.CloudOptions) (*CloudBoxClient, error) {
	return newCloudBoxClient(commonOpts, cloudOpts)
}

func (box *CloudBoxClient) Provider() model.BoxProvider {
	return model.Cloud
}

func (box *CloudBoxClient) Events() *event.EventBus {
	return box.eventBus
}

func (box *CloudBoxClient) Create(opts *model.CreateOptions) (*model.BoxInfo, error) {
	//defer box.close()
	return box.createBox(opts)
}

func (box *CloudBoxClient) Connect(opts *model.ConnectOptions) error {
	//defer box.close()
	return box.connectBox(opts)
}

func (box *CloudBoxClient) Copy(string, string, string) error {
	defer box.close()
	return errors.New("not implemented")
}

func (box *CloudBoxClient) Describe(name string) (*model.BoxDetails, error) {
	//defer box.close()
	return box.describeBox(name)
}

func (box *CloudBoxClient) List() ([]model.BoxInfo, error) {
	defer box.close()
	return box.listBoxes()
}

func (box *CloudBoxClient) Delete(names []string) ([]string, error) {
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
