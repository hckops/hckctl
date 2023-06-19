package config

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/hckops/hckctl/pkg/client/kubernetes"
)

func TestResourceSizes(t *testing.T) {
	assert.Equal(t, 5, len(resourceSizes))
	assert.Equal(t, "XS", resourceSizes[extraSmall])
	assert.Equal(t, "S", resourceSizes[small])
	assert.Equal(t, "M", resourceSizes[medium])
	assert.Equal(t, "L", resourceSizes[large])
	assert.Equal(t, "XL", resourceSizes[extraLarge])
}

func TestToKubeResource(t *testing.T) {
	kubeResource := &kubernetes.KubeResource{
		Memory: "512Mi",
		Cpu:    "500m",
	}

	assert.Equal(t, kubeResource, extraSmall.toKubeResource())
	assert.Equal(t, kubeResource, small.toKubeResource())
	assert.Equal(t, kubeResource, medium.toKubeResource())
	assert.Equal(t, kubeResource, large.toKubeResource())
	assert.Equal(t, kubeResource, extraLarge.toKubeResource())
}

func TestExistResourceSize(t *testing.T) {
	size, err := existResourceSize("s")
	assert.NoError(t, err)
	assert.Equal(t, small, size)

	_, err = existResourceSize("abc")
	assert.EqualError(t, err, "invalid resource size")
}
