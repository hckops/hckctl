package config

import (
	"fmt"
	"github.com/hckops/hckctl/internal/config"
	"github.com/hckops/hckctl/pkg/command/common"
	"os"
	"path/filepath"
	"strings"

	"github.com/adrg/xdg"
	"github.com/hckops/hckctl/pkg/util"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

const (
	DefaultDirectoryMod os.FileMode = 0755
	DefaultFileMod      os.FileMode = 0600
)

const (
	CliName       string = "hckctl"
	ConfigDir     string = "hck"
	ConfigNameEnv string = "HCK_CONFIG" //  overrides .config/hck/config.yml
)

var configPath string

func InitConfig(globalOpts *common.GlobalCmdOptions) error {
	log.Info().Msg("CONFIG")
	return nil
}

func InitCliConfig(globalOpts *common.GlobalCmdOptions) {
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

func ConfigHome() string {
	if env := os.Getenv(ConfigNameEnv); env != "" {
		return env
	}

	// https://xdgbasedirectoryspecification.com
	xdgHome, err := xdg.ConfigFile(ConfigDir)
	if err != nil {
		log.Fatal().Err(err).Msgf("unable to create configuration directory: %s", ConfigDir)
	}

	return xdgHome
}

func LogHome() string {
	return "TODO"
}

// tmp file
var DefaultLogFile = filepath.Join(os.TempDir(), fmt.Sprintf("%s.log", CliName))
