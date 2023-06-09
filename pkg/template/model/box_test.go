package model

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testBox = &BoxV1{
	Kind: "box/v1",
	Name: "my-test",
	Tags: []string{"my-test"},
	Image: struct {
		Repository string
		Version    string
	}{
		Repository: "hckops/my-test",
	},
	Shell: "/bin/bash",
	Network: struct{ Ports []string }{Ports: []string{
		"aaa:123",
		"bbb:456:789",
	}},
}

func TestGenerateName(t *testing.T) {
	boxId := testBox.GenerateName()
	assert.True(t, strings.HasPrefix(boxId, "box-my-test-"))
	assert.Equal(t, 17, len(boxId))
}

func TestImageName(t *testing.T) {
	assert.Equal(t, "hckops/my-test:latest", testBox.ImageName())
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

func TestPretty(t *testing.T) {
	json := `{
  "Kind": "box/v1",
  "Name": "my-test",
  "Tags": [
    "my-test"
  ],
  "Image": {
    "Repository": "hckops/my-test",
    "Version": ""
  },
  "Shell": "/bin/bash",
  "Network": {
    "Ports": [
      "aaa:123",
      "bbb:456:789"
    ]
  }
}`
	assert.Equal(t, json, testBox.Pretty())
}
