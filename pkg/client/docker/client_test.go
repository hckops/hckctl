package docker

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
)

func TestNewContainerInfo(t *testing.T) {
	containerInfo := NewContainerInfo("myId", "/myName", "running")
	expected := ContainerInfo{
		ContainerId:   "myId",
		ContainerName: "myName",
		Healthy:       true,
	}
	assert.Equal(t, expected, containerInfo)
}

func TestNewContainerDetails(t *testing.T) {

	portBindings := make(nat.PortMap)
	remotePort, _ := nat.NewPort("tcp", "7681")
	portBindings[remotePort] = []nat.PortBinding{{
		HostIP:   "0.0.0.0",
		HostPort: "7683",
	}}

	containerJson := types.ContainerJSON{
		ContainerJSONBase: &types.ContainerJSONBase{
			ID:      "myId",
			Name:    "/myName",
			Created: "yyyy-mm-dd",
			State: &types.ContainerState{
				Status: "exited",
			},
			HostConfig: &container.HostConfig{
				PortBindings: portBindings,
			},
		},
		Config: &container.Config{
			Labels: map[string]string{
				"com.hckops.test": "true",
			},
			Env: []string{
				"MY_KEY=MY_VALUE",
			},
		},
	}
	containerDetails := newContainerDetails(containerJson)

	expected := ContainerDetails{
		Info: ContainerInfo{
			ContainerId:   "myId",
			ContainerName: "myName",
			Healthy:       false,
		},
		Created: "yyyy-mm-dd",
		Labels: map[string]string{
			"com.hckops.test": "true",
		},
		Env: []string{
			"MY_KEY=MY_VALUE",
		},
		Ports: []ContainerPort{
			{Local: "7683", Remote: "7681"},
		},
	}

	assert.Equal(t, expected, containerDetails)
}
