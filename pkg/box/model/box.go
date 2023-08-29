package model

import (
	"fmt"
	"math"
	"strings"

	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"

	"github.com/hckops/hckctl/pkg/util"
)

const (
	BoxPrefixName        = "box-"
	BoxShellNone         = "none"     // distroless
	BoxPortNone          = "none"     // runtime only when tunnelling
	boxPrefixVirtualPort = "virtual-" // experimental cloud feature only
)

type BoxV1 struct {
	Kind  string
	Name  string
	Tags  []string
	Image struct {
		Repository string
		Version    string
	}
	Shell   string
	Env     []string
	Network struct {
		Ports []string
	}
}

type BoxPort struct {
	Alias  string
	Remote string // TODO int ?
	Local  string // TODO int ?
	Public bool   // TODO
}

func SortPorts(ports []BoxPort) []BoxPort {
	sorted := ports
	// ascending order
	slices.SortFunc(sorted, func(a, b BoxPort) bool {
		return slices.IsSorted([]string{a.Remote, b.Remote})
	})
	return sorted
}

type BoxEnv struct {
	Key   string
	Value string
}

func SortEnv(env []BoxEnv) []BoxEnv {
	sorted := env
	slices.SortFunc(sorted, func(a, b BoxEnv) bool {
		return slices.IsSorted([]string{a.Key, b.Key})
	})
	return sorted
}

func (box *BoxV1) GenerateName() string {
	return fmt.Sprintf("%s%s-%s", BoxPrefixName, box.Name, util.RandomAlphanumeric(5))
}

// ToBoxTemplateName returns the strictly validated template name, or the original trimmed name
func ToBoxTemplateName(boxName string) string {
	trimmed := strings.TrimSpace(boxName)
	if trimmed == "" {
		return trimmed
	}
	if strings.Count(trimmed, "-") < 2 {
		return trimmed
	}

	// removes prefix and suffix
	values := strings.Split(boxName, "-")
	prefix := values[0]
	names := values[1 : len(values)-1]
	name := strings.Join(names, "-")
	suffix := values[len(values)-1]

	if len(prefix) > 0 && len(name) > 0 && len(suffix) == 5 {
		return name
	} else {
		return trimmed
	}
}

func (box *BoxV1) ImageName() string {
	return fmt.Sprintf("%s:%s", box.Image.Repository, box.ImageVersion())
}

func (box *BoxV1) ImageVersion() string {
	var version string
	if box.Image.Version == "" {
		version = "latest"
	} else {
		version = box.Image.Version
	}
	return version
}

func (box *BoxV1) HasPorts() bool {
	return len(box.Network.Ports) > 0
}

// TODO return error validation?
func (box *BoxV1) NetworkPorts(includeVirtual bool) map[string]BoxPort {

	ports := map[string]BoxPort{}
	for _, portString := range box.Network.Ports {

		// name:remote[:local]
		values := strings.Split(portString, ":")

		var local string
		if len(values) == 2 {
			// local == remote
			local = values[1]
		} else if len(values) == 3 {
			local = values[2]
		} else {
			// silently ignore
			continue
		}

		// ports are not validated i.e. valid number and range
		port := BoxPort{
			Alias:  values[0],
			Remote: values[1],
			Local:  local,
			Public: false,
		}

		// by default ignore virtual-* ports
		if !strings.HasPrefix(port.Alias, boxPrefixVirtualPort) || includeVirtual {
			// remote is always unique
			ports[port.Remote] = port
		}
	}

	return ports
}

func (box *BoxV1) NetworkPortValues(includeVirtual bool) []BoxPort {
	return SortPorts(maps.Values(box.NetworkPorts(includeVirtual)))
}

func PortFormatPadding(ports []BoxPort) int {
	var max float64
	for _, port := range ports {
		max = math.Max(max, float64(len(port.Alias)))
	}
	return int(max)
}

func (box *BoxV1) EnvironmentVariables() map[string]BoxEnv {
	// TODO return error validation?
	return ToEnvironmentVariables(box.Env)
}

func ToEnvironmentVariables(values []string) map[string]BoxEnv {
	envs := map[string]BoxEnv{}
	for _, e := range values {

		// silently ignore errors
		if key, value, err := util.SplitKeyValue(e); err == nil {
			envs[key] = BoxEnv{
				Key:   key,
				Value: value,
			}
		}
	}
	return envs
}

func (box *BoxV1) EnvironmentVariableValues() []BoxEnv {
	return SortEnv(maps.Values(box.EnvironmentVariables()))
}

func (box *BoxV1) Pretty() string {
	value, _ := util.EncodeJsonIndent(box)
	return value
}
