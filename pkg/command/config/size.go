package config

import (
	"strings"

	"github.com/pkg/errors"

	"github.com/hckops/hckctl/pkg/client/kubernetes"
)

type resourceSize uint

const (
	extraSmall resourceSize = iota
	small
	medium
	large
	extraLarge
)

var resourceSizes = map[resourceSize]string{
	extraSmall: "XS",
	small:      "S",
	medium:     "M",
	large:      "L",
	extraLarge: "XL",
}

func (size resourceSize) String() string {
	return resourceSizes[size]
}

func (size resourceSize) toKubeResource() *kubernetes.KubeResource {
	smallSize := &kubernetes.KubeResource{
		Memory: "512Mi",
		Cpu:    "500m",
	}
	// TODO define sizes
	switch size {
	case small:
		return smallSize
	default:
		return smallSize
	}
}

func existResourceSize(str string) (resourceSize, error) {
	for size, value := range resourceSizes {
		if strings.ToUpper(str) == value {
			return size, nil
		}
	}
	return small, errors.New("invalid resource size")
}
