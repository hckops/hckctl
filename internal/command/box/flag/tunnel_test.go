package flag

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/hckops/hckctl/pkg/box/model"
)

func TestValidateTunnelFlag(t *testing.T) {
	err := ValidateTunnelFlag(&TunnelFlag{NoExec: true, NoTunnel: true}, model.Docker)
	assert.EqualError(t, err, "flag not supported: provider=docker no-exec=true no-tunnel=true")

	errExec := ValidateTunnelFlag(&TunnelFlag{NoExec: true, NoTunnel: false}, model.Docker)
	assert.EqualError(t, errExec, "flag not supported: provider=docker no-exec=true no-tunnel=false")

	errTunnel := ValidateTunnelFlag(&TunnelFlag{NoExec: false, NoTunnel: true}, model.Docker)
	assert.EqualError(t, errTunnel, "flag not supported: provider=docker no-exec=false no-tunnel=true")
}
