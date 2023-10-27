package model

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/hckops/hckctl/pkg/client/kubernetes"
)

func TestResourceSizes(t *testing.T) {
	assert.Equal(t, 5, len(resourceSizes))
	assert.Equal(t, "XS", resourceSizes[ExtraSmall])
	assert.Equal(t, "S", resourceSizes[Small])
	assert.Equal(t, "M", resourceSizes[Medium])
	assert.Equal(t, "L", resourceSizes[Large])
	assert.Equal(t, "XL", resourceSizes[ExtraLarge])
}

func TestToKubeResource(t *testing.T) {
	extraSmallSize := &kubernetes.KubeResource{
		Memory: "512Mi",
		Cpu:    "500m",
	}
	assert.Equal(t, extraSmallSize, ExtraSmall.ToKubeResource())

	smallSize := &kubernetes.KubeResource{
		Memory: "1024Mi",
		Cpu:    "1000m",
	}
	assert.Equal(t, smallSize, Small.ToKubeResource())

	mediumSize := &kubernetes.KubeResource{
		Memory: "2Gi",
		Cpu:    "2000m",
	}
	assert.Equal(t, mediumSize, Medium.ToKubeResource())

	largeSize := &kubernetes.KubeResource{
		Memory: "3Gi",
		Cpu:    "3000m",
	}
	assert.Equal(t, largeSize, Large.ToKubeResource())

	extraLargeSize := &kubernetes.KubeResource{
		Memory: "4Gi",
		Cpu:    "4000m",
	}
	assert.Equal(t, extraLargeSize, ExtraLarge.ToKubeResource())
}

func TestExistResourceSize(t *testing.T) {
	size, err := ExistResourceSize("s")
	assert.NoError(t, err)
	assert.Equal(t, Small, size)

	_, err = ExistResourceSize("abc")
	assert.EqualError(t, err, "invalid resource size value=abc")
}
