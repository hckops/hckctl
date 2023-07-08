package common

import (
	"regexp"
	"strings"
)

func DefaultShell(command string) string {
	if shellCmd := strings.TrimSpace(command); shellCmd != "" {
		return shellCmd
	} else {
		return "/bin/bash"
	}
}

// matches anything other than a letter, digit or underscore, equivalent to "[^a-zA-Z0-9_]"
var anyNonWordCharacterRegex = regexp.MustCompile(`\W+`)

func ToKebabCase(value string) string {
	return anyNonWordCharacterRegex.ReplaceAllString(value, "-")
}
