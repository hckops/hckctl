package config

import (
	"github.com/hckops/hckctl/pkg/client/ssh"
	"net"
	"strconv"

	"github.com/hckops/hckctl/pkg/box/model"
	"github.com/hckops/hckctl/pkg/client/kubernetes"
	"github.com/hckops/hckctl/pkg/command/common"
	"github.com/hckops/hckctl/pkg/logger"
	"github.com/hckops/hckctl/pkg/template/schema"
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
	Template TemplateConfig `yaml:"template"`
	Box      BoxConfig      `yaml:"box"`
	Provider ProviderConfig `yaml:"provider"`
}

type LogConfig struct {
	Level    string `yaml:"level"`
	FilePath string `yaml:"filePath"`
}

type TemplateConfig struct {
	Revision string `yaml:"revision"`
	CacheDir string `yaml:"cacheDir"`
}

type BoxConfig struct {
	Provider string `yaml:"provider"`
}

type ProviderConfig struct {
	Kube  KubeConfig  `yaml:"kube"`
	Cloud CloudConfig `yaml:"cloud"`
}

type KubeConfig struct {
	ConfigPath   string `yaml:"configPath"`
	Namespace    string `yaml:"namespace"`
	ResourceSize string `yaml:"resourceSize"`
}

func (c *KubeConfig) ToKubeClientConfig() (*kubernetes.KubeClientConfig, error) {
	if size, err := existResourceSize(c.ResourceSize); err != nil {
		return nil, err
	} else {
		return &kubernetes.KubeClientConfig{
			ConfigPath: c.ConfigPath,
			Namespace:  c.Namespace,
			Resource:   size.toKubeResource(),
		}, nil
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

func (c *CloudConfig) ToSshClientConfig() *ssh.SshClientConfig {
	return &ssh.SshClientConfig{
		Address:  c.address(),
		Username: c.Username,
		Token:    c.Token,
	}
}

func newConfig(logFile, cacheDir string) *ConfigV1 {
	return &ConfigV1{
		Kind: schema.KindConfigV1.String(),
		Log: LogConfig{
			Level:    logger.InfoLogLevel.String(),
			FilePath: logFile,
		},
		Template: TemplateConfig{
			Revision: common.TemplateSourceRevision,
			CacheDir: cacheDir,
		},
		Box: BoxConfig{
			Provider: model.Docker.String(),
		},
		Provider: ProviderConfig{
			Kube: KubeConfig{
				Namespace:    common.ProjectName,
				ConfigPath:   "",
				ResourceSize: small.String(),
			},
			Cloud: CloudConfig{
				Host:     "0.0.0.0",
				Port:     2222,
				Username: "",
				Token:    "",
			},
		},
	}
}
