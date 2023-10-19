package model

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	commonModel "github.com/hckops/hckctl/pkg/common/model"
)

var testBox = &BoxV1{
	Kind: "box/v1",
	Name: "my-name",
	Tags: []string{"my-tag"},
	Image: commonModel.Image{
		Repository: "hckops/my-image",
	},
	Shell: "/bin/bash",
	Env: []string{
		"TTYD_USERNAME=username",
		"TTYD_PASSWORD=password",
	},
	Network: struct{ Ports []string }{Ports: []string{
		"foo",
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

func TestMainContainerName(t *testing.T) {
	assert.Equal(t, "hckops-my-image", testBox.MainContainerName())
}

func TestHasPorts(t *testing.T) {
	assert.True(t, testBox.HasPorts())
}

func TestNetworkPorts(t *testing.T) {
	ports := map[string]BoxPort{
		"123": {Alias: "aaa", Remote: "123", Local: "123", Public: false},
		"456": {Alias: "bbb", Remote: "456", Local: "789", Public: false},
	}
	assert.Equal(t, ports, testBox.NetworkPorts(false))
}

func TestNetworkPortsInvalid(t *testing.T) {
	var testBox = &BoxV1{
		Network: struct{ Ports []string }{Ports: []string{
			"foo",
			"foo:bar:bizz:buzz",
		}},
	}
	assert.Equal(t, map[string]BoxPort{}, testBox.NetworkPorts(false))
}

func TestNetworkPortsUnique(t *testing.T) {
	var testBox = &BoxV1{
		Network: struct{ Ports []string }{Ports: []string{
			"foo:123",
			"bar:123:456",
		}},
	}
	ports := map[string]BoxPort{
		"123": {Alias: "bar", Remote: "123", Local: "456", Public: false},
	}
	assert.Equal(t, ports, testBox.NetworkPorts(false))
}

func TestNetworkPortsIncludeVirtual(t *testing.T) {
	ports := map[string]BoxPort{
		"123": {Alias: "aaa", Remote: "123", Local: "123", Public: false},
		"456": {Alias: "bbb", Remote: "456", Local: "789", Public: false},
		"321": {Alias: "virtual-ccc", Remote: "321", Local: "321", Public: false},
	}
	assert.Equal(t, ports, testBox.NetworkPorts(true))
}

func TestNetworkPortValues(t *testing.T) {
	ports := []BoxPort{
		{Alias: "aaa", Remote: "123", Local: "123", Public: false},
		{Alias: "bbb", Remote: "456", Local: "789", Public: false},
	}
	assert.Equal(t, ports, testBox.NetworkPortValues(false))
}

func TestPortFormatPadding(t *testing.T) {
	ports := []BoxPort{
		{Alias: "aaaa"},
		{Alias: "bbbbbbbbbb"},
		{Alias: "cccccc"},
	}
	assert.Equal(t, 10, PortFormatPadding(ports))
}

func TestEnvironmentVariables(t *testing.T) {
	env := map[string]BoxEnv{
		"TTYD_USERNAME": {Key: "TTYD_USERNAME", Value: "username"},
		"TTYD_PASSWORD": {Key: "TTYD_PASSWORD", Value: "password"},
	}
	assert.Equal(t, env, testBox.EnvironmentVariables())
}

func TestEnvironmentVariablesInvalid(t *testing.T) {
	var testBox = &BoxV1{
		Env: []string{
			"foo",
			"=no_key",
			"no_value=",
			"ok==?=",
		},
	}
	env := map[string]BoxEnv{
		"ok": {Key: "ok", Value: "=?="},
	}
	assert.Equal(t, env, testBox.EnvironmentVariables())
}

func TestEnvironmentVariablesUnique(t *testing.T) {
	var testBox = &BoxV1{
		Env: []string{
			"TTYD_USERNAME=first",
			"TTYD_USERNAME=last",
		},
	}
	env := map[string]BoxEnv{
		"TTYD_USERNAME": {Key: "TTYD_USERNAME", Value: "last"},
	}
	assert.Equal(t, env, testBox.EnvironmentVariables())
}

func TestEnvironmentVariableValues(t *testing.T) {
	ports := []BoxEnv{
		{Key: "TTYD_PASSWORD", Value: "password"},
		{Key: "TTYD_USERNAME", Value: "username"},
	}
	assert.Equal(t, ports, testBox.EnvironmentVariableValues())
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
  "Env": [
    "TTYD_USERNAME=username",
    "TTYD_PASSWORD=password"
  ],
  "Network": {
    "Ports": [
      "foo",
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
