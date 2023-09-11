package model

import (
	"fmt"
	"os"
	"strings"

	"golang.org/x/exp/slices"

	boxModel "github.com/hckops/hckctl/pkg/box/model"
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
func (box *LabBox) Expand(inputs map[string]string) (*LabBox, error) {

	// TODO optional
	if alias, err := expand(box.Alias, inputs); err != nil {
		return nil, err
	} else {
		box.Alias = alias
	}

	// TODO optional
	if vpn, err := expand(box.Vpn, inputs); err != nil {
		return nil, err
	} else {
		box.Vpn = vpn
	}

	for i, e := range box.Template.Env {
		if key, value, err := util.SplitKeyValue(e); err == nil {
			// ignore errors
			if env, err := expand(value, inputs); err == nil {
				box.Template.Env[i] = fmt.Sprintf("%s=%s", key, env)
			}
		}
	}

	return box, nil
}

func expand(raw string, inputs map[string]string) (string, error) {
	// reserved keyword
	const separator = ":"
	var err error
	expanded := os.Expand(raw, func(value string) string {

		// empty value
		if strings.TrimSpace(value) == "" {
			return ""
		}

		// optional field
		items := strings.Split(value, separator)
		if len(items) == 2 {
			// handle keywords
			switch items[1] {
			case "random":
				return util.RandomAlphanumeric(10)
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
