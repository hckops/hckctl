package cloud

import (
	"github.com/pkg/errors"

	v1 "github.com/hckops/hckctl/pkg/api/v1"
	"github.com/hckops/hckctl/pkg/box/model"
	"github.com/hckops/hckctl/pkg/client/ssh"
)

func newCloudBox(commonOpts *model.BoxCommonOptions, clientConfig *ssh.SshClientConfig) (*CloudBox, error) {
	commonOpts.EventBus.Publish(newClientInitCloudEvent())

	sshClient, err := ssh.NewSshClient(clientConfig)
	if err != nil {
		return nil, errors.Wrap(err, "error cloud box")
	}

	return &CloudBox{
		clientConfig: clientConfig,
		client:       sshClient,
		eventBus:     commonOpts.EventBus,
	}, nil
}

func (box *CloudBox) close() error {
	box.eventBus.Publish(newClientCloseCloudEvent())
	box.eventBus.Close()
	return box.client.Close()
}

func (box *CloudBox) createBox(opts *model.TemplateOptions) (*model.BoxInfo, error) {
	box.eventBus.Publish(newApiCreateCloudLoaderEvent(box.clientConfig.Address, opts.Template.Name))

	request := v1.NewBoxCreateRequest(box.clientConfig.Version, opts.Template.Name, opts.Size.String())
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
	box.eventBus.Publish(newApiCreateCloudEvent(opts.Template.Name, boxName, response.Body.Size))

	return &model.BoxInfo{Id: boxName, Name: boxName}, nil
}

func (box *CloudBox) execBox(template *model.BoxV1, tunnelOpts *model.TunnelOptions, name string) error {
	// TODO event
	box.eventBus.Publish(newApiExecCloudEvent(name))

	session := v1.NewBoxExecSession(box.clientConfig.Version, name)
	payload, err := session.Encode()
	if err != nil {
		return errors.Wrap(err, "error cloud exec session")
	}

	opts := &ssh.SshExecOpts{
		Payload: payload,
		OnStreamStartCallback: func() {
			// TODO stop loader
		},
		OnStreamErrorCallback: func(err error) {
			// TODO stop loader
		},
	}
	return box.client.Exec(opts)
}

func (box *CloudBox) listBoxes() ([]model.BoxInfo, error) {

	request := v1.NewBoxListRequest(box.clientConfig.Version)
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
		return nil, errors.Wrap(err, "error cloud list response")
	}

	var result []model.BoxInfo
	for index, item := range response.Body.Items {
		result = append(result, model.BoxInfo{Id: item.Id, Name: item.Name, Healthy: item.Healthy})
		box.eventBus.Publish(newApiListCloudEvent(index, item.Name))
	}
	return result, nil
}

func (box *CloudBox) deleteBoxes(names []string) ([]model.BoxInfo, error) {

	request := v1.NewBoxDeleteRequest(box.clientConfig.Version, names)
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
	for index, item := range response.Body.Items {
		result = append(result, model.BoxInfo{Id: item.Id, Name: item.Name, Healthy: false})
		box.eventBus.Publish(newApiDeleteCloudEvent(index, item.Name))
	}
	return result, nil
}

func (box *CloudBox) version() (string, error) {

	request := v1.NewPingMessage(box.clientConfig.Version)
	payload, err := request.Encode()
	box.eventBus.Publish(newApiRawCloudEvent(payload))
	if err != nil {
		return "", errors.Wrap(err, "error cloud ping request")
	}

	value, err := box.client.SendRequest(request.Protocol(), payload)
	if err != nil {
		return "", errors.Wrap(err, "error cloud ping")
	}

	response, err := v1.Decode[v1.PongBody](value)
	if err != nil {
		return "", errors.Wrap(err, "error cloud pong response")
	}
	box.eventBus.Publish(newApiRawCloudEvent(value))

	return response.Origin, nil
}
