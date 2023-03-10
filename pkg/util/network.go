package util

import (
	"fmt"
	"net"
	"strconv"

	"github.com/rs/zerolog/log"
)

func GetLocalPort(port string) string {
	if err := verifyOpenPort(port); err == nil {
		return port
	} else {
		p, errConv := strconv.Atoi(port)
		if errConv != nil {
			log.Fatal().Err(errConv).Msgf("port %s is not a valid int", port)
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
