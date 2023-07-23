package docker

import (
	"testing"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/go-connections/nat"
	"github.com/stretchr/testify/assert"

	"github.com/hckops/hckctl/pkg/box/model"
	"github.com/hckops/hckctl/pkg/client/docker"
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
		Labels: map[string]string{
			"a.b.c": "hello",
			"x.y.z": "world",
		},
	}
	opts := &containerConfigOptions{
		imageName:     "myImageName",
		containerName: "myContainerName",
		ports:         testPorts,
		labels: map[string]string{
			"a.b.c": "hello",
			"x.y.z": "world",
		},
	}

	result, err := buildContainerConfig(opts)
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

func TestBuildNetworkingConfig(t *testing.T) {
	expected := &network.NetworkingConfig{EndpointsConfig: map[string]*network.EndpointSettings{"myNetwork": {NetworkID: "123"}}}

	result := buildNetworkingConfig("myNetwork", "123")
	assert.Equal(t, expected, result)
}

func TestBoxLabel(t *testing.T) {
	expected := "com.hckops.schema.kind=box/v1"
	assert.Equal(t, expected, boxLabel())
}

func TestToBoxDetails(t *testing.T) {
	createdTime, _ := time.Parse(time.RFC3339, "2042-12-08T10:30:05.265113665Z")

	containerDetails := docker.ContainerDetails{
		Info: docker.ContainerInfo{
			ContainerId:   "myId",
			ContainerName: "myName",
			Healthy:       true,
		},
		Created: createdTime,
		Labels: map[string]string{
			"com.hckops.template.local":      "true",
			"com.hckops.template.cache.path": "/tmp/cache/myUuid",
			"com.hckops.box.size":            "m",
		},
		Env: []string{
			"MY_KEY_1=MY_VALUE_1",
			"MY_KEY_2=MY_VALUE_2",
		},
		Ports: []docker.ContainerPort{
			{Local: "123", Remote: "456"},
		},
	}
	expected := &model.BoxDetails{
		Info: model.BoxInfo{
			Id:      "myId",
			Name:    "myName",
			Healthy: true,
		},
		TemplateInfo: &model.BoxTemplateInfo{
			CachedTemplate: &model.CachedTemplateInfo{
				Path: "/tmp/cache/myUuid",
			},
		},
		ProviderInfo: &model.BoxProviderInfo{
			Provider: model.Docker,
			DockerProvider: &model.DockerProviderInfo{
				Network: "TODO",
			},
		},
		Size: model.Medium,
		Env: []model.BoxEnv{
			{Key: "MY_KEY_1", Value: "MY_VALUE_1"},
			{Key: "MY_KEY_2", Value: "MY_VALUE_2"},
		},
		Ports: []model.BoxPort{
			{Alias: "TODO", Local: "123", Remote: "456", Public: false},
		},
		Created: createdTime,
	}
	result, err := toBoxDetails(containerDetails)

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}
