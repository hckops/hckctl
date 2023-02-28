package model

import (
	"fmt"
	"strings"

	"github.com/dchest/uniuri"

	"github.com/hckops/hckctl/internal/common"
)

type BoxV1 struct {
	Kind  string
	Name  string
	Tags  []string
	Image struct {
		Repository string
		Version    string // TODO RawVersion (internal)
	}
	Network struct {
		Ports []string
	}
}

// TODO verify schema validation
type PortV1 struct {
	Alias  string
	Local  string
	Remote string
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

func (box *BoxV1) GenerateName() string {
	return fmt.Sprintf("box-%s-%s", box.Name, strings.ToLower(uniuri.NewLen(5)))
}

func (box *BoxV1) NetworkPorts() []PortV1 {

	ports := make([]PortV1, 0)
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

		port := PortV1{
			Alias:  values[0],
			Local:  values[1],
			Remote: remote,
		}

		ports = append(ports, port)
	}

	return ports
}

func (box *BoxV1) Pretty() string {
	value, _ := common.ToJson(box)
	return value
}
