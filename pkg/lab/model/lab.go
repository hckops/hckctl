package model

import (
	"github.com/hckops/hckctl/pkg/util"
)

type LabV1 struct {
	Kind string
	Name string
	Tags []string
	Box  LabBox
}

type LabBox struct {
	Alias    string
	Template string
	Env      []string
	Size     string
	Vpn      string
	Ports    []string // cloud only
	Dumps    []string // cloud only
}

func (box *LabV1) Pretty() string {
	value, _ := util.EncodeJsonIndent(box)
	return value
}
