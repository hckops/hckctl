package flag

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/hckops/hckctl/pkg/common/model"
)

func TestValidateNetworkVpnFlag(t *testing.T) {
	emptyVpn, emptyErr := ValidateNetworkVpnFlag("", map[string]model.NetworkVpnInfo{})
	assert.Nil(t, emptyVpn)
	assert.Nil(t, emptyErr)

	validVpn, validErr := ValidateNetworkVpnFlag("default", map[string]model.NetworkVpnInfo{"default": {
		Name:        "myDefault",
		LocalPath:   "myLocalPath",
		ConfigValue: "myConfigValue",
	}})
	assert.Equal(t, &model.NetworkVpnInfo{Name: "myDefault", LocalPath: "myLocalPath", ConfigValue: "myConfigValue"}, validVpn)
	assert.Nil(t, validErr)

	invalidVpn, invalidErr := ValidateNetworkVpnFlag("default", map[string]model.NetworkVpnInfo{})
	assert.Nil(t, invalidVpn)
	assert.EqualError(t, invalidErr, "vpn network [default] config not found")
}
