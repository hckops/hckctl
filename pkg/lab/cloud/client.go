package cloud

import (
	"github.com/hckops/hckctl/pkg/client/ssh"
	"github.com/hckops/hckctl/pkg/event"
	"github.com/hckops/hckctl/pkg/lab/model"
	"github.com/hckops/hckctl/pkg/provider"
)

type CloudLabClient struct {
	client     *ssh.SshClient
	clientOpts *provider.CloudOptions
	eventBus   *event.EventBus
}

func NewCloudLabClient(commonOpts *model.CommonLabOptions, cloudOpts *provider.CloudOptions) (*CloudLabClient, error) {
	return newCloudLabClient(commonOpts, cloudOpts)
}

func (lab *CloudLabClient) Provider() model.LabProvider {
	return model.Cloud
}

func (lab *CloudLabClient) Events() *event.EventBus {
	return lab.eventBus
}

func (lab *CloudLabClient) Create(opts *model.CreateOptions) (*model.LabInfo, error) {
	return lab.createLab(opts)
}
