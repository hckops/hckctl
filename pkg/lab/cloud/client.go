package cloud

import (
	"github.com/hckops/hckctl/pkg/client/ssh"
	commonModel "github.com/hckops/hckctl/pkg/common/model"
	"github.com/hckops/hckctl/pkg/event"
	labModel "github.com/hckops/hckctl/pkg/lab/model"
)

type CloudLabClient struct {
	client     *ssh.SshClient
	clientOpts *commonModel.CloudOptions
	eventBus   *event.EventBus
}

func NewCloudLabClient(commonOpts *labModel.CommonLabOptions, cloudOpts *commonModel.CloudOptions) (*CloudLabClient, error) {
	return newCloudLabClient(commonOpts, cloudOpts)
}

func (lab *CloudLabClient) Provider() labModel.LabProvider {
	return labModel.Cloud
}

func (lab *CloudLabClient) Events() *event.EventBus {
	return lab.eventBus
}

func (lab *CloudLabClient) Create(opts *labModel.CreateOptions) (*labModel.LabInfo, error) {
	return lab.createLab(opts)
}
