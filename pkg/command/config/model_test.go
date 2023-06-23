package config

import (
	"github.com/hckops/hckctl/pkg/client/ssh"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/hckops/hckctl/pkg/client/kubernetes"
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
		},
		Provider: ProviderConfig{
			Kube: KubeConfig{
				Namespace:    "hckops",
				ConfigPath:   "",
				ResourceSize: "S",
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

func TestToKubeClientConfig(t *testing.T) {
	kubeConfig := &KubeConfig{
		ConfigPath:   "/tmp/config.yml",
		Namespace:    "namespace",
		ResourceSize: "XL",
	}
	expected := &kubernetes.KubeClientConfig{
		ConfigPath: "/tmp/config.yml",
		Namespace:  "namespace",
		Resource: &kubernetes.KubeResource{
			Memory: "512Mi",
			Cpu:    "500m",
		},
	}

	result, err := kubeConfig.ToKubeClientConfig()
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestToSshClientConfig(t *testing.T) {
	cloudConfig := &CloudConfig{
		Host:     "0.0.0.0",
		Port:     2222,
		Username: "myUsername",
		Token:    "myToken",
	}
	expected := &ssh.SshClientConfig{
		Address:  "0.0.0.0:2222",
		Username: "myUsername",
		Token:    "myToken",
	}
	assert.Equal(t, expected, cloudConfig.ToSshClientConfig())
}
