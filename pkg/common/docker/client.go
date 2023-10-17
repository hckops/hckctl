package docker

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"

	"github.com/hckops/hckctl/pkg/client/docker"
	commonModel "github.com/hckops/hckctl/pkg/common/model"
	"github.com/hckops/hckctl/pkg/event"
	"github.com/hckops/hckctl/pkg/schema"
	"github.com/hckops/hckctl/pkg/util"
)

type DockerCommonClient struct {
	client     *docker.DockerClient
	clientOpts *commonModel.DockerOptions
	eventBus   *event.EventBus
}

func NewDockerCommonClient(dockerOpts *commonModel.DockerOptions, eventBus *event.EventBus) (*DockerCommonClient, error) {
	eventBus.Publish(newInitDockerClientEvent())

	dockerClient, err := docker.NewDockerClient()
	if err != nil {
		return nil, errors.Wrap(err, "error docker common client")
	}

	return &DockerCommonClient{
		client:     dockerClient,
		clientOpts: dockerOpts,
		eventBus:   eventBus,
	}, nil
}

func (common *DockerCommonClient) GetClient() *docker.DockerClient {
	return common.client
}

func (common *DockerCommonClient) Close() error {
	common.eventBus.Publish(newCloseDockerClientEvent())
	common.eventBus.Close()
	return common.client.Close()
}

func (common *DockerCommonClient) PullImageOffline(imageName string, onImagePullCallback func()) error {

	imagePullOpts := &docker.ImagePullOpts{
		ImageName:           imageName,
		OnImagePullCallback: onImagePullCallback,
	}
	common.eventBus.Publish(newImagePullDockerEvent(imageName))
	if err := common.client.ImagePull(imagePullOpts); err != nil {
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
	if err := common.client.ImageRemoveDangling(imageRemoveOpts); err != nil {
		return err
	}

	return nil
}

func sidecarLabel() string {
	return fmt.Sprintf("%s=%s", commonModel.LabelSchemaKind, schema.KindSidecarV1.String())
}

func (common *DockerCommonClient) SidecarList(containerName string) ([]commonModel.SidecarInfo, error) {

	// filter by prefix and label
	containers, err := common.client.ContainerList(commonModel.SidecarPrefixName, sidecarLabel())
	if err != nil {
		return nil, err
	}

	var sidecars []commonModel.SidecarInfo
	for _, c := range containers {
		// expect valid name always
		tokens := strings.Split(c.ContainerName, "-")
		// include only associated containers
		if strings.HasSuffix(containerName, tokens[len(tokens)-1]) {
			sidecars = append(sidecars, commonModel.SidecarInfo{Id: c.ContainerId, Name: c.ContainerName})
			// TODO common.eventBus.Publish
		}
	}
	return sidecars, nil
}

func buildSidecarVpnName(containerName string) string {
	// expect valid name always
	tokens := strings.Split(containerName, "-")
	return fmt.Sprintf("%svpn-%s", commonModel.SidecarPrefixName, tokens[len(tokens)-1])
}

func (common *DockerCommonClient) SidecarVpnInject(opts *commonModel.SidecarVpnInjectOpts, portConfig *docker.ContainerPortConfigOpts) (string, error) {

	// sidecarName
	containerName := buildSidecarVpnName(opts.MainContainerId)

	// constants
	imageName := commonModel.SidecarVpnImageName
	// base directory "/usr/share" must exist
	vpnConfigPath := "/usr/share/client.ovpn"

	if err := common.PullImageOffline(imageName, func() {
		common.eventBus.Publish(newSidecarVpnConnectDockerEvent(opts.NetworkVpn.Name))
		common.eventBus.Publish(newSidecarVpnConnectDockerLoaderEvent(opts.NetworkVpn.Name))
	}); err != nil {
		return "", err
	}

	containerConfig, err := docker.BuildContainerConfig(&docker.ContainerConfigOpts{
		ImageName:  imageName,
		Hostname:   opts.MainContainerId,
		Env:        []docker.ContainerEnv{{Key: "OPENVPN_CONFIG", Value: vpnConfigPath}},
		Ports:      portConfig.Ports,
		Tty:        false,
		Entrypoint: nil,
		Cmd:        []string{},
		Labels:     commonModel.NewSidecarLabels().AddSidecarMain(opts.MainContainerId),
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
		CaptureInterrupt: false, // edge case: killing this while creating will leave an orphan sidecar container
		OnContainerCreateCallback: func(containerId string) error {
			// upload openvpn config file
			return common.client.CopyFileToContainer(containerId, opts.NetworkVpn.LocalPath, vpnConfigPath)
		},
		OnContainerStatusCallback: func(status string) {
			common.eventBus.Publish(newSidecarVpnCreateStatusDockerEvent(status))
		},
		OnContainerStartCallback: func() {},
	}
	// sidecarId
	containerId, err := common.client.ContainerCreate(containerOpts)
	if err != nil {
		return "", err
	}
	common.eventBus.Publish(newSidecarVpnCreateDockerEvent(containerName, containerId))
	// block to give time to connect
	util.Sleep(3)

	return containerId, nil
}
