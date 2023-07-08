package common

import (
	"strings"

	"github.com/pkg/errors"

	"github.com/hckops/hckctl/pkg/client/kubernetes"
)

type ResourceSize uint

const (
	ExtraSmall ResourceSize = iota
	Small
	Medium
	Large
	ExtraLarge
)

var resourceSizes = map[ResourceSize]string{
	ExtraSmall: "XS",
	Small:      "S",
	Medium:     "M",
	Large:      "L",
	ExtraLarge: "XL",
}

func (size ResourceSize) String() string {
	return resourceSizes[size]
}

func (size ResourceSize) ToKubeResource() *kubernetes.KubeResource {
	smallSize := &kubernetes.KubeResource{
		Memory: "512Mi",
		Cpu:    "500m",
	}
	// TODO define sizes
	switch size {
	case Small:
		return smallSize
	default:
		return smallSize
	}
}

func ExistResourceSize(str string) (ResourceSize, error) {
	for size, value := range resourceSizes {
		// case insensitive
		if strings.ToUpper(str) == value {
			return size, nil
		}
	}
	return Small, errors.New("invalid resource size")
}
