package model

import (
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/dchest/uniuri"

	"github.com/hckops/hckctl/pkg/util"
)

type BoxStreams struct {
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
	IsTty  bool // TODO tty false for tunnel only
}

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

// matches anything other than a letter, digit or underscore, equivalent to "[^a-zA-Z0-9_]"
var anyNonWordCharacterRegex = regexp.MustCompile(`\W+`)

func (box *BoxV1) SafeName() string {
	return anyNonWordCharacterRegex.ReplaceAllString(box.Image.Repository, "-")
}

func (box *BoxV1) GenerateFullName() string {
	// e.g. "box-organization-image-RANDOM"
	return fmt.Sprintf("box-%s-%s", box.SafeName(), strings.ToLower(uniuri.NewLen(5)))
}

func (box *BoxV1) HasPorts() bool {
	return len(box.Network.Ports) > 0
}

// TODO validate that ports are valid ints + verify schema type
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
	value, _ := util.ToJson(box)
	return value
}
