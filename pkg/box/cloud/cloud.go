package cloud

import (
	"github.com/pkg/errors"

	"github.com/hckops/hckctl/pkg/box/cloud/api/v1"
	"github.com/hckops/hckctl/pkg/box/model"
	"github.com/hckops/hckctl/pkg/client/ssh"
)

func newCloudBox(internalOpts *model.BoxInternalOptions, clientConfig *ssh.SshClientConfig) (*CloudBox, error) {
	internalOpts.EventBus.Publish(newClientInitCloudEvent())

	sshClient, err := ssh.NewSshClient(clientConfig)
	if err != nil {
		return nil, errors.Wrap(err, "error cloud box")
	}

	return &CloudBox{
		clientVersion: internalOpts.ClientVersion,
		clientConfig:  clientConfig,
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

func (box *CloudBox) execBox(template *model.BoxV1, name string) error {
	// TODO box.eventBus.Publish

	request := v1.NewBoxExecRequest(box.clientVersion, name)
	payload, err := request.Encode()
	if err != nil {
		return errors.Wrap(err, "error cloud exec request")
	}

	opts := &ssh.ExecOpts{
		Payload: payload, // TODO BoxExecRequestBody
		OnStreamStartCallback: func() {
			// TODO
		},
		OnStreamErrorCallback: func(err error) {
			// TODO
		},
	}
	return box.client.Exec(opts)
}

// empty "names" means all
func (box *CloudBox) deleteBoxes(names []string) ([]model.BoxInfo, error) {
	// TODO box.eventBus.Publish
	// TODO delete namespace if empty

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
	return toBoxes(response.Body.Names), nil
}

func (box *CloudBox) listBoxes() ([]model.BoxInfo, error) {
	// TODO box.eventBus.Publish

	request := v1.NewBoxListRequest(box.clientVersion)
	payload, err := request.Encode()
	if err != nil {
		return nil, errors.Wrap(err, "error cloud list request")
	}
	value, err := box.client.SendRequest(request.Protocol(), payload)
	if err != nil {
		return nil, errors.Wrap(err, "error cloud list")
	}

	response, err := v1.Decode[v1.BoxListResponseBody](value)
	if err != nil {
		return nil, errors.Wrap(err, "error cloud delete response")
	}
	return toBoxes(response.Body.Names), nil
}

func toBoxes(names []string) []model.BoxInfo {
	var result []model.BoxInfo
	for _, name := range names {
		result = append(result, model.BoxInfo{Id: name, Name: name})
		// TODO box.eventBus.Publish + index
	}
	return result
}
