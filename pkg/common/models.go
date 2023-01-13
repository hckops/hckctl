package common

import (
	"fmt"

	"github.com/dchest/uniuri"
)

type BoxV1 struct {
	Kind  string
	Name  string
	Tags  []string
	Image struct {
		Repository string
		Version    string
	}
}

func (box *BoxV1) ImageName() string {
	var version string
	if box.Image.Version == "" {
		version = "latest"
	} else {
		version = box.Image.Version
	}
	return fmt.Sprintf("%s:%s", box.Image.Repository, version)
}

func (box *BoxV1) GenerateName() string {
	return fmt.Sprintf("%s-%s", box.Name, uniuri.NewLen(5))
}
