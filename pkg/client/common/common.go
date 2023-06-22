package common

import (
	"strings"
)

func DefaultShell(command string) string {
	if shellCmd := strings.TrimSpace(command); shellCmd != "" {
		return shellCmd
	} else {
		return "/bin/bash"
	}
}
