package config

import (
	"net"
	"strconv"

	"github.com/hckops/hckctl/internal/command/common"
	boxModel "github.com/hckops/hckctl/pkg/box/model"
	labModel "github.com/hckops/hckctl/pkg/lab/model"
	"github.com/hckops/hckctl/pkg/logger"
	"github.com/hckops/hckctl/pkg/provider"
	"github.com/hckops/hckctl/pkg/schema"
)

// ConfigRef is a wrapper used to avoid global variables.
// It's used to reference the config value in the commands
// before they are actually loaded with viper in each PersistentPreRunE.
type ConfigRef struct {
	Config *ConfigV1
}

type ConfigV1 struct {
	Kind     string         `yaml:"kind"`
	Log      LogConfig      `yaml:"log"`
	Provider ProviderConfig `yaml:"provider"`
	Network  NetworkConfig  `yaml:"network"`
	Template TemplateConfig `yaml:"template"`
	Box      BoxConfig      `yaml:"box"`
	Lab      LabConfig      `yaml:"lab"`
}

type LogConfig struct {
	Level    string `yaml:"level"`
	FilePath string `yaml:"filePath"`
}

type ProviderConfig struct {
	Docker DockerConfig `yaml:"docker"`
	Kube   KubeConfig   `yaml:"kube"`
	Cloud  CloudConfig  `yaml:"cloud"`
}

type DockerConfig struct {
	NetworkName string `yaml:"networkName"`
}

func (c *DockerConfig) ToDockerOptions() *provider.DockerOptions {
	return &provider.DockerOptions{
		NetworkName:          c.NetworkName,
		IgnoreImagePullError: true, // always allow to start offline/obsolete images
	}
}

type KubeConfig struct {
	ConfigPath string `yaml:"configPath"`
	Namespace  string `yaml:"namespace"`
}

func (c *KubeConfig) ToKubeOptions() *provider.KubeOptions {
	return &provider.KubeOptions{
		InCluster:  false,
		ConfigPath: c.ConfigPath,
		Namespace:  c.Namespace,
	}
}

type CloudConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Token    string `yaml:"token"`
}

func (c *CloudConfig) address() string {
	return net.JoinHostPort(c.Host, strconv.Itoa(c.Port))
}

func (c *CloudConfig) ToCloudOptions(version string) *provider.CloudOptions {
	return &provider.CloudOptions{
		Version:  version,
		Address:  c.address(),
		Username: c.Username,
		Token:    c.Token,
	}
}

type NetworkConfig struct {
	Vpn []VpnConfig `yaml:"vpn"`
}

type VpnConfig struct {
	Name string `yaml:"name"`
	Path string `yaml:"path"`
}

type TemplateConfig struct {
	Revision string `yaml:"revision"`
	CacheDir string `yaml:"cacheDir"`
}

type BoxConfig struct {
	Provider string `yaml:"provider"`
	Size     string `yaml:"size"`
}

type LabConfig struct {
	Provider string `yaml:"provider"`
	Vpn      string `yaml:"vpn"`
}

func newConfig(logFile, cacheDir string) *ConfigV1 {
	return &ConfigV1{
		Kind: schema.KindConfigV1.String(),
		Log: LogConfig{
			Level:    logger.InfoLogLevel.String(),
			FilePath: logFile,
		},
		Provider: ProviderConfig{
			Docker: DockerConfig{
				NetworkName: common.ProjectName,
			},
			Kube: KubeConfig{
				Namespace:  common.ProjectName,
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
				{Name: common.DefaultVpnName, Path: "/path/to/client.ovpn"},
			},
		},
		Template: TemplateConfig{
			Revision: common.TemplateSourceRevision,
			CacheDir: cacheDir,
		},
		Box: BoxConfig{
			Provider: boxModel.Docker.String(),
			Size:     boxModel.Small.String(),
		},
		Lab: LabConfig{
			Provider: labModel.Cloud.String(),
			Vpn:      common.DefaultVpnName,
		},
	}
}
