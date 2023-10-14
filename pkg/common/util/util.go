package util

import (
	"strings"
)

func DefaultShellCommand(command string) []string {
	if shellCmd := strings.TrimSpace(command); shellCmd != "" {
		return []string{shellCmd}
	} else {
		return []string{"/bin/bash"}
	}
}
