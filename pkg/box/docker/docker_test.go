package docker

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/hckops/hckctl/pkg/box/model"
	"github.com/hckops/hckctl/pkg/client/docker"
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
		Env: []string{
			"MY_KEY_2=MY_VALUE_2",
			"MY_KEY_1=MY_VALUE_1",
			"MY_KEY_3=MY_VALUE_3",
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
				Network: "myNetworkName",
				Ip:      "myNetworkIp",
			},
		},
		Size: model.Medium,
		Env: []model.BoxEnv{
			{Key: "MY_KEY_1", Value: "MY_VALUE_1"},
			{Key: "MY_KEY_2", Value: "MY_VALUE_2"},
			{Key: "MY_KEY_3", Value: "MY_VALUE_3"},
		},
		Ports: []model.BoxPort{
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
