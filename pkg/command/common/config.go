package common

import (
	"net"
	"strconv"

	"github.com/hckops/hckctl/pkg/template/schema"
)

// ConfigRef is a wrapper used to avoid global variables and to reference the config value in the commands
// before they are actually loaded with viper in each PersistentPreRunE.
// The ConfigV1 model is in a common package to avoid "import cycle not allowed"
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
	DirPath  string `yaml:"dirPath"`
}

type BoxConfig struct {
	Provider Provider    `yaml:"provider"`
	Kube     KubeConfig  `yaml:"kube"`
	Cloud    CloudConfig `yaml:"cloud"`
}

type Provider string

const (
	Docker     Provider = "docker"
	Kubernetes Provider = "kube"
	Argo       Provider = "argo"
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

func NewConfig(logFile, sourceDir string) *ConfigV1 {
	return &ConfigV1{
		Kind: schema.KindConfigV1.String(),
		Log: LogConfig{
			Level:    "info", // TODO enum
			FilePath: logFile,
		},
		Template: TemplateConfig{
			Revision: TemplateRevision,
			DirPath:  sourceDir,
		},
		Box: BoxConfig{
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
