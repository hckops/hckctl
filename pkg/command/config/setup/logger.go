package setup

import (
	"fmt"
	"os/user"
	"path/filepath"

	"github.com/adrg/xdg"
	"github.com/pkg/errors"

	"github.com/hckops/hckctl/pkg/command/common"
	"github.com/hckops/hckctl/pkg/logger"
)

const (
	logDirName string = "hck"
)

func SetupLogger(configRef *common.ConfigRef) (func() error, error) {
	logConfig := configRef.Config.Log
	logger.SetTimestamp()
	logger.SetLevel(logger.ParseLevel(logConfig.Level))
	logger.SetContext(common.CliName)
	return logger.SetFileOutput(logConfig.FilePath)
}

func getLogFile() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", errors.Wrap(err, "unable to retrieve current user")
	}

	logFile := filepath.Join(xdg.StateHome, logDirName, fmt.Sprintf("%s-%s.log", common.CliName, usr.Username))

	return logFile, nil
}
