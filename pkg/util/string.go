package util

import (
	"fmt"
	"os"
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

// TODO refactor test in box and task
func Expand(raw string, inputs map[string]string) (string, error) {
	// reserved keyword
	const separator = ":"
	var err error
	expanded := os.Expand(raw, func(value string) string {

		// empty value
		if strings.TrimSpace(value) == "" {
			return ""
		}

		// optional field
		items := strings.SplitN(value, separator, 2)
		if len(items) == 2 {
			// handle keywords
			switch items[1] {
			case "random":
				return RandomAlphanumeric(10)
			}

			key := items[0]
			if input, ok := inputs[key]; !ok {
				// default
				return items[1]
			} else {
				// input
				return input
			}
		}

		// required field
		if raw == fmt.Sprintf("$%s", value) || raw == fmt.Sprintf("${%s}", value) {
			if input, ok := inputs[value]; !ok {
				err = fmt.Errorf("%s required", value)
				return ""
			} else {
				return input
			}
		}

		err = fmt.Errorf("%s unexpected error", value)
		return ""
	})

	return expanded, err
}
