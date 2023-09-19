package util

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/dchest/uniuri"
	"github.com/pkg/errors"
)

func IotaToValues[T comparable](kv map[T]string) []string {
	var values []string
	for _, v := range kv {
		values = append(values, v)
	}
	return values
}

func RandomAlphanumeric(size int) string {
	return strings.ToLower(uniuri.NewLen(size))
}

func SplitKeyValue(kv string) (string, string, error) {
	// split on first equal
	values := strings.Split(kv, "=")
	if len(values) >= 2 {
		// value might contains equals
		v := strings.TrimPrefix(kv, fmt.Sprintf("%s=", values[0]))

		key := strings.TrimSpace(values[0])
		if key == "" {
			return "", "", errors.New("empty key")
		}

		value := strings.TrimSpace(v)
		if value == "" {
			return "", "", errors.New("empty value")
		}

		return key, value, nil
	}
	return "", "", errors.New("invalid key-value pair")
}

// matches anything other than a letter, digit or underscore, equivalent to "[^a-zA-Z0-9_]"
var anyNonWordCharacterRegex = regexp.MustCompile(`\W+`)

func ToLowerKebabCase(value string) string {
	return anyNonWordCharacterRegex.ReplaceAllString(strings.ToLower(strings.TrimSpace(value)), "-")
}