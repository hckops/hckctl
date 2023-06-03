package util

import (
	"fmt"
	"net"
	"strconv"
)

func GetLocalPort(port string) (string, error) {
	if err := verifyOpenPort(port); err == nil {
		return port, nil
	} else {
		p, errConv := strconv.Atoi(port)
		if errConv != nil {
			return "", fmt.Errorf("port %s is not a valid int: %v", port, errConv)
		}
		nextPort := strconv.Itoa(p + 1)
		// WARN fmt.Printf("port %s is not available, attempt %s", port, nextPort)

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
