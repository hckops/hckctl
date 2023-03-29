package config

import (
	"os"
	"os/user"
	"path/filepath"

	"github.com/adrg/xdg"
	"github.com/rs/zerolog/log"
)

// returns configs home directory
func ConfigHome() string {
	if env := os.Getenv(ConfigNameEnv); env != "" {
		return env
	}

	// https://xdgbasedirectoryspecification.com
	xdgHome, err := xdg.ConfigFile(ConfigDir)
	if err != nil {
		log.Fatal().Err(err).Msgf("unable to create configuration directory: %s", ConfigDir)
	}

	return xdgHome
}

func GetUserOrDie() string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal().Err(err).Msg("unable to retrieve current user")
	}
	return usr.Username
}

// makes sure a directory exists from the given path
func EnsurePathOrDie(path string, mod os.FileMode) {
	dir := filepath.Dir(path)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err = os.MkdirAll(dir, mod); err != nil {
			log.Fatal().Err(err).Msgf("Unable to create dir %q", dir)
		}
	}
}
