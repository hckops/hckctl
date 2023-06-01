package config

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/adrg/xdg"
	"github.com/pkg/errors"
	"github.com/spf13/viper"

	"github.com/hckops/hckctl/pkg/command/common"
	"github.com/hckops/hckctl/pkg/util"
)

const (
	cliName string = "hckctl"
	dirName string = "hck"

	configEnvName string = "HCK_CONFIG" // overrides .config/hck/config.yml
	configDirEnv  string = "HCK_CONFIG_DIR"
	logDirEnv     string = "HCK_LOG_DIR"
)

func SetupConfig() error {
	configDir, err := loadConfigDir()
	if err != nil {
		return errors.Wrap(err, "error loading config dir")
	}
	configName := "config"
	configType := "yml"
	configPath := filepath.Join(configDir, configName+"."+configType)

	// see https://github.com/spf13/viper/issues/430
	if err := util.CreateBaseDir(configPath); err != nil {
		return errors.Wrap(err, "error creating config dir")
	}

	logFile, err := loadLogFile()
	if err != nil {
		return errors.Wrap(err, "error loading log file")
	}
	if err := util.CreateBaseDir(logFile); err != nil {
		return errors.Wrap(err, "error creating log dir")
	}

	viper.SetConfigName(configName)
	viper.SetConfigType(configType)
	viper.AddConfigPath(configDir)

	viper.SetEnvPrefix(configEnvName)
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {

		// first time only
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// default config
			cliConfig := common.NewConfig(logFile)

			var configString string
			if configString, err = util.ToYaml(&cliConfig); err != nil {
				return errors.Wrap(err, "error encoding config")
			}
			if err := viper.ReadConfig(strings.NewReader(configString)); err != nil {
				return errors.Wrap(err, "error reading config")
			}
			if err := viper.SafeWriteConfigAs(configPath); err != nil {
				return errors.Wrap(err, "error writing config")
			}
		} else {
			return errors.Wrap(err, "invalid config file")
		}
	}
	//viper.Debug()
	return nil
}

func LoadConfig() (*common.Config, error) {
	var configV1 *common.Config
	if err := viper.Unmarshal(&configV1); err != nil {
		return nil, errors.Wrap(err, "error decoding config")
	}
	return configV1, nil
}

func loadConfigDir() (string, error) {
	// override
	if env := os.Getenv(configDirEnv); env != "" {
		return env, nil
	}

	// https://xdgbasedirectoryspecification.com
	xdgDir, err := xdg.ConfigFile(dirName)
	if err != nil {
		return "", errors.Wrapf(err, "unable to create xdg config directory %s", dirName)
	}
	return xdgDir, nil
}

func loadLogFile() (string, error) {
	// override
	if env := os.Getenv(logDirEnv); env != "" {
		return env, nil
	}

	usr, err := user.Current()
	if err != nil {
		return "", errors.Wrap(err, "unable to retrieve current user")
	}

	logFile := filepath.Join(xdg.StateHome, dirName, fmt.Sprintf("%s-%s.log", cliName, usr.Username))

	return logFile, nil
}
