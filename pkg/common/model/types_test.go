package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestImageName(t *testing.T) {
	image := Image{
		Repository: "hckops/my-image",
	}
	assert.Equal(t, "hckops/my-image:latest", image.Name())
}

func TestImageVersion(t *testing.T) {
	image := Image{
		Repository: "hckops/my-image",
	}
	assert.Equal(t, "latest", image.ResolveVersion())

	image.Version = "my-version"
	assert.Equal(t, "my-version", image.ResolveVersion())
}
