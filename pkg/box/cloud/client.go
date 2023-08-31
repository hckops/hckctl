package cloud

import (
	"github.com/pkg/errors"

	boxModel "github.com/hckops/hckctl/pkg/box/model"
	"github.com/hckops/hckctl/pkg/client/ssh"
	commonModel "github.com/hckops/hckctl/pkg/common/model"
	"github.com/hckops/hckctl/pkg/event"
)

type CloudBoxClient struct {
	client     *ssh.SshClient
	clientOpts *commonModel.CloudOptions
	eventBus   *event.EventBus
}

func NewCloudBoxClient(commonOpts *boxModel.CommonBoxOptions, cloudOpts *commonModel.CloudOptions) (*CloudBoxClient, error) {
	return newCloudBoxClient(commonOpts, cloudOpts)
}

func (box *CloudBoxClient) Provider() boxModel.BoxProvider {
	return boxModel.Cloud
}

func (box *CloudBoxClient) Events() *event.EventBus {
	return box.eventBus
}

func (box *CloudBoxClient) Create(opts *boxModel.CreateOptions) (*boxModel.BoxInfo, error) {
	//defer box.close()
	return box.createBox(opts)
}

func (box *CloudBoxClient) Connect(opts *boxModel.ConnectOptions) error {
	//defer box.close()
	return box.connectBox(opts)
}

func (box *CloudBoxClient) Copy(string, string, string) error {
	defer box.close()
	return errors.New("not implemented")
}

func (box *CloudBoxClient) Describe(name string) (*boxModel.BoxDetails, error) {
	//defer box.close()
	return box.describeBox(name)
}

func (box *CloudBoxClient) List() ([]boxModel.BoxInfo, error) {
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
