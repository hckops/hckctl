package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/thediveo/enumflag/v2"

	"github.com/hckops/hckctl/internal/common"
)

type CliConfig struct {
	Kind string    `yaml:"kind"`
	Box  BoxConfig `yaml:"box"`
	Log  LogConfig `yaml:"log"`
}

// tmp file
var DefaultLogFile = filepath.Join(os.TempDir(), fmt.Sprintf("%s-%s.log", common.CliName, common.GetUserOrDie()))

type LogConfig struct {
	Level    string `yaml:"level"`
	FilePath string `yaml:"filePath"`
}

type BoxConfig struct {
	Revision string     `yaml:"revision"`
	Provider string     `yaml:"provider"`
	Kube     KubeConfig `yaml:"kube"`
}

// see enum https://stackoverflow.com/questions/50824554/permitted-flag-values-for-cobra
type Provider enumflag.Flag

const (
	Docker Provider = iota
	Kubernetes
	Cloud
)

var ProviderIds = map[Provider][]string{
	Docker:     {"docker"},
	Kubernetes: {"kube"},
	Cloud:      {"cloud"},
}

func ProviderToId(provider Provider) string {
	return ProviderIds[provider][0]
}

func StringToProvider(value string) (Provider, error) {
	switch value {
	case "docker":
		return Docker, nil
	case "kube":
		return Kubernetes, nil
	case "cloud":
		return Cloud, nil
	default:
		return 999, fmt.Errorf("invalid provider")
	}
}

type KubeConfig struct {
	Namespace  string `yaml:"namespace"`
	ConfigPath string `yaml:"configPath"`
}

func newCliConfig() *CliConfig {
	return &CliConfig{
		Kind: "config/v1",
		Box: BoxConfig{
			Revision: "main",
			Provider: ProviderToId(Docker),
			Kube: KubeConfig{
				Namespace:  "labs",
				ConfigPath: "~/.kube/config",
			},
		},
		Log: LogConfig{
			Level:    "info",
			FilePath: DefaultLogFile,
		},
	}
}

func InitCliConfig() {
	configHome := common.ConfigHome()
	configName := "config"
	configType := "yml"
	configPath := filepath.Join(configHome, configName+"."+configType)
	// see https://github.com/spf13/viper/issues/430
	common.EnsurePathOrDie(configPath, common.DefaultDirectoryMod)

	viper.AddConfigPath(configHome)
	viper.SetConfigName(configName)
	viper.SetConfigType(configType)

	viper.AutomaticEnv()
	viper.SetEnvPrefix(common.ConfigNameEnv)

	// first time only
	if err := viper.ReadInConfig(); err != nil {

		// default config
		cliConfig := newCliConfig()

		var configString string
		if configString, err = common.ToYaml(&cliConfig); err != nil {
			log.Fatal().Err(fmt.Errorf("error encoding config: %w", err))
		}
		if err := viper.ReadConfig(strings.NewReader(configString)); err != nil {
			log.Fatal().Err(fmt.Errorf("error reading config: %w", err))
		}
		if err := viper.SafeWriteConfigAs(configPath); err != nil {
			log.Fatal().Err(fmt.Errorf("error writing config: %w", err))
		}
	}
}

// TODO add command to set/reset configs
func NewConfigCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "config",
		Short: "prints current configurations",
		Run: func(cmd *cobra.Command, args []string) {
			// TODO override validation
			GetCliConfig().print()
		},
	}
}

func GetCliConfig() *CliConfig {
	var cliConfig *CliConfig
	if err := viper.Unmarshal(&cliConfig); err != nil {
		log.Fatal().Err(fmt.Errorf("error decoding config: %w", err))
	}
	return cliConfig
}

func (config *CliConfig) print() {
	value, err := common.ToYaml(&config)
	if err != nil {
		log.Warn().Msg("invalid config")
	}
	fmt.Print(value)
}
