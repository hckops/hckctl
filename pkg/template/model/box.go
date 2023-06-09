package model

import (
	"fmt"
	"strings"

	"github.com/dchest/uniuri"
	"github.com/hckops/hckctl/pkg/util"
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
	Local  string
	Remote string
	Public bool // TODO not used, always false
}

func (box *BoxV1) GenerateName() string {
	return fmt.Sprintf("box-%s-%s", box.Name, strings.ToLower(uniuri.NewLen(5)))
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

func (box *BoxV1) NetworkPorts() []BoxPort {

	ports := make([]BoxPort, 0)
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

		// ports are not validated
		port := BoxPort{
			Alias:  values[0],
			Local:  values[1],
			Remote: remote,
			Public: false,
		}

		ports = append(ports, port)
	}

	return ports
}

func (box *BoxV1) Pretty() string {
	value, _ := util.EncodeJsonIndent(box)
	return value
}
