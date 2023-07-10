package model

import (
	"strings"

	"github.com/pkg/errors"

	"github.com/hckops/hckctl/pkg/client/kubernetes"
)

// TODO comparable e.g. "M > S"

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
	defaultSize := &kubernetes.KubeResource{
		Memory: "1024Mi",
		Cpu:    "1000m",
	}
	// TODO define sizes
	switch size {
	case ExtraSmall:
		return &kubernetes.KubeResource{
			Memory: "512Mi",
			Cpu:    "500m",
		}
	case Small:
		return defaultSize
	default:
		return defaultSize
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
