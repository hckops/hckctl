package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBoxProvider(t *testing.T) {
	assert.Equal(t, "docker", Docker.String())
	assert.Equal(t, "kube", Kubernetes.String())
	assert.Equal(t, "cloud", Cloud.String())
}
