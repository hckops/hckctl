package schema

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testBox = &BoxV1{
	Kind: "box/v1",
	Name: "mybox",
	Tags: []string{"my-test"},
	Image: struct {
		Repository string
		Version    string
	}{
		Repository: "hckops/box-mybox",
	},
	Network: struct{ Ports []string }{Ports: []string{
		"aaa:123",
		"bbb:456:789",
	}},
}

func TestImageName(t *testing.T) {
	assert.Equal(t, "hckops/box-mybox:latest", testBox.ImageName())
}

func TestImageVersion(t *testing.T) {
	// FIXME mutable test
	testBox.Image.Version = "myversion"
	assert.Equal(t, "myversion", testBox.ImageVersion())

	testBox.Image.Version = ""
	assert.Equal(t, "latest", testBox.ImageVersion())
}

func TestGenerateName(t *testing.T) {
	assert.True(t, strings.HasPrefix(testBox.GenerateName(), "box-mybox-"))
}

func TestSafeName(t *testing.T) {
	assert.Equal(t, "hckops-box-mybox", testBox.SafeName())
}

func TestGenerateFullName(t *testing.T) {
	name := testBox.GenerateFullName()
	prefix := "box-hckops-box-mybox-latest-"
	suffix := strings.ReplaceAll(name, prefix, "")

	assert.True(t, strings.HasPrefix(name, prefix))
	assert.True(t, len(suffix) == 5)
	assert.True(t, strings.ToLower(suffix) == suffix)
}

func TestHasPorts(t *testing.T) {
	assert.True(t, testBox.HasPorts())
}

func TestNetworkPorts(t *testing.T) {
	ports := []PortV1{
		{Alias: "aaa", Local: "123", Remote: "123"},
		{Alias: "bbb", Local: "456", Remote: "789"},
	}
	assert.Equal(t, ports, testBox.NetworkPorts())
}

func TestPretty(t *testing.T) {
	json := `{
  "Kind": "box/v1",
  "Name": "mybox",
  "Tags": [
    "my-test"
  ],
  "Image": {
    "Repository": "hckops/box-mybox",
    "Version": ""
  },
  "Network": {
    "Ports": [
      "aaa:123",
      "bbb:456:789"
    ]
  }
}`
	assert.Equal(t, json, testBox.Pretty())
}
