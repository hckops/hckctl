package model

import (
	"fmt"

	"github.com/hckops/hckctl/internal/common"
	"github.com/hckops/hckctl/pkg/util"
	"github.com/rs/zerolog/log"
)

type ConfigV1 struct {
	Kind string    `yaml:"kind"`
	Box  BoxConfig `yaml:"box"`
	Log  LogConfig `yaml:"log"`
}

type LogConfig struct {
	Level    string `yaml:"level"`
	FilePath string `yaml:"filePath"`
}

type BoxConfig struct {
	Revision string     `yaml:"revision"`
	Provider Provider   `yaml:"provider"`
	Kube     KubeConfig `yaml:"kube"`
}

type Provider string

const (
	Docker     Provider = "docker"
	Kubernetes Provider = "kube"
	Cloud      Provider = "cloud"
)

type KubeConfig struct {
	Namespace  string        `yaml:"namespace"`
	ConfigPath string        `yaml:"configPath"`
	Resources  KubeResources `yaml:"resources"`
}

type KubeResources struct {
	Memory string `yaml:"memory"`
	Cpu    string `yaml:"cpu"`
}

func NewConfig() *ConfigV1 {
	return &ConfigV1{
		Kind: "config/v1",
		Box: BoxConfig{
			Revision: "main",
			Provider: Docker,
			Kube: KubeConfig{
				Namespace:  "labs",
				ConfigPath: "~/.kube/config",
				Resources: KubeResources{
					Memory: "512Mi",
					Cpu:    "500m",
				},
			},
		},
		Log: LogConfig{
			Level:    "info",
			FilePath: common.DefaultLogFile,
		},
	}
}

func (config *ConfigV1) Print() {
	value, err := util.ToYaml(&config)
	if err != nil {
		log.Warn().Msg("invalid config")
	}
	fmt.Print(value)
}
