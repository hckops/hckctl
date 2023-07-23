package docker

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
)

func TestNewContainerInfo(t *testing.T) {
	containerInfo := newContainerInfo("myId", "/myName", "running")
	expected := ContainerInfo{
		ContainerId:   "myId",
		ContainerName: "myName",
		Healthy:       true,
	}
	assert.Equal(t, expected, containerInfo)
}

func TestNewContainerDetails(t *testing.T) {
	created := "2042-12-08T10:30:05.265113665Z"
	createdTime, _ := time.Parse(time.RFC3339, created)

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
			Created: created,
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

	expected := ContainerDetails{
		Info: ContainerInfo{
			ContainerId:   "myId",
			ContainerName: "myName",
			Healthy:       false,
		},
		Created: createdTime,
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
	containerDetails, err := newContainerDetails(containerJson)

	assert.NoError(t, err)
	assert.Equal(t, expected, containerDetails)
}
