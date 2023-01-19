package cmd

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	DefaultLogLevel = "info"
)

func InitLogger(flags *Flags) {
	setTimestamp()
	setLevel(parseLevel(flags.LogLevel))
	setFormat()
	setContext()
}

func setTimestamp() {
	zerolog.TimestampFunc = func() time.Time {
		return time.Now().UTC()
	}
}

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

func setFormat() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})
}

func setContext() {
	log.Logger = log.With().Caller().Str("source", "hckctl").Logger()
}
