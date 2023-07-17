package cloud

import (
	"github.com/pkg/errors"

	v1 "github.com/hckops/hckctl/pkg/api/v1"
	"github.com/hckops/hckctl/pkg/box/model"
	"github.com/hckops/hckctl/pkg/client/ssh"
)

func newCloudBoxClient(commonOpts *model.CommonBoxOptions, cloudOpts *model.CloudBoxOptions) (*CloudBoxClient, error) {
	commonOpts.EventBus.Publish(newClientInitCloudEvent())

	clientConfig := &ssh.SshClientConfig{
		Address:  cloudOpts.Address,
		Username: cloudOpts.Username,
		Token:    cloudOpts.Token,
	}
	sshClient, err := ssh.NewSshClient(clientConfig)
	if err != nil {
		return nil, errors.Wrap(err, "error cloud box")
	}

	return &CloudBoxClient{
		client:     sshClient,
		clientOpts: cloudOpts,
		eventBus:   commonOpts.EventBus,
	}, nil
}

func (box *CloudBoxClient) close() error {
	box.eventBus.Publish(newClientCloseCloudEvent())
	box.eventBus.Close()
	return box.client.Close()
}

func (box *CloudBoxClient) createBox(opts *model.TemplateOptions) (*model.BoxInfo, error) {
	box.eventBus.Publish(newApiCreateCloudLoaderEvent(box.clientOpts.Address, opts.Template.Name))

	request := v1.NewBoxCreateRequest(box.clientOpts.Version, opts.Template.Name, opts.Size.String())
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

func (box *CloudBoxClient) execBox(template *model.BoxV1, tunnelOpts *model.TunnelOptions, name string) error {
	// TODO event
	box.eventBus.Publish(newApiExecCloudEvent(name))

	session := v1.NewBoxExecSession(box.clientOpts.Version, name)
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

func (box *CloudBoxClient) listBoxes() ([]model.BoxInfo, error) {

	request := v1.NewBoxListRequest(box.clientOpts.Version)
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

func (box *CloudBoxClient) deleteBoxes(names []string) ([]string, error) {

	request := v1.NewBoxDeleteRequest(box.clientOpts.Version, names)
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

	var result []string
	for index, item := range response.Body.Items {
		result = append(result, item.Name)
		box.eventBus.Publish(newApiDeleteCloudEvent(index, item.Name))
	}
	return result, nil
}

func (box *CloudBoxClient) version() (string, error) {

	request := v1.NewPingMessage(box.clientOpts.Version)
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
