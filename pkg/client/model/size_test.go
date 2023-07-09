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
	kubeResource := &kubernetes.KubeResource{
		Memory: "512Mi",
		Cpu:    "500m",
	}

	assert.Equal(t, kubeResource, ExtraSmall.ToKubeResource())
	assert.Equal(t, kubeResource, Small.ToKubeResource())
	assert.Equal(t, kubeResource, Medium.ToKubeResource())
	assert.Equal(t, kubeResource, Large.ToKubeResource())
	assert.Equal(t, kubeResource, ExtraLarge.ToKubeResource())
}

func TestExistResourceSize(t *testing.T) {
	size, err := ExistResourceSize("s")
	assert.NoError(t, err)
	assert.Equal(t, Small, size)

	_, err = ExistResourceSize("abc")
	assert.EqualError(t, err, "invalid resource size")
}
