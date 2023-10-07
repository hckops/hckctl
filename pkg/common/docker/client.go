package docker

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"

	"github.com/hckops/hckctl/pkg/client/docker"
	commonModel "github.com/hckops/hckctl/pkg/common/model"
	"github.com/hckops/hckctl/pkg/event"
)

type DockerCommonClient struct {
	Client     *docker.DockerClient
	clientOpts *commonModel.DockerOptions
	eventBus   *event.EventBus
}

func NewDockerCommonClient(eventBus *event.EventBus, dockerOpts *commonModel.DockerOptions) (*DockerCommonClient, error) {
	eventBus.Publish(newInitDockerClientEvent())

	dockerClient, err := docker.NewDockerClient()
	if err != nil {
		return nil, errors.Wrap(err, "error docker common client")
	}

	return &DockerCommonClient{
		Client:     dockerClient,
		clientOpts: dockerOpts,
		eventBus:   eventBus,
	}, nil
}

func (common *DockerCommonClient) Close() error {
	common.eventBus.Publish(newCloseDockerClientEvent())
	common.eventBus.Close()
	return common.Client.Close()
}

func (common *DockerCommonClient) PullImageOffline(imageName string, onImagePullCallback func()) error {

	imagePullOpts := &docker.ImagePullOpts{
		ImageName:           imageName,
		OnImagePullCallback: onImagePullCallback,
	}
	common.eventBus.Publish(newImagePullDockerEvent(imageName))
	if err := common.Client.ImagePull(imagePullOpts); err != nil {
		// ignore error and try to use an existing image if exists
		if common.clientOpts.IgnoreImagePullError {
			common.eventBus.Publish(newImagePullIgnoreDockerEvent(imageName))
		} else {
			// do not allow offline
			return err
		}
	}

	// cleanup obsolete nightly images
	imageRemoveOpts := &docker.ImageRemoveOpts{
		OnImageRemoveCallback: func(imageId string) {
			common.eventBus.Publish(newImageRemoveDockerEvent(imageId))
		},
		OnImageRemoveErrorCallback: func(imageId string, err error) {
			// ignore error: keep images used by existing containers
			common.eventBus.Publish(newImageRemoveIgnoreDockerEvent(imageId, err))
		},
	}
	if err := common.Client.ImageRemoveDangling(imageRemoveOpts); err != nil {
		return err
	}

	return nil
}

func buildVpnSidecarName(containerName string) string {
	// expect valid name always
	tokens := strings.Split(containerName, "-")
	return fmt.Sprintf("sidecar-vpn-%s", tokens[len(tokens)-1])
}

type StartVpnSidecarOptions struct {
	MainContainerName  string
	VpnInfo            *commonModel.VpnNetworkInfo
	ContainerPorts     []docker.ContainerPort
	OnPortBindCallback func(port docker.ContainerPort)
}

func (common *DockerCommonClient) StartVpnSidecar(baseContainerName string, vpnInfo *commonModel.VpnNetworkInfo, portConfig *docker.ContainerPortConfigOpts) (string, error) {

	// sidecarName
	containerName := buildVpnSidecarName(baseContainerName)

	// constants
	imageName := commonModel.SidecarVpnImageName
	// base directory "/usr/share" must exist
	vpnConfigPath := "/usr/share/client.ovpn"

	if err := common.PullImageOffline(imageName, func() {
		common.eventBus.Publish(newVpnSidecarConnectDockerEvent(vpnInfo.Name))
		common.eventBus.Publish(newVpnSidecarConnectDockerLoaderEvent(vpnInfo.Name))
	}); err != nil {
		return "", err
	}

	containerConfig, err := docker.BuildContainerConfig(&docker.ContainerConfigOpts{
		ImageName: imageName,
		Hostname:  baseContainerName,
		Env:       []docker.ContainerEnv{{Key: "OPENVPN_CONFIG", Value: vpnConfigPath}},
		Ports:     portConfig.Ports,
		Tty:       false,
		Cmd:       []string{},
		Labels:    commonModel.Labels{},
	})
	if err != nil {
		return "", err
	}

	hostConfig, err := docker.BuildVpnHostConfig(portConfig)
	if err != nil {
		return "", err
	}

	containerOpts := &docker.ContainerCreateOpts{
		ContainerName:    containerName,
		ContainerConfig:  containerConfig,
		HostConfig:       hostConfig,
		WaitStatus:       false,
		CaptureInterrupt: false, // edge case: killing this while creating will produce an orphan sidecar container
		OnContainerCreateCallback: func(containerId string) error {
			// upload openvpn config file
			return common.Client.CopyFileToContainer(containerId, vpnInfo.LocalPath, vpnConfigPath)
		},
		OnContainerStatusCallback: func(status string) {
			common.eventBus.Publish(newVpnSidecarCreateStatusDockerEvent(status))
		},
		OnContainerStartCallback: func() {},
	}
	// sidecarId
	containerId, err := common.Client.ContainerCreate(containerOpts)
	if err != nil {
		return "", err
	}
	common.eventBus.Publish(newVpnSidecarCreateDockerEvent(containerName, containerId))

	return containerId, nil
}
