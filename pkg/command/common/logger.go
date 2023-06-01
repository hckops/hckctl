package common

import (
	"github.com/rs/zerolog/log"
)

const (
	LogLevelFlag = "log-level"
)

func InitFileLogger(global *GlobalCmdOptions, config *LogConfig) error {
	log.Info().Msgf("opts LEVEL=%s", global.LogLevel)
	log.Info().Msgf("config LEVEL=%s", config.Level)
	return nil
}
