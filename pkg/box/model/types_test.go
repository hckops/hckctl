package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBoxProviderValues(t *testing.T) {
	assert.Equal(t, 4, len(providerValue))

	assert.Equal(t, "docker", providerValue[0])
	assert.Equal(t, "docker", Docker.String())

	assert.Equal(t, "kube", providerValue[1])
	assert.Equal(t, "kube", Kubernetes.String())

	assert.Equal(t, "argo-cd", providerValue[2])
	assert.Equal(t, "argo-cd", ArgoCd.String())

	assert.Equal(t, "cloud", providerValue[3])
	assert.Equal(t, "cloud", Cloud.String())
}
