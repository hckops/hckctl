package cloud

import (
	"github.com/pkg/errors"

	v1 "github.com/hckops/hckctl/pkg/api/v1"
	"github.com/hckops/hckctl/pkg/client/ssh"
	"github.com/hckops/hckctl/pkg/lab/model"
	"github.com/hckops/hckctl/pkg/provider"
)

func newCloudLabClient(commonOpts *model.CommonLabOptions, cloudOpts *provider.CloudOptions) (*CloudLabClient, error) {
	commonOpts.EventBus.Publish(newInitCloudClientEvent())

	clientConfig := &ssh.SshClientConfig{
		Address:  cloudOpts.Address,
		Username: cloudOpts.Username,
		Token:    cloudOpts.Token,
	}
	sshClient, err := ssh.NewSshClient(clientConfig)
	if err != nil {
		return nil, errors.Wrap(err, "error cloud lab")
	}

	return &CloudLabClient{
		client:     sshClient,
		clientOpts: cloudOpts,
		eventBus:   commonOpts.EventBus,
	}, nil
}

func (lab *CloudLabClient) createLab(opts *model.CreateOptions) (*model.LabInfo, error) {
	lab.eventBus.Publish(newApiCreateCloudLoaderEvent(lab.clientOpts.Address, opts.Template.Name))

	request := v1.NewLabCreateRequest(lab.clientOpts.Version, opts.Template.Name, opts.Parameters)
	payload, err := request.Encode()
	if err != nil {
		return nil, errors.Wrap(err, "error cloud lab create request")
	}
	value, err := lab.client.SendRequest(request.Protocol(), payload)
	if err != nil {
		return nil, errors.Wrap(err, "error cloud lab create")
	}

	response, err := v1.Decode[v1.LabCreateResponseBody](value)
	if err != nil {
		return nil, errors.Wrap(err, "error cloud lab create response")
	}
	labName := response.Body.Name
	lab.eventBus.Publish(newApiCreateCloudEvent(opts.Template.Name, labName))

	return &model.LabInfo{Id: labName, Name: labName}, nil
}
