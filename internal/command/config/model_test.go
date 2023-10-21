package config

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/hckops/hckctl/pkg/common/model"
	"github.com/hckops/hckctl/pkg/util"
)

func TestNewConfig(t *testing.T) {

	configOpts := &configOptions{
		logFile:    "/tmp/example.log",
		cacheDir:   "/tmp/cache/",
		shareDir:   "/tmp/share/",
		taskLogDir: "/tmp/task/log/",
	}

	expected := &ConfigV1{
		Kind:    "config/v1",
		Version: "1.0",
		Log: LogConfig{
			Level:    "info",
			FilePath: "/tmp/example.log",
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
			Privileged: false,
			Vpn: []VpnConfig{
				{Name: "default", Path: "/path/to/client.ovpn"},
			},
		},
		Template: TemplateConfig{
			Revision: "main",
			CacheDir: "/tmp/cache/",
		},
		Common: CommonConfig{
			ShareDir: "/tmp/share/",
		},
		Box: BoxConfig{
			Provider: "docker",
			Size:     "S",
		},
		Task: TaskConfig{
			Provider: "docker",
			LogDir:   "/tmp/task/log/",
		},
	}

	result := newConfig(configOpts)
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

func TestVpnNetworks(t *testing.T) {
	networkConfig := NetworkConfig{
		Vpn: []VpnConfig{
			{Name: "readme", Path: "../../../README.md"},
			{Name: "license", Path: "../../../LICENSE"},
		},
	}
	assert.Equal(t, 2, len(networkConfig.VpnNetworks()))
}

func TestToNetworkVpnInfo(t *testing.T) {
	networkConfig := NetworkConfig{
		Vpn: []VpnConfig{
			{Name: "readme", Path: "../../../README.md"},
			{Name: "license", Path: "../../../LICENSE"},
		},
	}

	emptyVpn, emptyErr := networkConfig.ToNetworkVpnInfo("")
	assert.Nil(t, emptyVpn)
	assert.Nil(t, emptyErr)

	validVpn, validErr := networkConfig.ToNetworkVpnInfo("readme")
	configFile, _ := util.ReadFile("../../../README.md")
	expected := &model.NetworkVpnInfo{
		Name:        "readme",
		LocalPath:   "../../../README.md",
		ConfigValue: configFile,
	}
	assert.Equal(t, expected, validVpn)
	assert.Nil(t, validErr)

	invalidVpn, invalidErr := networkConfig.ToNetworkVpnInfo("foo")
	assert.Nil(t, invalidVpn)
	assert.EqualError(t, invalidErr, "vpn not found name=foo")
}

func TestToShareDirInfo(t *testing.T) {
	commonConfig := &CommonConfig{
		ShareDir: "myShareDir",
	}

	expected := &model.ShareDirInfo{
		LocalPath:  "myShareDir",
		RemotePath: "/hck/share",
		LockDir:    true,
	}
	assert.Equal(t, expected, commonConfig.ToShareDirInfo(true))
}
