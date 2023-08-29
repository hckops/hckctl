package model

import (
	"fmt"
	"os"
	"strings"

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

func (lab *LabV1) Pretty() string {
	value, _ := util.EncodeJsonIndent(lab)
	return value
}

// TODO use reflection to expand all fields
func (lab *LabV1) ExpandBox(inputs map[string]string) (*LabV1, error) {

	if alias, err := expand(lab.Box.Alias, inputs); err != nil {
		return nil, err
	} else {
		lab.Box.Alias = alias
	}

	if vpn, err := expand(lab.Box.Vpn, inputs); err != nil {
		return nil, err
	} else {
		lab.Box.Vpn = vpn
	}

	for i, e := range lab.Box.Env {
		items := strings.Split(e, "=")
		if len(items) == 2 {
			// ignore errors
			if env, err := expand(items[1], inputs); err == nil {
				lab.Box.Env[i] = fmt.Sprintf("%s=%s", items[0], env)
			}
		}
	}

	return lab, nil
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
