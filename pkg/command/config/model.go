package config

import (
	"net"
	"strconv"

	"github.com/hckops/hckctl/pkg/box/model"
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
	Provider string      `yaml:"provider"`
	Kube     KubeConfig  `yaml:"kube"`
	Cloud    CloudConfig `yaml:"cloud"`
}

type KubeConfig struct {
	Namespace    string `yaml:"namespace"`
	ConfigPath   string `yaml:"configPath"`
	ResourceSize string `yaml:"resourceSize"`
}

type CloudConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Token    string `yaml:"token"`
}

func (c *CloudConfig) Address() string {
	return net.JoinHostPort(c.Host, strconv.Itoa(c.Port))
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
			Kube: KubeConfig{
				Namespace:    common.ProjectName,
				ConfigPath:   "",
				ResourceSize: "s",
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
