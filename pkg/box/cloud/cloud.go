package cloud

import (
	"github.com/hckops/hckctl/pkg/util"
	"strings"
	"time"

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

// TODO common (no exec info)
func (box *CloudBoxClient) openBox(templateOpts *model.TemplateOptions, tunnelOpts *model.TunnelOptions) error {
	if info, err := box.createBox(templateOpts); err != nil {
		return err
	} else {
		return box.execBox(templateOpts.Template, tunnelOpts, info.Name, true)
	}
}

func (box *CloudBoxClient) execBox(template *model.BoxV1, tunnelOpts *model.TunnelOptions, name string, removeOnExit bool) error {
	box.eventBus.Publish(newApiExecCloudEvent(name))

	// TODO print environment variables

	if tunnelOpts.TunnelOnly {
		return box.tunnelBox(template, name)
	}

	if !tunnelOpts.NoTunnel {
		if err := box.tunnelBox(template, name); err != nil {
			return err
		}
	}

	session := v1.NewBoxExecSession(box.clientOpts.Version, name)
	payload, err := session.Encode()
	if err != nil {
		return errors.Wrap(err, "error cloud exec session")
	}

	opts := &ssh.SshExecOpts{
		Payload: payload,
		OnStreamStartCallback: func() {
			// stop loader
			box.eventBus.Publish(newApiExecCloudLoaderEvent())
		},
		OnStreamErrorCallback: func(err error) {
			box.eventBus.Publish(newApiExecErrorCloudEvent(name, err))
		},
	}

	if removeOnExit {
		defer box.deleteBoxes([]string{name})
	}

	return box.client.Exec(opts)
}

func (box *CloudBoxClient) tunnelBox(template *model.BoxV1, name string) error {

	if !template.HasPorts() {
		box.eventBus.Publish(newApiTunnelSkippedCloudEvent(name))
		// exit, no service/port available to bind
		return nil
	}

	networkPorts := template.NetworkPorts(true)
	portPadding := model.PortFormatPadding(networkPorts)

	for _, p := range networkPorts {
		port, err := bindPort(p)
		if err != nil {
			return err
		}

		// TODO print remote url
		box.eventBus.Publish(newApiTunnelBindingCloudEvent(name, port))
		box.eventBus.Publish(newApiTunnelBindingCloudConsoleEvent(name, port, portPadding))

		sshTunnelOpts := &ssh.SshTunnelOpts{
			LocalPort:  port.Local,
			RemoteHost: name,
			RemotePort: port.Remote,
			OnTunnelErrorCallback: func(err error) {
				box.eventBus.Publish(newApiTunnelErrorCloudEvent(name, err))
			},
		}
		go box.client.Tunnel(sshTunnelOpts)
	}
	return nil
}

func bindPort(port model.BoxPort) (model.BoxPort, error) {
	localPort, err := util.FindOpenPort(port.Local)
	if err != nil {
		return model.BoxPort{}, errors.Wrapf(err, "error bind local port %s", port.Local)
	}
	// update actual port
	port.Local = localPort

	return port, nil
}

func (box *CloudBoxClient) describe(name string) (*model.BoxDetails, error) {
	box.eventBus.Publish(newApiDescribeCloudEvent(name))

	request := v1.NewBoxDescribeRequest(box.clientOpts.Version, name)
	payload, err := request.Encode()
	if err != nil {
		return nil, errors.Wrap(err, "error cloud describe request")
	}
	value, err := box.client.SendRequest(request.Protocol(), payload)
	if err != nil {
		return nil, errors.Wrap(err, "error cloud describe")
	}

	response, err := v1.Decode[v1.BoxDescribeResponseBody](value)
	if err != nil {
		return nil, errors.Wrap(err, "error cloud describe response")
	}

	return toBoxDetails(response)
}

func toBoxDetails(response *v1.Message[v1.BoxDescribeResponseBody]) (*model.BoxDetails, error) {

	size, err := model.ExistResourceSize(response.Body.Size)
	if err != nil {
		return nil, errors.Wrap(err, "error cloud box details size")
	}

	var env []model.BoxEnv
	for _, e := range response.Body.Env {
		items := strings.Split(e, "=")
		// silently ignore invalid
		if len(items) == 2 {
			env = append(env, model.BoxEnv{
				Key:   items[0],
				Value: items[1],
			})
		}
	}

	var ports []model.BoxPort
	for _, p := range response.Body.Ports {
		items := strings.Split(p, "/")
		// silently ignore invalid
		if len(items) == 2 {
			ports = append(ports, model.BoxPort{
				Alias:  items[0],
				Local:  "TODO", // TODO runtime only
				Remote: items[1],
				Public: false, // TODO
			})
		}
	}

	created, err := time.Parse(time.RFC3339, response.Body.Created)
	if err != nil {
		return nil, errors.Wrap(err, "error cloud box details created")
	}

	return &model.BoxDetails{
		Info: model.BoxInfo{
			Id:      response.Body.Id,
			Name:    response.Body.Name,
			Healthy: response.Body.Healthy,
		},
		TemplateInfo: &model.BoxTemplateInfo{
			// TODO only if response.Body.Template.Public
			GitTemplate: &model.GitTemplateInfo{
				Url:      response.Body.Template.Url,
				Revision: response.Body.Template.Revision,
				Commit:   response.Body.Template.Commit,
				Name:     response.Body.Template.Name,
			},
		},
		ProviderInfo: &model.BoxProviderInfo{
			Provider: model.Cloud,
		},
		Size:    size,
		Env:     env,
		Ports:   ports,
		Created: created,
	}, nil
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
	for index, name := range response.Body.Names {
		result = append(result, name)
		box.eventBus.Publish(newApiDeleteCloudEvent(index, name))
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
