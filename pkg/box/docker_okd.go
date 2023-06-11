package box

//
//import (
//	"context"
//	"io"
//	"sync"
//
//	"github.com/docker/docker/api/types"
//	"github.com/docker/docker/api/types/container"
//	"github.com/docker/docker/client"
//	"github.com/docker/docker/pkg/stdcopy"
//	"github.com/docker/go-connections/nat"
//	"github.com/pkg/errors"
//
//	"github.com/hckops/hckctl/pkg/template/model"
//	"github.com/hckops/hckctl/pkg/util"
//)
//
//type DockerClient struct {
//	ctx       context.Context
//	dockerApi *client.Client
//	template  *model.BoxV1
//	streams   *BoxStreams
//	eventBus  *EventBus
//}
//
//func NewDockerClient(template *model.BoxV1, streams *BoxStreams, eventBus *EventBus) (*DockerClient, error) {
//
//	dockerApi, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
//	if err != nil {
//		return nil, errors.Wrap(err, "error docker client")
//	}
//
//	return &DockerClient{
//		ctx:       context.Background(),
//		dockerApi: dockerApi,
//		template:  template,
//		streams:   streams,
//		eventBus:  eventBus,
//	}, nil
//}
//
//func (c *DockerClient) Events() *EventBus {
//	return c.eventBus
//}
//
//// TODO
//func (c *DockerClient) Open() error {
//	defer c.close()
//
//	if err := c.setup(); err != nil {
//		return err
//	}
//
//	boxName := c.template.GenerateName()
//	boxId, err := c.createContainer(boxName)
//	if err != nil {
//		return err
//	}
//	c.eventBus.PublishDebugEvent("open", "open box: templateName=%s boxName=%s boxId=%s", c.template.Name, boxName, boxId)
//
//	if err := c.Exec(boxId); err != nil {
//		return err
//	}
//
//	return nil
//}
//
//func (c *DockerClient) Create() (*BoxInfo, error) {
//	defer c.close()
//
//	if err := c.setup(); err != nil {
//		return nil, err
//	}
//
//	boxName := c.template.GenerateName()
//	boxId, err := c.createContainer(boxName)
//	if err != nil {
//		return nil, err
//	}
//	c.eventBus.PublishDebugEvent("create", "create box: templateName=%s boxName=%s boxId=%s", c.template.Name, boxName, boxId)
//
//	return &BoxInfo{Id: boxId, Name: boxName}, nil
//}
//
//func (c *DockerClient) close() error {
//	c.eventBus.PublishDebugEvent("close", "closing event bus and docker client")
//	c.eventBus.Close()
//	return c.dockerApi.Close()
//}
//
//func (c *DockerClient) setup() error {
//	imageName := c.template.ImageName()
//	c.eventBus.PublishDebugEvent("setup", "setup image: imageName=%s", imageName)
//	// TODO delete dangling images
//	reader, err := c.dockerApi.ImagePull(c.ctx, imageName, types.ImagePullOptions{})
//	if err != nil {
//		return errors.Wrap(err, "error image pull")
//	}
//	defer reader.Close()
//
//	c.eventBus.PublishInfoEvent("setup", "pulling %s", imageName)
//
//	// suppress default output
//	if _, err := io.Copy(io.Discard, reader); err != nil {
//		return errors.Wrap(err, "error image pull output message")
//	}
//
//	return nil
//}
//
//func (c *DockerClient) createContainer(containerName string) (string, error) {
//	c.eventBus.PublishDebugEvent("createContainer", "create container: containerName=%s", containerName)
//
//	containerConfig, err := buildContainerConfig(
//		c.template.ImageName(),
//		containerName,
//		c.template.NetworkPorts(),
//	)
//	if err != nil {
//		return "", err
//	}
//
//	onPortBindCallback := func(port model.BoxPort) {
//		c.eventBus.PublishConsoleEvent("createContainer",
//			"[%s][%s]   \texpose (container) %s -> (local) http://localhost:%s",
//			containerName, port.Alias, port.Remote, port.Local)
//	}
//
//	hostConfig, err := buildHostConfig(c.template.NetworkPorts(), onPortBindCallback)
//	if err != nil {
//		return "", err
//	}
//
//	newContainer, err := c.dockerApi.ContainerCreate(
//		c.ctx,
//		containerConfig,
//		hostConfig,
//		nil, // networkingConfig
//		nil, // platform
//		containerName)
//	if err != nil {
//		return "", errors.Wrap(err, "error container create")
//	}
//
//	if err := c.dockerApi.ContainerStart(c.ctx, newContainer.ID, types.ContainerStartOptions{}); err != nil {
//		return "", errors.Wrap(err, "error container start")
//	}
//
//	return newContainer.ID, nil
//}
//
//func buildContainerConfig(imageName string, containerName string, ports []model.BoxPort) (*container.Config, error) {
//
//	exposedPorts := make(nat.PortSet)
//	for _, port := range ports {
//		p, err := nat.NewPort("tcp", port.Remote)
//		if err != nil {
//			return nil, errors.Wrap(err, "error docker port: containerConfig")
//		}
//		exposedPorts[p] = struct{}{}
//	}
//
//	return &container.Config{
//		Hostname:     containerName,
//		Image:        imageName,
//		AttachStdin:  true,
//		AttachStdout: true,
//		AttachStderr: true,
//		OpenStdin:    true,
//		StdinOnce:    true,
//		Tty:          true,
//		ExposedPorts: exposedPorts,
//	}, nil
//}
//
//func buildHostConfig(ports []model.BoxPort, onPortBindCallback func(port model.BoxPort)) (*container.HostConfig, error) {
//
//	portBindings := make(nat.PortMap)
//	for _, port := range ports {
//
//		localPort, err := util.FindOpenPort(port.Local)
//		if err != nil {
//			return nil, errors.Wrap(err, "error docker local port: hostConfig")
//		}
//
//		remotePort, err := nat.NewPort("tcp", port.Remote)
//		if err != nil {
//			return nil, errors.Wrap(err, "error docker remote port: hostConfig")
//		}
//
//		// actual binded port
//		onPortBindCallback(model.BoxPort{
//			Alias:  port.Alias,
//			Local:  localPort,
//			Remote: port.Remote,
//		})
//
//		portBindings[remotePort] = []nat.PortBinding{{
//			HostIP:   "0.0.0.0",
//			HostPort: localPort,
//		}}
//	}
//
//	return &container.HostConfig{
//		PortBindings: portBindings,
//	}, nil
//}
//
//// TODO close
//func (c *DockerClient) Exec(boxId string) error {
//	return c.execContainer(boxId)
//}
//
//func (c *DockerClient) execContainer(containerId string) error {
//	c.eventBus.PublishDebugEvent("execContainer", "exec container: containerId=%s", containerId)
//
//	execCreateResponse, err := c.dockerApi.ContainerExecCreate(c.ctx, containerId, types.ExecConfig{
//		AttachStdin:  true,
//		AttachStdout: true,
//		AttachStderr: true,
//		Detach:       false,
//		Tty:          c.streams.IsTty,
//		Cmd:          []string{c.template.Shell},
//	})
//	if err != nil {
//		return errors.Wrap(err, "error container exec create")
//	}
//
//	execAttachResponse, err := c.dockerApi.ContainerExecAttach(c.ctx, execCreateResponse.ID, types.ExecStartCheck{
//		Tty: c.streams.IsTty,
//	})
//	if err != nil {
//		return errors.Wrap(err, "error container exec attach")
//	}
//	defer execAttachResponse.Close()
//
//	removeContainerCallback := func() {
//		if err := c.removeContainer(containerId); err != nil {
//			c.eventBus.PublishDebugEvent("removeContainer", "error removing container: containerId=%s error=%v", containerId, err)
//		}
//	}
//
//	onStreamErrorCallback := func(err error, message string) {
//		c.eventBus.PublishErrorEvent("execContainer", message)
//	}
//
//	handleStreams(&execAttachResponse, c.streams, removeContainerCallback, onStreamErrorCallback)
//
//	// fixes echoes and handle SIGTERM interrupt properly
//	if terminal, err := util.NewRawTerminal(c.streams.In); err == nil {
//		defer terminal.Restore()
//	}
//
//	c.eventBus.PublishConsoleEvent("execContainer", "waiting")
//
//	// waits for interrupt signals
//	statusCh, errCh := c.dockerApi.ContainerWait(c.ctx, containerId, container.WaitConditionNotRunning)
//	select {
//	case err := <-errCh:
//		if err != nil {
//			return errors.Wrap(err, "error container wait")
//		}
//	case <-statusCh:
//	}
//	return nil
//}
//
//func (c *DockerClient) removeContainer(containerId string) error {
//	c.eventBus.PublishDebugEvent("removeContainer", "removing container: containerId=%s", containerId)
//
//	if err := c.dockerApi.ContainerRemove(c.ctx, containerId, types.ContainerRemoveOptions{Force: true}); err != nil {
//		return errors.Wrap(err, "error docker remove")
//	}
//	return nil
//}
//
//func handleStreams(
//	execAttachResponse *types.HijackedResponse,
//	streams *BoxStreams,
//	onCloseCallback func(),
//	onStreamErrorCallback func(error, string),
//) {
//
//	var once sync.Once
//	go func() {
//
//		if streams.IsTty {
//			if _, err := io.Copy(streams.Out, execAttachResponse.Reader); err != nil {
//				onStreamErrorCallback(err, "error copy stdout docker->local")
//			}
//		} else {
//			if _, err := stdcopy.StdCopy(streams.Out, streams.Err, execAttachResponse.Reader); err != nil {
//				onStreamErrorCallback(err, "error copy stdout and stderr docker->local")
//			}
//		}
//
//		once.Do(onCloseCallback)
//	}()
//	go func() {
//		if _, err := io.Copy(execAttachResponse.Conn, streams.In); err != nil {
//			onStreamErrorCallback(err, "error copy stdin local->docker")
//		}
//
//		once.Do(onCloseCallback)
//	}()
//}
