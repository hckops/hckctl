package flag

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/hckops/hckctl/pkg/common/model"
)

func TestValidateSourceFlag(t *testing.T) {
	provider := CloudProviderFlag

	errRevision := ValidateSourceFlag(&provider, &SourceFlag{Revision: "invalid"})
	assert.EqualError(t, errRevision, "flag not supported: provider=cloud revision=invalid")

	errLocal := ValidateSourceFlag(&provider, &SourceFlag{Revision: "main", Local: true})
	assert.EqualError(t, errLocal, "flag not supported: provider=cloud local=true")
}

func TestValidateNetworkVpnFlag(t *testing.T) {
	assert.Nil(t, ValidateNetworkVpnFlag("default", map[string]model.VpnNetworkInfo{"default": {}}))

	errMissing := ValidateNetworkVpnFlag("default", map[string]model.VpnNetworkInfo{})
	assert.EqualError(t, errMissing, "vpn network [default] config not found")
}
