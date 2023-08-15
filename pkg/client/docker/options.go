package docker

import (
	"io"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
)

type ImagePullOpts struct {
	ImageName           string
	OnImagePullCallback func()
}

type ImageRemoveOpts struct {
	OnImageRemoveCallback      func(imageId string)
	OnImageRemoveErrorCallback func(imageId string, err error)
}

type ContainerCreateOpts struct {
	ContainerName            string
	ContainerConfig          *container.Config
	HostConfig               *container.HostConfig
	NetworkingConfig         *network.NetworkingConfig
	OnContainerStartCallback func()
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
