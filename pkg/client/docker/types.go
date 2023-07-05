package docker

import (
	"context"
	"io"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
)

type DockerClient struct {
	ctx    context.Context
	docker *client.Client
}

type DockerClientConfig struct {
	NetworkName string
}

type ImagePullOpts struct {
	ImageName           string
	OnImagePullCallback func()
}

type ImageRemoveOpts struct {
	OnImageRemoveCallback      func(imageId string)
	OnImageRemoveErrorCallback func(imageId string, err error)
}

type ContainerCreateOpts struct {
	ContainerName    string
	ContainerConfig  *container.Config
	HostConfig       *container.HostConfig
	NetworkingConfig *network.NetworkingConfig
}

type ContainerExecOpts struct {
	ContainerId             string
	Shell                   string
	InStream                io.ReadCloser
	OutStream               io.Writer
	ErrStream               io.Writer
	IsTty                   bool
	OnContainerExecCallback func()
	OnStreamCloseCallback   func()
	OnStreamErrorCallback   func(error)
}

type ContainerInfo struct {
	ContainerId   string
	ContainerName string
}

type ContainerLogsOpts struct {
	ContainerId           string
	OutStream             io.Writer
	ErrStream             io.Writer
	OnStreamCloseCallback func()
	OnStreamErrorCallback func(error)
}
