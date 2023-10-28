package model

import (
	"fmt"
	"strings"

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
	defaultSize := &kubernetes.KubeResource{
		Memory: "1024Mi",
		Cpu:    "1000m",
	}
	switch size {
	case ExtraSmall:
		return &kubernetes.KubeResource{
			Memory: "512Mi",
			Cpu:    "500m",
		}
	case Small:
		return defaultSize
	case Medium:
		return &kubernetes.KubeResource{
			Memory: "2Gi",
			Cpu:    "2000m",
		}
	case Large:
		return &kubernetes.KubeResource{
			Memory: "3Gi",
			Cpu:    "3000m",
		}
	case ExtraLarge:
		return &kubernetes.KubeResource{
			Memory: "4Gi",
			Cpu:    "4000m",
		}
	default:
		return defaultSize
	}
}

func ExistResourceSize(value string) (ResourceSize, error) {
	for size, str := range resourceSizes {
		// case insensitive
		if strings.ToUpper(value) == str {
			return size, nil
		}
	}
	return Small, fmt.Errorf("invalid resource size value=%s", value)
}
