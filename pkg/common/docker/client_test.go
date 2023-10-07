package docker

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildSidecarVpnName(t *testing.T) {
	expected := "sidecar-vpn-12345"
	assert.Equal(t, expected, buildSidecarVpnName("aaa-bbb-ccc-ddd-12345"))
}
