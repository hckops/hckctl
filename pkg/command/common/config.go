package common

import (
	"net"
	"strconv"
)

// ConfigRef is a wrapper used to avoid global variables and to reference the config value in the commands
// before they are actually loaded with viper in each PersistentPreRunE.
// The Config model is in a common package to avoid "import cycle not allowed"
type ConfigRef struct {
	Config *Config
}

type Config struct {
	Kind string     `yaml:"kind"`
	Box  BoxConfig  `yaml:"box"`
	Log  *LogConfig `yaml:"log"`
}

type LogConfig struct {
	Level    string `yaml:"level"`
	FilePath string `yaml:"filePath"`
}

type BoxConfig struct {
	Revision string      `yaml:"revision"`
	Provider Provider    `yaml:"provider"`
	Kube     KubeConfig  `yaml:"kube"`
	Cloud    CloudConfig `yaml:"cloud"`
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

type CloudConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Token    string `yaml:"token"`
}

func (c *CloudConfig) Address() string {
	return net.JoinHostPort(c.Host, strconv.Itoa(c.Port))
}

type KubeResources struct {
	Memory string `yaml:"memory"`
	Cpu    string `yaml:"cpu"`
}

func NewConfig(logFile string) *Config {
	return &Config{
		Kind: "config/v1",
		Log: &LogConfig{
			Level:    "info",
			FilePath: logFile,
		},
		Box: BoxConfig{
			Revision: "main",
			Provider: Docker,
			Kube: KubeConfig{
				Namespace:  "labs",
				ConfigPath: "",
				Resources: KubeResources{
					Memory: "512Mi",
					Cpu:    "500m",
				},
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
