package flag

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/hckops/hckctl/internal/command/common/flag"
	"github.com/hckops/hckctl/pkg/box/model"
)

func TestValidateSourceFlag(t *testing.T) {
	errRevision := ValidateSourceFlag(model.Cloud, &flag.SourceFlag{Revision: "invalid"})
	assert.EqualError(t, errRevision, "flag not supported: provider=cloud revision=invalid")

	errLocal := ValidateSourceFlag(model.Cloud, &flag.SourceFlag{Revision: "main", Local: true})
	assert.EqualError(t, errLocal, "flag not supported: provider=cloud local=true")
}

func TestValidateTunnelFlag(t *testing.T) {
	err := ValidateTunnelFlag(model.Docker, &TunnelFlag{NoExec: true, NoTunnel: true})
	assert.EqualError(t, err, "flag not supported: provider=docker no-exec=true no-tunnel=true")

	errExec := ValidateTunnelFlag(model.Docker, &TunnelFlag{NoExec: true, NoTunnel: false})
	assert.EqualError(t, errExec, "flag not supported: provider=docker no-exec=true no-tunnel=false")

	errTunnel := ValidateTunnelFlag(model.Docker, &TunnelFlag{NoExec: false, NoTunnel: true})
	assert.EqualError(t, errTunnel, "flag not supported: provider=docker no-exec=false no-tunnel=true")
}
