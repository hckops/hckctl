package common

import (
	"encoding/json"
	"os"
	"os/user"
	"path/filepath"

	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v2"
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

func ToJson(data interface{}) (string, error) {
	bytes, err := json.MarshalIndent(data, "", "  ")
	return string(bytes), err
}

func ToYaml(data interface{}) (string, error) {
	// v2 prints 2 spaces
	bytes, err := yaml.Marshal(data)
	return string(bytes), err
}
