package flag

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/hckops/hckctl/pkg/box/model"
)

func TestValidateTunnelFlag(t *testing.T) {
	err := ValidateTunnelFlag(model.Docker, &TunnelFlag{NoExec: true, NoTunnel: true})
	assert.EqualError(t, err, "flag not supported: provider=docker no-exec=true no-tunnel=true")

	errExec := ValidateTunnelFlag(model.Docker, &TunnelFlag{NoExec: true, NoTunnel: false})
	assert.EqualError(t, errExec, "flag not supported: provider=docker no-exec=true no-tunnel=false")

	errTunnel := ValidateTunnelFlag(model.Docker, &TunnelFlag{NoExec: false, NoTunnel: true})
	assert.EqualError(t, errTunnel, "flag not supported: provider=docker no-exec=false no-tunnel=true")
}
