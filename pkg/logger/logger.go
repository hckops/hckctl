package logger

import (
	"io"
	"os"
	"strings"
	"time"

	"github.com/dchest/uniuri"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/hckops/hckctl/pkg/util"
)

func SetTimestamp() {
	zerolog.TimestampFunc = func() time.Time {
		return time.Now().UTC()
	}
}

func SetLevel(value string) {
	level := parseLevel(value)

	if level == zerolog.NoLevel {
		// default info
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	} else {
		zerolog.SetGlobalLevel(level)
	}
}

func parseLevel(value string) zerolog.Level {
	switch value {
	case DebugLogLevel.String():
		return zerolog.DebugLevel
	case InfoLogLevel.String():
		return zerolog.InfoLevel
	case WarningLogLevel.String():
		return zerolog.WarnLevel
	case ErrorLogLevel.String():
		return zerolog.ErrorLevel
	default:
		// silent fallback
		return zerolog.NoLevel
	}
}

func SetFormat(value string, out io.Writer) {
	format := parseFormat(value)

	// default is json
	if format == TextLogFormat {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: out, TimeFormat: time.RFC3339})
	}
}

func parseFormat(value string) LogFormat {
	switch value {
	case JsonLogFormat.String():
		return JsonLogFormat
	case TextLogFormat.String():
		return TextLogFormat
	default:
		// silent fallback
		return JsonLogFormat
	}
}

func SetContext(source string) {
	log.Logger = log.With().
		Caller().
		Stack().
		Str("source", source).
		Logger()
}

func SetSession() {
	log.Logger = log.With().
		Str("session", generateSession()).
		Logger()
}

func generateSession() string {
	return strings.ToLower(uniuri.NewLen(5))
}

func SetFileOutput(filePath string) (func() error, error) {

	if err := util.CreateBaseDir(filePath); err != nil {
		return nil, errors.Wrap(err, "error creating log dir")
	}

	const fileMod os.FileMode = 0600
	mod := os.O_CREATE | os.O_APPEND | os.O_WRONLY
	file, err := os.OpenFile(filePath, mod, fileMod)
	if err != nil {
		return nil, errors.Wrap(err, "error creating log file")
	}

	// default
	SetFormat(TextLogFormat.String(), file)

	return closeFileCallback(file), nil
}

func closeFileCallback(file *os.File) func() error {
	return func() error {
		if file != nil {
			log.Debug().Msg("closing log file")
			return file.Close()
		}
		log.Warn().Msg("log file already closed")
		return nil
	}
}
