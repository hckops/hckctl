package flag

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/hckops/hckctl/pkg/common/model"
)

func TestValidateNetworkVpnFlag(t *testing.T) {
	assert.Nil(t, ValidateNetworkVpnFlag("", map[string]model.VpnNetworkInfo{}))
	assert.Nil(t, ValidateNetworkVpnFlag("default", map[string]model.VpnNetworkInfo{"default": {}}))

	errMissing := ValidateNetworkVpnFlag("default", map[string]model.VpnNetworkInfo{})
	assert.EqualError(t, errMissing, "vpn network [default] config not found")
}
