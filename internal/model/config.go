package model

import (
	"fmt"

	"github.com/hckops/hckctl/internal/common"
	"github.com/rs/zerolog/log"
)

type CliConfig struct {
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
	Namespace  string `yaml:"namespace"`
	ConfigPath string `yaml:"configPath"`
}

func NewCliConfig() *CliConfig {
	return &CliConfig{
		Kind: "config/v1",
		Box: BoxConfig{
			Revision: "main",
			Provider: Docker,
			Kube: KubeConfig{
				Namespace:  "labs",
				ConfigPath: "~/.kube/config",
			},
		},
		Log: LogConfig{
			Level:    "info",
			FilePath: common.DefaultLogFile,
		},
	}
}

func (config *CliConfig) Print() {
	value, err := common.ToYaml(&config)
	if err != nil {
		log.Warn().Msg("invalid config")
	}
	fmt.Print(value)
}
