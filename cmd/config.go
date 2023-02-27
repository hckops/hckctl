package cmd

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/thediveo/enumflag/v2"

	"github.com/hckops/hckctl/internal/common"
	"github.com/hckops/hckctl/internal/model"
)

// see enum https://stackoverflow.com/questions/50824554/permitted-flag-values-for-cobra
type ProviderFlag enumflag.Flag

const (
	DockerFlag ProviderFlag = iota
	KubernetesFlag
	CloudFlag
)

var ProviderIds = map[ProviderFlag][]string{
	DockerFlag:     {string(model.Docker)},
	KubernetesFlag: {string(model.Kubernetes)},
	CloudFlag:      {string(model.Cloud)},
}

func ProviderToId(provider ProviderFlag) string {
	return ProviderIds[provider][0]
}

func ProviderToFlag(value model.Provider) (ProviderFlag, error) {
	switch value {
	case model.Docker:
		return DockerFlag, nil
	case model.Kubernetes:
		return KubernetesFlag, nil
	case model.Cloud:
		return CloudFlag, nil
	default:
		return 999, fmt.Errorf("invalid provider")
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
		cliConfig := model.NewCliConfig()

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
			GetCliConfig().Print()
		},
	}
}

func GetCliConfig() *model.CliConfig {
	var cliConfig *model.CliConfig
	if err := viper.Unmarshal(&cliConfig); err != nil {
		log.Fatal().Err(fmt.Errorf("error decoding config: %w", err))
	}
	return cliConfig
}
