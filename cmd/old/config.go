package old

import (
	"fmt"
	"github.com/hckops/hckctl/pkg/util"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/thediveo/enumflag/v2"

	"github.com/hckops/hckctl/internal/config"
)

type ProviderFlag enumflag.Flag

const (
	DockerFlag ProviderFlag = iota
	KubernetesFlag
	CloudFlag
)

var ProviderIds = map[ProviderFlag][]string{
	DockerFlag:     {string(config.Docker)},
	KubernetesFlag: {string(config.Kubernetes)},
	CloudFlag:      {string(config.Cloud)},
}

func ProviderToId(provider ProviderFlag) string {
	return ProviderIds[provider][0]
}

func ProviderToFlag(value config.Provider) (ProviderFlag, error) {
	switch value {
	case config.Docker:
		return DockerFlag, nil
	case config.Kubernetes:
		return KubernetesFlag, nil
	case config.Cloud:
		return CloudFlag, nil
	default:
		return 999, fmt.Errorf("invalid provider")
	}
}

var configPath string

func InitCliConfig() {
	configHome := config.ConfigHome()
	configName := "config"
	configType := "yml"
	configPath = filepath.Join(configHome, configName+"."+configType)
	// see https://github.com/spf13/viper/issues/430
	config.EnsurePathOrDie(configPath, config.DefaultDirectoryMod)

	viper.SetConfigName(configName)
	viper.SetConfigType(configType)
	viper.AddConfigPath(configHome)

	viper.SetEnvPrefix(config.ConfigNameEnv)
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {

		// first time only
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// default config
			cliConfig := config.NewConfig()

			var configString string
			if configString, err = util.ToYaml(&cliConfig); err != nil {
				log.Fatal().Err(err).Msg("error encoding config")
			}
			if err := viper.ReadConfig(strings.NewReader(configString)); err != nil {
				log.Fatal().Err(err).Msg("error reading config")
			}
			if err := viper.SafeWriteConfigAs(configPath); err != nil {
				log.Fatal().Err(err).Msg("error writing config")
			}
		} else {
			log.Fatal().Err(err).Msg("invalid config file")
		}
	}
	//viper.Debug()
}

// TODO add command to set/reset configs
func NewConfigCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "config",
		Short: "prints current configurations",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(fmt.Sprintf("# %s", configPath))
			GetCliConfig().Print()
		},
	}
}

func GetCliConfig() *config.ConfigV1 {
	var cliConfig *config.ConfigV1
	if err := viper.Unmarshal(&cliConfig); err != nil {
		log.Fatal().Err(err).Msg("error decoding config")
	}
	return cliConfig
}
