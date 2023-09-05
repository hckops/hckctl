package config

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/hckops/hckctl/pkg/common/model"
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
		Network: NetworkConfig{
			Vpn: []VpnConfig{
				{Name: "default", Path: "/path/to/client.ovpn"},
			},
		},
		Template: TemplateConfig{
			Revision: "main",
			CacheDir: cacheDir,
		},
		Box: BoxConfig{
			Provider: "docker",
			Size:     "S",
		},
		Lab: LabConfig{
			Provider: "cloud",
			Vpn:      "default",
		},
	}

	result := newConfig(logFile, cacheDir)
	assert.Equal(t, expected, result)
}

func TestToDockerOptions(t *testing.T) {
	dockerConfig := &DockerConfig{
		NetworkName: "myNetwork",
	}
	expected := &model.DockerOptions{
		NetworkName:          "myNetwork",
		IgnoreImagePullError: true,
	}
	assert.Equal(t, expected, dockerConfig.ToDockerOptions())
}

func TestToKubeOptions(t *testing.T) {
	kubeConfig := &KubeConfig{
		ConfigPath: "/tmp/config.yml",
		Namespace:  "namespace",
	}
	expected := &model.KubeOptions{
		InCluster:  false,
		ConfigPath: "/tmp/config.yml",
		Namespace:  "namespace",
	}
	assert.Equal(t, expected, kubeConfig.ToKubeOptions())
}

func TestToCloudOptions(t *testing.T) {
	cloudConfig := &CloudConfig{
		Host:     "0.0.0.0",
		Port:     2222,
		Username: "myUsername",
		Token:    "myToken",
	}
	expected := &model.CloudOptions{
		Version:  "hckctl-dev",
		Address:  "0.0.0.0:2222",
		Username: "myUsername",
		Token:    "myToken",
	}
	assert.Equal(t, expected, cloudConfig.ToCloudOptions("hckctl-dev"))
}
