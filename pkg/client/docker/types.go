package docker

import (
	"context"
	"io"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
)

const (
	ContainerStatusRunning = "running"
)

type DockerClient struct {
	ctx    context.Context
	docker *client.Client
}

type ContainerInfo struct {
	ContainerId   string
	ContainerName string
	Healthy       bool
}

type ContainerDetails struct {
	Info    ContainerInfo
	Created time.Time
	Labels  map[string]string
	Env     []string
	Ports   []ContainerPort
	Network NetworkInfo
}

type NetworkInfo struct {
	Id         string
	Name       string
	IpAddress  string
	MacAddress string
}

type ContainerPort struct {
	Local  string
	Remote string
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

type ContainerRestartOpts struct {
	ContainerId       string
	OnRestartCallback func(string)
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

type ContainerLogsOpts struct {
	ContainerId           string
	OutStream             io.Writer
	ErrStream             io.Writer
	OnStreamCloseCallback func()
	OnStreamErrorCallback func(error)
}
