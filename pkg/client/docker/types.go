package docker

import (
	"context"
	"time"

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
