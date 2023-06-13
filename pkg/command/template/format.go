package template

import (
	"github.com/thediveo/enumflag/v2"
)

type formatFlag enumflag.Flag

const (
	yamlFlag formatFlag = iota
	jsonFlag
)

var formatIds = map[formatFlag][]string{
	yamlFlag: {"yaml", "yml"},
	jsonFlag: {"json"},
}

func (f formatFlag) String() string {
	return formatIds[f][0]
}

func formatValues() []string {
	var values []string
	for _, formatId := range formatIds {
		for _, format := range formatId {
			values = append(values, format)
		}
	}
	return values
}
