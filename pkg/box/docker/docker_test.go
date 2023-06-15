package docker

import (
	"testing"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
	"github.com/stretchr/testify/assert"

	"github.com/hckops/hckctl/pkg/box/model"
)

var testPorts = []model.BoxPort{
	{Alias: "aaa", Local: "123", Remote: "123"},
	{Alias: "bbb", Local: "456", Remote: "789"},
}

func TestBuildContainerConfig(t *testing.T) {
	expected := &container.Config{
		Hostname:     "myContainerName",
		Image:        "myImageName",
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		OpenStdin:    true,
		StdinOnce:    true,
		Tty:          true,
		ExposedPorts: nat.PortSet{
			"123/tcp": struct{}{},
			"789/tcp": struct{}{},
		},
	}

	result, err := buildContainerConfig("myImageName", "myContainerName", testPorts)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestBuildHostConfig(t *testing.T) {
	// 1024 is the first user port available https://superuser.com/questions/1631280/what-exactly-are-user-ports
	expected := &container.HostConfig{
		PortBindings: nat.PortMap{
			"123/tcp": []nat.PortBinding{{HostIP: "0.0.0.0", HostPort: "1024"}},
			"789/tcp": []nat.PortBinding{{HostIP: "0.0.0.0", HostPort: "1024"}},
		},
	}

	result, err := buildHostConfig(testPorts, func(port model.BoxPort) {})
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}
