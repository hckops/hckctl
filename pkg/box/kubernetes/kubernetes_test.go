package kubernetes

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/hckops/hckctl/pkg/box/model"
	"github.com/hckops/hckctl/pkg/client/kubernetes"
)

func TestNewResources(t *testing.T) {
	namespace := "my-namespace"
	boxName := "my-box-name"
	template := &model.BoxV1{
		Kind: "box/v1",
		Name: "my-name",
		Tags: []string{"my-tag"},
		Image: struct {
			Repository string
			Version    string
		}{
			Repository: "hckops/my-image",
		},
		Shell: "/bin/bash",
		Network: struct{ Ports []string }{Ports: []string{
			"aaa:123",
			"bbb:456:789",
		}},
	}
	opts := &model.CreateOptions{
		Template: template,
		Size:     model.ExtraSmall,
		Labels: map[string]string{
			"a.b.c": "hello",
			"x.y.z": "world",
		},
	}
	expected := &kubernetes.ResourcesOpts{
		Namespace:   namespace,
		Name:        boxName,
		Annotations: opts.Labels,
		Labels: map[string]string{
			"app.kubernetes.io/instance":   "hckops-my-image",
			"app.kubernetes.io/managed-by": "hckops",
			"app.kubernetes.io/name":       "my-box-name",
			"app.kubernetes.io/version":    "latest",
			"com.hckops.schema.kind":       "box-v1",
		},
		Ports: []kubernetes.KubePort{
			{Name: "aaa", Port: "123"},
			{Name: "bbb", Port: "789"},
		},
		PodInfo: &kubernetes.PodInfo{
			Namespace:     namespace,
			PodName:       "INVALID_POD_NAME",
			ContainerName: "hckops/my-image", // sanitized in builder
			ImageName:     "hckops/my-image:latest",
			Env:           nil,
			Resource: &kubernetes.KubeResource{
				Memory: "512Mi",
				Cpu:    "500m",
			},
		},
	}

	assert.Equal(t, expected, newResources(namespace, boxName, opts))
}

func TestLabelSelector(t *testing.T) {
	expected := "com.hckops.schema.kind=box-v1,app.kubernetes.io/name=myName"

	assert.Equal(t, expected, boxNameLabelSelector("myName"))
}

func TestToBoxDetails(t *testing.T) {
	createdTime, _ := time.Parse(time.RFC3339, "2042-12-08T10:30:05.265113665Z")

	deployment := &kubernetes.DeploymentDetails{
		Info: &kubernetes.DeploymentInfo{
			Namespace: "myDeploymentNamespace",
			Name:      "myDeploymentName",
			Healthy:   false,
			PodInfo: &kubernetes.PodInfo{
				Namespace:     "myPodNamespace",
				PodName:       "myPodName",
				ContainerName: "myContainerName",
				Env: map[string]string{
					"MY_KEY_2": "MY_VALUE_2",
					"MY_KEY_1": "MY_VALUE_1",
					"MY_KEY_3": "MY_VALUE_3",
				},
			},
		},
		Created: createdTime,
		Annotations: map[string]string{
			"com.hckops.template.git":          "true",
			"com.hckops.template.git.url":      "myUrl",
			"com.hckops.template.git.revision": "myRevision",
			"com.hckops.template.git.commit":   "myCommit",
			"com.hckops.template.git.name":     "box/base/arch",
			"com.hckops.template.cache.path":   "/tmp/cache/myUuid",
			"com.hckops.box.size":              "m",
		},
	}
	serviceInfo := &kubernetes.ServiceInfo{
		Namespace: "myServiceNamespace",
		Name:      "myServiceName",
		Ports: []kubernetes.KubePort{
			{Name: "name-x", Port: "remote-2"},
			{Name: "name-y", Port: "remote-1"},
			{Name: "name-z", Port: "remote-3"},
		},
	}
	expected := &model.BoxDetails{
		Info: model.BoxInfo{
			Id:      "myPodName",
			Name:    "myDeploymentName",
			Healthy: false,
		},
		TemplateInfo: &model.BoxTemplateInfo{
			GitTemplate: &model.GitTemplateInfo{
				Url:      "myUrl",
				Revision: "myRevision",
				Commit:   "myCommit",
				Name:     "box/base/arch",
			},
		},
		ProviderInfo: &model.BoxProviderInfo{
			Provider: model.Kubernetes,
			KubeProvider: &model.KubeProviderInfo{
				Namespace: "myDeploymentNamespace",
			},
		},
		Size: model.Medium,
		Env: []model.BoxEnv{
			{Key: "MY_KEY_1", Value: "MY_VALUE_1"},
			{Key: "MY_KEY_2", Value: "MY_VALUE_2"},
			{Key: "MY_KEY_3", Value: "MY_VALUE_3"},
		},
		Ports: []model.BoxPort{
			{Alias: "name-y", Local: "none", Remote: "remote-1", Public: false},
			{Alias: "name-x", Local: "none", Remote: "remote-2", Public: false},
			{Alias: "name-z", Local: "none", Remote: "remote-3", Public: false},
		},
		Created: createdTime,
	}
	result, err := toBoxDetails(deployment, serviceInfo)

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}
