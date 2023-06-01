package common

import (
	"github.com/rs/zerolog/log"
)

func SetupLogger(global *GlobalCmdOptions, config *LogConfig) error {
	log.Info().Msgf("opts LEVEL=%s", global.LogLevel)
	log.Info().Msgf("config LEVEL=%s", config.Level)
	return nil
}
