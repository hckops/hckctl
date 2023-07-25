package config

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/hckops/hckctl/pkg/box/model"
)

func TestNewConfig(t *testing.T) {
	logFile := "/tmp/example.log"
	cacheDir := "/tmp/cache/"

	expected := &ConfigV1{
		Kind: "config/v1",
		Log: LogConfig{
			Level:    "info",
			FilePath: logFile,
		},
		Template: TemplateConfig{
			Revision: "main",
			CacheDir: cacheDir,
		},
		Box: BoxConfig{
			Provider: "docker",
			Size:     "S",
		},
		Provider: ProviderConfig{
			Docker: DockerConfig{
				NetworkName: "hckops",
			},
			Kube: KubeConfig{
				Namespace:  "hckops",
				ConfigPath: "",
			},
			Cloud: CloudConfig{
				Host:     "0.0.0.0",
				Port:     2222,
				Username: "",
				Token:    "",
			},
		},
	}

	result := newConfig(logFile, cacheDir)
	assert.Equal(t, expected, result)
}

func TestToDockerBoxOptions(t *testing.T) {
	dockerConfig := &DockerConfig{
		NetworkName: "myNetwork",
	}
	expected := &model.DockerBoxOptions{
		NetworkName:          "myNetwork",
		IgnoreImagePullError: true,
	}
	assert.Equal(t, expected, dockerConfig.ToDockerBoxOptions())
}

func TestToKubeBoxOptions(t *testing.T) {
	kubeConfig := &KubeConfig{
		ConfigPath: "/tmp/config.yml",
		Namespace:  "namespace",
	}
	expected := &model.KubeBoxOptions{
		InCluster:  false,
		ConfigPath: "/tmp/config.yml",
		Namespace:  "namespace",
	}
	assert.Equal(t, expected, kubeConfig.ToKubeBoxOptions())
}

func TestToCloudBoxOptions(t *testing.T) {
	cloudConfig := &CloudConfig{
		Host:     "0.0.0.0",
		Port:     2222,
		Username: "myUsername",
		Token:    "myToken",
	}
	expected := &model.CloudBoxOptions{
		Version:  "hckctl-dev",
		Address:  "0.0.0.0:2222",
		Username: "myUsername",
		Token:    "myToken",
	}
	assert.Equal(t, expected, cloudConfig.ToCloudBoxOptions("hckctl-dev"))
}
