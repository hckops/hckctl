package config

import (
	"fmt"
	"net"
	"strconv"

	"github.com/hckops/hckctl/internal/command/common"
	boxModel "github.com/hckops/hckctl/pkg/box/model"
	commonModel "github.com/hckops/hckctl/pkg/common/model"
	"github.com/hckops/hckctl/pkg/logger"
	"github.com/hckops/hckctl/pkg/schema"
	taskModel "github.com/hckops/hckctl/pkg/task/model"
	"github.com/hckops/hckctl/pkg/util"
)

// ConfigRef is a wrapper used to avoid global variables.
// It's used to reference the config value in the commands
// before they are actually loaded with viper in each PersistentPreRunE.
type ConfigRef struct {
	Config *ConfigV1
}

// TODO not used, useful for migrations
const (
	currentVersion = "1.0"
)

type ConfigV1 struct {
	Kind     string         `yaml:"kind"`
	Version  string         `yaml:"version"`
	Log      LogConfig      `yaml:"log"`
	Provider ProviderConfig `yaml:"provider"`
	Network  NetworkConfig  `yaml:"network"`
	Template TemplateConfig `yaml:"template"`
	Common   CommonConfig   `yaml:"common"`
	Box      BoxConfig      `yaml:"box"`
	Task     TaskConfig     `yaml:"task"`
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

func (c *DockerConfig) ToDockerOptions() *commonModel.DockerOptions {
	return &commonModel.DockerOptions{
		NetworkName:          c.NetworkName,
		IgnoreImagePullError: true, // always allow to start offline/obsolete images
	}
}

type KubeConfig struct {
	ConfigPath string `yaml:"configPath"`
	Namespace  string `yaml:"namespace"`
}

func (c *KubeConfig) ToKubeOptions() *commonModel.KubeOptions {
	return &commonModel.KubeOptions{
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

func (c *CloudConfig) ToCloudOptions(version string) *commonModel.CloudOptions {
	return &commonModel.CloudOptions{
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

func (c *NetworkConfig) VpnNetworks() map[string]commonModel.VpnNetworkInfo {
	info := map[string]commonModel.VpnNetworkInfo{}
	for _, network := range c.Vpn {
		// ignores invalid paths
		if configFile, err := util.ReadFile(network.Path); err == nil {
			info[network.Name] = commonModel.VpnNetworkInfo{
				Name:        network.Name,
				LocalPath:   network.Path,
				ConfigValue: configFile,
			}
		}
	}
	return info
}

func (c *NetworkConfig) ToVpnNetworkInfo(vpnName string) (*commonModel.VpnNetworkInfo, error) {
	if vpnName != "" {
		if vpnNetworkInfo, ok := c.VpnNetworks()[vpnName]; ok {
			return &vpnNetworkInfo, nil
		} else {
			return nil, fmt.Errorf("vpn not found name=%s", vpnName)
		}
	}
	return nil, nil
}

type TemplateConfig struct {
	Revision string `yaml:"revision"`
	CacheDir string `yaml:"cacheDir"`
}

type CommonConfig struct {
	ShareDir string `yaml:"shareDir"`
}

type BoxConfig struct {
	Provider string `yaml:"provider"`
	Size     string `yaml:"size"`
}

type TaskConfig struct {
	Provider string `yaml:"provider"`
	LogDir   string `yaml:"logDir"`
}

type configOptions struct {
	logFile    string
	cacheDir   string
	shareDir   string
	taskLogDir string
}

func newConfig(opts *configOptions) *ConfigV1 {
	return &ConfigV1{
		Kind:    schema.KindConfigV1.String(),
		Version: currentVersion,
		Log: LogConfig{
			Level:    logger.InfoLogLevel.String(),
			FilePath: opts.logFile,
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
			CacheDir: opts.cacheDir,
		},
		Common: CommonConfig{
			ShareDir: opts.shareDir,
		},
		Box: BoxConfig{
			Provider: boxModel.Docker.String(),
			Size:     boxModel.Small.String(),
		},

		Task: TaskConfig{
			Provider: taskModel.Docker.String(),
			LogDir:   opts.taskLogDir,
		},
	}
}
