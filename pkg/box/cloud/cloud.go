package cloud

import (
	"github.com/pkg/errors"

	"github.com/hckops/hckctl/pkg/box/cloud/api/v1"
	"github.com/hckops/hckctl/pkg/box/model"
	"github.com/hckops/hckctl/pkg/client/ssh"
)

func newCloudBox(internalOpts *model.BoxInternalOpts, sshConfig *ssh.SshClientConfig) (*CloudBox, error) {
	internalOpts.EventBus.Publish(newClientInitCloudEvent())

	sshClient, err := ssh.NewSshClient(sshConfig)
	if err != nil {
		return nil, errors.Wrap(err, "error cloud box")
	}

	return &CloudBox{
		clientVersion: internalOpts.ClientVersion,
		clientConfig:  sshConfig,
		client:        sshClient,
		streams:       internalOpts.Streams,
		eventBus:      internalOpts.EventBus,
	}, nil
}

func (box *CloudBox) close() error {
	box.eventBus.Publish(newClientCloseCloudEvent())
	box.eventBus.Close()
	return box.client.Close()
}

func (box *CloudBox) createBox(template *model.BoxV1) (*model.BoxInfo, error) {
	box.eventBus.Publish(newApiCreateCloudLoaderEvent(box.clientConfig.Address, template.Name))

	request := v1.NewBoxCreateRequest(box.clientVersion, template.Name)
	payload, err := request.Encode()
	if err != nil {
		return nil, errors.Wrap(err, "error cloud create request")
	}
	_, err = box.client.SendRequest(request.Protocol(), payload)

	// TODO box.eventBus.Publish(newContainerCreateDockerEvent(template.Name, containerName, containerId))

	return nil, nil
}
