package common

import (
	"os"

	"github.com/adrg/xdg"
	"github.com/rs/zerolog/log"
)

// returns configs home directory
func ConfigHome() string {
	if env := os.Getenv(ConfigNameEnv); env != "" {
		return env
	}

	xdgHome, err := xdg.ConfigFile(ConfigName)
	if err != nil {
		log.Fatal().Err(err).Msgf("unable to create configuration directory for %s", ConfigName)
	}

	return xdgHome
}
