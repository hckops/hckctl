package common

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"os/user"
	"path/filepath"
	"strconv"

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
			log.Fatal().Err(err).Msgf("Unable to create dir %q", dir)
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

func GetLocalPort(port string) string {
	if err := verifyOpenPort(port); err == nil {
		return port
	} else {
		p, err := strconv.Atoi(port)
		if err != nil {
			log.Fatal().Err(err).Msgf("port %s is not a valid int", port)
		}
		nextPort := strconv.Itoa(p + 1)
		log.Warn().Err(err).Msgf("port %s is not available, attempt %s", port, nextPort)

		return GetLocalPort(nextPort)
	}
}

func verifyOpenPort(port string) error {
	listener, err := net.Listen("tcp", fmt.Sprintf("[::]:%s", port))
	if err != nil {
		return fmt.Errorf("unable to listen on port %s: %v", port, err)
	}

	if err := listener.Close(); err != nil {
		return fmt.Errorf("failed to close port %s: %v", port, err)
	}

	return nil
}
