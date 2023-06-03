package template

import (
	"github.com/thediveo/enumflag/v2"
)

const (
	yamlFormat Format = "yaml"
	jsonFormat Format = "json"
)

type formatFlag enumflag.Flag

const (
	yamlFlag formatFlag = iota
	jsonFlag
)

var formats = map[formatFlag]Format{
	yamlFlag: yamlFormat,
	jsonFlag: jsonFormat,
}

func (f formatFlag) value() Format {
	return formats[f]
}

func toFormatIds() map[formatFlag][]string {
	var formatIds = make(map[formatFlag][]string)
	for flag, format := range formats {
		formatIds[flag] = []string{string(format)}
	}
	return formatIds
}
