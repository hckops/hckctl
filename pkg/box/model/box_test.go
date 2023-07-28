package model

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testBox = &BoxV1{
	Kind: "box/v1",
	Name: "my-name",
	Tags: []string{"my-tag"},
	Image: struct {
		Repository string
		Version    string
	}{
		Repository: "hckops/my-image",
	},
	Shell: "/bin/bash",
	Network: struct{ Ports []string }{Ports: []string{
		"aaa:123",
		"bbb:456:789",
		"virtual-ccc:321",
	}},
}

func TestGenerateName(t *testing.T) {
	boxId := testBox.GenerateName()
	assert.True(t, strings.HasPrefix(boxId, "box-my-name-"))
	assert.Equal(t, 17, len(boxId))
}

func TestToBoxTemplateName(t *testing.T) {
	assert.Equal(t, "my-long-name-example", ToBoxTemplateName("box-my-long-name-example-12345"))
	assert.Equal(t, "b", ToBoxTemplateName("a-b-12345"))
}

func TestToBoxTemplateNameInvalid(t *testing.T) {
	assert.Equal(t, "", ToBoxTemplateName("  \n \t  "))
	assert.Equal(t, "a", ToBoxTemplateName("a"))
	assert.Equal(t, "-", ToBoxTemplateName("-"))
	assert.Equal(t, "--", ToBoxTemplateName("--"))
	assert.Equal(t, "a--", ToBoxTemplateName("a--"))
	assert.Equal(t, "a-b-", ToBoxTemplateName("a-b-"))
	assert.Equal(t, "a-b-c", ToBoxTemplateName("a-b-c"))
	assert.Equal(t, "a--12345", ToBoxTemplateName("a--12345"))
	assert.Equal(t, "-b-12345", ToBoxTemplateName("-b-12345"))
}

func TestImageName(t *testing.T) {
	assert.Equal(t, "hckops/my-image:latest", testBox.ImageName())
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
	assert.Equal(t, ports, testBox.NetworkPorts(false))
}

func TestNetworkPortsIncludeVirtual(t *testing.T) {
	ports := []BoxPort{
		{Alias: "aaa", Local: "123", Remote: "123", Public: false},
		{Alias: "bbb", Local: "456", Remote: "789", Public: false},
		{Alias: "virtual-ccc", Local: "321", Remote: "321", Public: false},
	}
	assert.Equal(t, ports, testBox.NetworkPorts(true))
}

func TestPortFormatPadding(t *testing.T) {
	ports := []BoxPort{
		{Alias: "aaaa"},
		{Alias: "bbbbbbbbbb"},
		{Alias: "cccccc"},
	}
	assert.Equal(t, 10, PortFormatPadding(ports))
}

func TestPretty(t *testing.T) {
	json := `{
  "Kind": "box/v1",
  "Name": "my-name",
  "Tags": [
    "my-tag"
  ],
  "Image": {
    "Repository": "hckops/my-image",
    "Version": ""
  },
  "Shell": "/bin/bash",
  "Network": {
    "Ports": [
      "aaa:123",
      "bbb:456:789",
      "virtual-ccc:321"
    ]
  }
}`
	assert.Equal(t, json, testBox.Pretty())
}

func TestSortPorts(t *testing.T) {
	ports := []BoxPort{
		{Remote: "remote-d"},
		{Remote: "remote-a"},
		{Remote: "remote-c"},
		{Remote: "remote-b"},
	}
	expected := []BoxPort{
		{Remote: "remote-a"},
		{Remote: "remote-b"},
		{Remote: "remote-c"},
		{Remote: "remote-d"},
	}
	assert.Equal(t, expected, SortPorts(ports))
}

func TestSortEnv(t *testing.T) {
	env := []BoxEnv{
		{Key: "MY_KEY_4"},
		{Key: "MY_KEY_1"},
		{Key: "MY_KEY_3"},
		{Key: "MY_KEY_2"},
	}
	expected := []BoxEnv{
		{Key: "MY_KEY_1"},
		{Key: "MY_KEY_2"},
		{Key: "MY_KEY_3"},
		{Key: "MY_KEY_4"},
	}
	assert.Equal(t, expected, SortEnv(env))
}
