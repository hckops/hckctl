package docker

import (
	"io"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	commonModel "github.com/hckops/hckctl/pkg/common/model"
)

type ContainerConfigOpts struct {
	ImageName string
	Hostname  string
	Env       []ContainerEnv
	Ports     []ContainerPort
	Tty       bool
	Cmd       []string
	Labels    commonModel.Labels
}

type ContainerHostConfigOpts struct {
	NetworkMode        string
	Ports              []ContainerPort
	OnPortBindCallback func(port ContainerPort)
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
	ContainerName                string
	ContainerConfig              *container.Config
	HostConfig                   *container.HostConfig
	NetworkingConfig             *network.NetworkingConfig
	WaitStatus                   bool
	CaptureInterrupt             bool
	OnContainerInterruptCallback func(containerId string)
	OnContainerCreateCallback    func(containerId string) error
	OnContainerWaitCallback      func(containerId string) error
	OnContainerStatusCallback    func(status string)
	OnContainerStartCallback     func()
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
