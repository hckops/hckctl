package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/hckops/hckctl/internal/common"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// tmp file
var DefaultLogFile = filepath.Join(os.TempDir(), fmt.Sprintf("hckctl-%s.log", common.GetUserOrDie()))

func InitFileLogger(flags *Flags) {
	setTimestamp()
	setLevel(parseLevel(flags.LogLevel))
	setContext()
	setFileOutput()
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
	log.Logger = log.With().Caller().Str("source", "hckctl").Logger()
}

// TODO close file
func setFileOutput() {
	common.EnsurePathOrDie(DefaultLogFile, common.DefaultDirectoryMod)
	mod := os.O_CREATE | os.O_APPEND | os.O_WRONLY
	file, err := os.OpenFile(DefaultLogFile, mod, common.DefaultFileMod)
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
