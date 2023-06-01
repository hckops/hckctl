package common

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"time"

	"github.com/adrg/xdg"
	"github.com/dchest/uniuri"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/hckops/hckctl/internal/config"
	"github.com/hckops/hckctl/pkg/util"
)

const (
	logDirName string = "hck"
)

func GetLogFile() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", errors.Wrap(err, "unable to retrieve current user")
	}

	logFile := filepath.Join(xdg.StateHome, logDirName, fmt.Sprintf("%s-%s.log", CliName, usr.Username))

	return logFile, nil
}

func SetupLogger(config *LogConfig) error {
	setTimestamp()
	setLevel(parseLevel(config.Level))
	setContext()
	return setFileOutput(config.FilePath)
}

func setTimestamp() {
	zerolog.TimestampFunc = func() time.Time {
		return time.Now().UTC()
	}
}

func parseLevel(value string) zerolog.Level {
	switch value {
	case "debug":
		return zerolog.DebugLevel
	case "info":
		return zerolog.InfoLevel
	case "warning":
		return zerolog.WarnLevel
	case "error":
		return zerolog.ErrorLevel
	default:
		return zerolog.NoLevel
	}
}

// default info
func setLevel(level zerolog.Level) {
	if level == zerolog.NoLevel {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	} else {
		zerolog.SetGlobalLevel(level)
	}
}

func setContext() {
	log.Logger = log.With().Caller().
		Str("source", config.CliName).
		Str("session", generateSession()).
		Logger()
}

func generateSession() string {
	return strings.ToLower(uniuri.NewLen(5))
}

func setFileOutput(filePath string) error {

	if err := util.CreateBaseDir(filePath); err != nil {
		return errors.Wrap(err, "error creating log dir")
	}

	const fileMod os.FileMode = 0600
	mod := os.O_CREATE | os.O_APPEND | os.O_WRONLY
	file, err := os.OpenFile(filePath, mod, fileMod)
	if err != nil {
		return errors.Wrap(err, "error creating log file")
	}

	// TODO close file in rootCmd.run
	// defer func() {
	// 	if file != nil {
	// 		_ = file.Close()
	// 	}
	// }()

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: file, TimeFormat: time.RFC3339})
	return nil
}
