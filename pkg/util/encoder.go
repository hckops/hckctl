package util

import (
	"encoding/json"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

func EncodeJson(data interface{}) (string, error) {
	bytes, err := json.Marshal(data)
	return string(bytes), errors.Wrap(err, "error encoding json")
}

func EncodeJsonIndent(data interface{}) (string, error) {
	bytes, err := json.MarshalIndent(data, "", "  ")
	return string(bytes), errors.Wrap(err, "error encoding json")
}

func EncodeYaml(data interface{}) (string, error) {
	// v2 prints 2 spaces
	bytes, err := yaml.Marshal(data)
	return string(bytes), errors.Wrap(err, "error encoding yaml")
}
