package cmd

import (
	"os"
	"time"

	"github.com/hckops/hckctl/internal/common"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func InitFileLogger() {
	config := GetCliConfig().Log

	setTimestamp()
	setLevel(parseLevel(config.Level))
	setContext()
	setFileOutput(config.FilePath)
}

func setTimestamp() {
	zerolog.TimestampFunc = func() time.Time {
		return time.Now().UTC()
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

func setContext() {
	log.Logger = log.With().Caller().Str("source", common.CliName).Logger()
}

// TODO close file in rootCmd.run
func setFileOutput(filePath string) {
	common.EnsurePathOrDie(filePath, common.DefaultDirectoryMod)
	mod := os.O_CREATE | os.O_APPEND | os.O_WRONLY
	file, err := os.OpenFile(filePath, mod, common.DefaultFileMod)
	if err != nil {
		panic(err)
	}
	// defer func() {
	// 	if file != nil {
	// 		_ = file.Close()
	// 	}
	// }()

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: file, TimeFormat: time.RFC3339})
}
