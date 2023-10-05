package docker

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildVpnSidecarName(t *testing.T) {
	expected := "sidecar-vpn-12345"
	assert.Equal(t, expected, buildVpnSidecarName("aaa-bbb-ccc-ddd-12345"))
}
