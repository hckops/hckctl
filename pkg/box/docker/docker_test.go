package docker

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	boxModel "github.com/hckops/hckctl/pkg/box/model"
	"github.com/hckops/hckctl/pkg/client/docker"
	commonModel "github.com/hckops/hckctl/pkg/common/model"
)

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
		Env: []docker.ContainerEnv{
			{Key: "MY_KEY_2", Value: "MY_VALUE_2"},
			{Key: "MY_KEY_1", Value: "MY_VALUE_1"},
			{Key: "MY_KEY_3", Value: "MY_VALUE_3"},
		},
		Ports: []docker.ContainerPort{
			{Local: "local-x", Remote: "remote-2"},
			{Local: "local-y", Remote: "remote-1"},
			{Local: "local-z", Remote: "remote-3"},
		},
		Network: docker.NetworkInfo{
			Name:      "myNetworkName",
			IpAddress: "myNetworkIp",
		},
	}
	expected := &boxModel.BoxDetails{
		Info: boxModel.BoxInfo{
			Id:      "myId",
			Name:    "myName",
			Healthy: true,
		},
		TemplateInfo: &boxModel.BoxTemplateInfo{
			CachedTemplate: &commonModel.CachedTemplateInfo{
				Path: "/tmp/cache/myUuid",
			},
		},
		ProviderInfo: &boxModel.BoxProviderInfo{
			Provider: boxModel.Docker,
			DockerProvider: &commonModel.DockerProviderInfo{
				Network: "myNetworkName",
				Ip:      "myNetworkIp",
			},
		},
		Size: boxModel.Medium,
		Env: []boxModel.BoxEnv{
			{Key: "MY_KEY_1", Value: "MY_VALUE_1"},
			{Key: "MY_KEY_2", Value: "MY_VALUE_2"},
			{Key: "MY_KEY_3", Value: "MY_VALUE_3"},
		},
		Ports: []boxModel.BoxPort{
			{Alias: "none", Local: "local-y", Remote: "remote-1", Public: false},
			{Alias: "none", Local: "local-x", Remote: "remote-2", Public: false},
			{Alias: "none", Local: "local-z", Remote: "remote-3", Public: false},
		},
		Created: createdTime,
	}
	result, err := toBoxDetails(containerDetails)

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}
