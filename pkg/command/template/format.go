package template

import (
	templateFormat "github.com/hckops/hckctl/pkg/template/source"
	"github.com/thediveo/enumflag/v2"
)

type formatFlag enumflag.Flag

const (
	yamlFlag formatFlag = iota
	jsonFlag
)

var formatIds = map[formatFlag][]string{
	yamlFlag: {templateFormat.YamlFormat.String(), "yml"},
	jsonFlag: {templateFormat.JsonFormat.String()},
}

func (f formatFlag) value() string {
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
