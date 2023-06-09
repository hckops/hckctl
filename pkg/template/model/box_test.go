package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var testBox = &BoxV1{
	Kind: "box/v1",
	Name: "my-box",
	Tags: []string{"my-test"},
	Image: struct {
		Repository string
		Version    string
	}{
		Repository: "hckops/my-box",
	},
	Shell: "/bin/bash",
	Network: struct{ Ports []string }{Ports: []string{
		"aaa:123",
		"bbb:456:789",
	}},
}

func TestImageName(t *testing.T) {
	assert.Equal(t, "hckops/my-box:latest", testBox.ImageName())
}

func TestImageVersion(t *testing.T) {
	testBox.Image.Version = "my-version"
	assert.Equal(t, "my-version", testBox.ImageVersion())

	testBox.Image.Version = ""
	assert.Equal(t, "latest", testBox.ImageVersion())
}

func TestHasPorts(t *testing.T) {
	assert.True(t, testBox.HasPorts())
}

func TestNetworkPorts(t *testing.T) {
	ports := []BoxPort{
		{Alias: "aaa", Local: "123", Remote: "123", Public: false},
		{Alias: "bbb", Local: "456", Remote: "789", Public: false},
	}
	assert.Equal(t, ports, testBox.NetworkPorts())
}
