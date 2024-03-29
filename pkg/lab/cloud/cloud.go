package cloud

import (
	"github.com/pkg/errors"

	v1 "github.com/hckops/hckctl/pkg/api/v1"
	"github.com/hckops/hckctl/pkg/client/ssh"
	commonModel "github.com/hckops/hckctl/pkg/common/model"
	labModel "github.com/hckops/hckctl/pkg/lab/model"
)

func newCloudLabClient(commonOpts *labModel.CommonLabOptions, cloudOpts *commonModel.CloudOptions) (*CloudLabClient, error) {
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

func (lab *CloudLabClient) createLab(opts *labModel.CreateOptions) (*labModel.LabInfo, error) {
	lab.eventBus.Publish(newApiCreateCloudLoaderEvent(lab.clientOpts.Address, opts.LabTemplate.Name))

	request := v1.NewLabCreateRequest(lab.clientOpts.Version, opts.LabTemplate.Name, opts.Parameters)
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
	lab.eventBus.Publish(newApiCreateCloudEvent(opts.LabTemplate.Name, labName))

	return &labModel.LabInfo{Id: labName, Name: labName}, nil
}
