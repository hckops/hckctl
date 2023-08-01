package model

import (
	"fmt"
	"math"
	"strings"

	"github.com/dchest/uniuri"
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
	Network struct {
		Ports []string
	}
}

type BoxPort struct {
	Alias  string
	Local  string // TODO int ?
	Remote string // TODO int ?
	Public bool   // TODO not used, always false
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
	return fmt.Sprintf("%s%s-%s", BoxPrefixName, box.Name, strings.ToLower(uniuri.NewLen(5)))
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

func (box *BoxV1) NetworkPorts(includeVirtual bool) map[string]BoxPort {

	ports := map[string]BoxPort{}
	for _, portString := range box.Network.Ports {

		// name:local[:remote]
		values := strings.Split(portString, ":")

		var remote string
		if len(values) == 2 {
			// local == remote
			remote = values[1]
		} else {
			remote = values[2]
		}

		// ports are not validated i.e. valid number and range
		port := BoxPort{
			Alias:  values[0],
			Local:  values[1],
			Remote: remote,
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
	return maps.Values(box.NetworkPorts(includeVirtual))
}

func PortFormatPadding(ports []BoxPort) int {
	var max float64
	for _, port := range ports {
		max = math.Max(max, float64(len(port.Alias)))
	}
	return int(max)
}

func (box *BoxV1) Pretty() string {
	value, _ := util.EncodeJsonIndent(box)
	return value
}
