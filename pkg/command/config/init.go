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

// TODO move ???
const (
	defaultDirectoryMod os.FileMode = 0755
	defaultFileMod      os.FileMode = 0600
)

const (
	cliName string = "hckctl"
	dirName string = "hck"

	configEnvName string = "HCK_CONFIG" // overrides .config/hck/config.yml
	configDirEnv  string = "HCK_CONFIG_DIR"
	logDirEnv     string = "HCK_LOG_DIR"
)

func InitConfig() error {
	configDir, err := loadConfigDir()
	if err != nil {
		return errors.Wrap(err, "error loading config dir")
	}
	configName := "config"
	configType := "yml"
	configPath := filepath.Join(configDir, configName+"."+configType)

	// see https://github.com/spf13/viper/issues/430
	if err := initDir(configPath, defaultDirectoryMod); err != nil {
		return errors.Wrap(err, "error init config dir")
	}

	logFile, err := loadLogFile()
	if err != nil {
		return errors.Wrap(err, "error loading log file")
	}
	if err := initDir(logFile, defaultDirectoryMod); err != nil {
		return errors.Wrap(err, "error init log dir")
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

func LoadConfig() (*common.ConfigV1, error) {
	var configV1 *common.ConfigV1
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

func initDir(path string, mod os.FileMode) error {
	dir := filepath.Dir(path)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, mod); err != nil {
			return errors.Wrapf(err, "unable to create dir %s", dir)
		}
	}
	return nil
}
