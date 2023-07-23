package ssh

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSshTunnelOptsNetwork(t *testing.T) {
	assert.Equal(t, "tcp", (&SshTunnelOpts{}).Network())
}

func TestSshTunnelOptsLocalAddress(t *testing.T) {
	opts := &SshTunnelOpts{
		LocalPort: "123",
	}
	assert.Equal(t, "0.0.0.0:123", opts.LocalAddress())
}

func TestSshTunnelOptsRemoteAddress(t *testing.T) {
	opts := &SshTunnelOpts{
		RemoteHost: "myHost",
		RemotePort: "456",
	}
	assert.Equal(t, "myHost:456", opts.RemoteAddress())
}
