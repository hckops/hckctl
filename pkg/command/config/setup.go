package config

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/adrg/xdg"
	"github.com/pkg/errors"
	"github.com/spf13/viper"

	"github.com/hckops/hckctl/pkg/command/common"
	"github.com/hckops/hckctl/pkg/util"
)

const (
	configDirName string = "hck"
	configDirEnv  string = "HCK_CONFIG_DIR" // overrides .config/hck
	configEnvName string = "HCK_CONFIG"
)

// SetupConfig loads the config or initialize the default
func SetupConfig() (*common.ConfigV1, error) {
	err := initConfig(false)
	if err != nil {
		return nil, err
	}
	return loadConfig()
}

func initConfig(force bool) error {
	configDir, err := getConfigDir()
	if err != nil {
		return errors.Wrap(err, "invalid config dir")
	}
	configName := "config"
	configType := "yml"
	configPath := filepath.Join(configDir, configName+"."+configType)

	// see https://github.com/spf13/viper/issues/430
	if err := util.CreateBaseDir(configPath); err != nil {
		return errors.Wrap(err, "error creating config dir")
	}

	viper.SetConfigName(configName)
	viper.SetConfigType(configType)
	viper.AddConfigPath(configDir)

	viper.SetEnvPrefix(configEnvName)
	viper.AutomaticEnv()

	// reset config
	if force {
		return createDefaultConfig(configPath)
	} else if err := viper.ReadInConfig(); err != nil {

		// first time only
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return createDefaultConfig(configPath)
		} else {
			return errors.Wrap(err, "invalid config file")
		}
	}
	//viper.Debug()
	return nil
}

func getConfigDir() (string, error) {
	// override
	if env := os.Getenv(configDirEnv); env != "" {
		return env, nil
	}

	// https://xdgbasedirectoryspecification.com
	xdgDir, err := xdg.ConfigFile(configDirName)
	if err != nil {
		return "", errors.Wrapf(err, "unable to create xdg config directory %s", configDirName)
	}
	return xdgDir, nil
}

func createDefaultConfig(configPath string) error {
	// default log file
	logFile, err := common.GetLogFile()
	if err != nil {
		return errors.Wrap(err, "invalid log file")
	}

	// default config
	cliConfig := common.NewConfig(logFile)

	var configString string
	if configString, err = util.ToYaml(&cliConfig); err != nil {
		return errors.Wrap(err, "error encoding config")
	}
	if err := viper.ReadConfig(strings.NewReader(configString)); err != nil {
		return errors.Wrap(err, "error reading config")
	}
	// SafeWriteConfigAs prevents override
	if err := viper.WriteConfigAs(configPath); err != nil {
		return errors.Wrap(err, "error writing config")
	}
	return nil
}

func loadConfig() (*common.ConfigV1, error) {
	var configV1 *common.ConfigV1
	// "exact" makes sure to fail if fields are invalid
	if err := viper.UnmarshalExact(&configV1); err != nil {
		return nil, errors.Wrap(err, "error decoding config")
	}
	return configV1, nil
}
