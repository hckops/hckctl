package common

import (
	"os"
	"os/user"
	"path/filepath"

	"github.com/rs/zerolog/log"
)

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
			log.Fatal().Msgf("Unable to create dir %q %v", dir, err)
		}
	}
}
