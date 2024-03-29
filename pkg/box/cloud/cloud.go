package cloud

import (
	"strings"
	"time"

	"github.com/pkg/errors"

	v1 "github.com/hckops/hckctl/pkg/api/v1"
	boxModel "github.com/hckops/hckctl/pkg/box/model"
	"github.com/hckops/hckctl/pkg/client/ssh"
	commonModel "github.com/hckops/hckctl/pkg/common/model"
	"github.com/hckops/hckctl/pkg/util"
)

func newCloudBoxClient(commonOpts *boxModel.CommonBoxOptions, cloudOpts *commonModel.CloudOptions) (*CloudBoxClient, error) {
	commonOpts.EventBus.Publish(newInitCloudClientEvent())

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

// TODO issue "use of closed network connection" after multiple invocation e.g. info + open
func (box *CloudBoxClient) close() error {
	box.eventBus.Publish(newCloseCloudClientEvent())
	box.eventBus.Close()
	return box.client.Close()
}

func (box *CloudBoxClient) createBox(opts *boxModel.CreateOptions) (*boxModel.BoxInfo, error) {
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

	// internal id not exposed
	return &boxModel.BoxInfo{Id: boxName, Name: boxName}, nil
}

func (box *CloudBoxClient) connectBox(opts *boxModel.ConnectOptions) error {

	if opts.DisableExec && opts.DisableTunnel {
		return errors.New("invalid connection options")
	}

	for _, e := range opts.Template.EnvironmentVariables() {
		box.eventBus.Publish(newApiEnvCloudEvent(opts.Name, e))
		box.eventBus.Publish(newApiEnvCloudConsoleEvent(opts.Name, e))
	}

	// tunnel only
	if opts.DisableExec {
		return box.tunnelBox(opts.Template, opts.Name, true)
	}

	if !opts.DisableTunnel {
		if err := box.tunnelBox(opts.Template, opts.Name, false); err != nil {
			return err
		}
	}

	return box.execBox(opts)
}

func (box *CloudBoxClient) execBox(opts *boxModel.ConnectOptions) error {
	box.eventBus.Publish(newApiExecCloudEvent(opts.Name))

	session := v1.NewBoxExecSession(box.clientOpts.Version, opts.Name)
	payload, err := session.Encode()
	if err != nil {
		return errors.Wrap(err, "error cloud exec session")
	}

	// TODO close streams on opts.OnInterruptCallback

	if opts.DeleteOnExit {
		defer box.deleteBoxes([]string{opts.Name})
	}

	execOpts := &ssh.SshExecOpts{
		Payload: payload,
		OnStreamStartCallback: func() {
			// stop loader
			box.eventBus.Publish(newApiStopCloudLoaderEvent())
		},
		OnStreamErrorCallback: func(err error) {
			box.eventBus.Publish(newApiExecErrorCloudEvent(opts.Name, err))
		},
	}
	return box.client.Exec(execOpts)
}

func (box *CloudBoxClient) tunnelBox(template *boxModel.BoxV1, name string, isWait bool) error {

	if !template.HasPorts() {
		box.eventBus.Publish(newApiTunnelIgnoreCloudEvent(name))
		// exit, no service/port available to bind
		return nil
	}

	networkPorts := template.NetworkPortValues(true)
	portPadding := boxModel.PortFormatPadding(networkPorts)

	errorChannel := make(chan struct{})

	for _, p := range networkPorts {
		port, err := bindPort(p)
		if err != nil {
			return err
		}

		// TODO print remote url
		box.eventBus.Publish(newApiTunnelBindingCloudEvent(name, port))
		box.eventBus.Publish(newApiTunnelBindingCloudConsoleEvent(name, port, portPadding))
		box.eventBus.Publish(newApiTunnelListenCloudLoaderEvent())

		sshTunnelOpts := &ssh.SshTunnelOpts{
			LocalPort:  port.Local,
			RemoteHost: name,
			RemotePort: port.Remote,
			OnTunnelStartCallback: func(connection string) {
				box.eventBus.Publish(newApiTunnelStartCloudEvent(name, port, connection))
				// stop loader
				box.eventBus.Publish(newApiStopCloudLoaderEvent())
			},
			OnTunnelStopCallback: func(connection string) {
				box.eventBus.Publish(newApiTunnelStopCloudEvent(name, port, connection))
			},
			OnTunnelErrorCallback: func(err error) {
				box.eventBus.Publish(newApiTunnelErrorCloudEvent(name, err))
				close(errorChannel)
			},
		}
		go box.client.Tunnel(sshTunnelOpts)
	}

	if isWait {
		// waits until it's interrupted
		select {
		case <-errorChannel:
		}
	}

	return nil
}

func bindPort(port boxModel.BoxPort) (boxModel.BoxPort, error) {
	localPort, err := util.FindOpenPort(port.Local)
	if err != nil {
		return boxModel.BoxPort{}, errors.Wrapf(err, "error bind local port %s", port.Local)
	}
	// update actual port
	port.Local = localPort

	return port, nil
}

func (box *CloudBoxClient) describeBox(name string) (*boxModel.BoxDetails, error) {
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

func toBoxDetails(response *v1.Message[v1.BoxDescribeResponseBody]) (*boxModel.BoxDetails, error) {

	size, err := boxModel.ExistResourceSize(response.Body.Size)
	if err != nil {
		return nil, errors.Wrap(err, "error cloud box details size")
	}

	var envs []boxModel.BoxEnv
	for _, env := range response.Body.Env {
		// silently ignore invalid envs
		if key, value, err := util.SplitKeyValue(env); err == nil {
			envs = append(envs, boxModel.BoxEnv{
				Key:   key,
				Value: value,
			})
		}
	}

	var ports []boxModel.BoxPort
	for _, port := range response.Body.Ports {
		items := strings.Split(port, "/")
		// silently ignore invalid ports
		if len(items) == 2 {
			ports = append(ports, boxModel.BoxPort{
				Alias:  items[0],
				Local:  boxModel.BoxPortNone, // runtime only
				Remote: items[1],
				Public: false, // TODO url
			})
		}
	}

	created, err := time.Parse(time.RFC3339, response.Body.Created)
	if err != nil {
		return nil, errors.Wrap(err, "error cloud box details created")
	}

	return &boxModel.BoxDetails{
		Info: boxModel.BoxInfo{
			Id:      response.Body.Id,
			Name:    response.Body.Name,
			Healthy: response.Body.Healthy,
		},
		TemplateInfo: &boxModel.BoxTemplateInfo{
			// TODO valid only if response.Body.Template.Public
			GitTemplate: &commonModel.GitTemplateInfo{
				Url:      response.Body.Template.Url,
				Revision: response.Body.Template.Revision,
				Commit:   response.Body.Template.Commit,
				Name:     response.Body.Template.Name,
			},
		},
		ProviderInfo: &boxModel.BoxProviderInfo{
			Provider: boxModel.Cloud,
		},
		Size:    size,
		Env:     boxModel.SortEnv(envs),
		Ports:   boxModel.SortPorts(ports),
		Created: created,
	}, nil
}

func (box *CloudBoxClient) listBoxes() ([]boxModel.BoxInfo, error) {

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

	var result []boxModel.BoxInfo
	for index, item := range response.Body.Items {
		result = append(result, boxModel.BoxInfo{Id: item.Id, Name: item.Name, Healthy: item.Healthy})
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
