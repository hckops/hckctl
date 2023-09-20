package model

import (
	"fmt"

	"golang.org/x/exp/slices"

	boxModel "github.com/hckops/hckctl/pkg/box/model"
	commonModel "github.com/hckops/hckctl/pkg/common/model"
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
	Template BoxTemplate
	Size     string
	Vpn      string
	Ports    []string // cloud only
	Dumps    []string // cloud only
}

type BoxTemplate struct {
	Name string
	Env  []string
}

func (t *BoxTemplate) Merge(original *boxModel.BoxV1) *boxModel.BoxV1 {

	override := boxModel.ToEnvironmentVariables(t.Env)

	var envs []string
	for key, originalEnv := range original.EnvironmentVariables() {
		if overrideEnv, ok := override[key]; ok {
			envs = append(envs, fmt.Sprintf("%s=%s", key, overrideEnv.Value))
		} else {
			envs = append(envs, fmt.Sprintf("%s=%s", key, originalEnv.Value))
		}
	}
	slices.Sort(envs)

	original.Env = envs
	return original
}

func (lab *LabV1) Pretty() string {
	value, _ := util.EncodeJsonIndent(lab)
	return value
}

// TODO use reflection to expand all fields
func (box *LabBox) Expand(parameters commonModel.Parameters) (*LabBox, error) {

	// TODO optional
	if alias, err := util.Expand(box.Alias, parameters); err != nil {
		return nil, err
	} else {
		box.Alias = alias
	}

	// TODO optional
	if vpn, err := util.Expand(box.Vpn, parameters); err != nil {
		return nil, err
	} else {
		box.Vpn = vpn
	}

	for i, e := range box.Template.Env {
		if key, value, err := util.SplitKeyValue(e); err == nil {
			// ignore errors
			if env, err := util.Expand(value, parameters); err == nil {
				box.Template.Env[i] = fmt.Sprintf("%s=%s", key, env)
			}
		}
	}

	return box, nil
}
