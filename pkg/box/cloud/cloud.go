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
	value, err := box.client.SendRequest(request.Protocol(), payload)
	if err != nil {
		return nil, errors.Wrap(err, "error cloud create")
	}

	response, err := v1.Decode[v1.BoxCreateResponseBody](value)
	if err != nil {
		return nil, errors.Wrap(err, "error cloud create response")
	}
	boxName := response.Body.Name
	box.eventBus.Publish(newApiCreateCloudEvent(template.Name, boxName))

	return &model.BoxInfo{Id: boxName, Name: boxName}, nil
}

// empty "names" means all
func (box *CloudBox) deleteBoxes(names []string) ([]model.BoxInfo, error) {
	// TODO box.eventBus.Publish

	request := v1.NewBoxDeleteRequest(box.clientVersion, names)
	payload, err := request.Encode()
	if err != nil {
		return nil, errors.Wrap(err, "error cloud delete request")
	}
	value, err := box.client.SendRequest(request.Protocol(), payload)
	if err != nil {
		return nil, errors.Wrap(err, "error cloud delete")
	}

	response, err := v1.Decode[v1.BoxDeleteResponseBody](value)
	if err != nil {
		return nil, errors.Wrap(err, "error cloud delete response")
	}
	var result []model.BoxInfo
	for _, name := range response.Body.Names {
		result = append(result, model.BoxInfo{Id: name, Name: name})
	}

	return result, nil
}
