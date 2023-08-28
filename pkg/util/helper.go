package util

import (
	"strings"

	"github.com/dchest/uniuri"
)

func IotaToValues[T comparable](kv map[T]string) []string {
	var values []string
	for _, v := range kv {
		values = append(values, v)
	}
	return values
}

func Random(size int) string {
	return strings.ToLower(uniuri.NewLen(size))
}
