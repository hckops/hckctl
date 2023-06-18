package util

import (
	"encoding/json"
	"regexp"

	"gopkg.in/yaml.v2"
)

func EncodeJson(data interface{}) (string, error) {
	bytes, err := json.Marshal(data)
	return string(bytes), err
}

func EncodeJsonIndent(data interface{}) (string, error) {
	bytes, err := json.MarshalIndent(data, "", "  ")
	return string(bytes), err
}

func EncodeYaml(data interface{}) (string, error) {
	// v2 prints 2 spaces
	bytes, err := yaml.Marshal(data)
	return string(bytes), err
}

// matches anything other than a letter, digit or underscore, equivalent to "[^a-zA-Z0-9_]"
var anyNonWordCharacterRegex = regexp.MustCompile(`\W+`)

func ToKebabCase(value string) string {
	return anyNonWordCharacterRegex.ReplaceAllString(value, "-")
}
