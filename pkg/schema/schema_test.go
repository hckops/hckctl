package schema

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProviderFlag(t *testing.T) {
	assert.Equal(t, 8, len(kinds))
	assert.Equal(t, "config/v1", KindConfigV1.String())
	assert.Equal(t, "api/v1", KindApiV1.String())
	assert.Equal(t, "sidecar/v1", KindSidecarV1.String())
	assert.Equal(t, "box/v1", KindBoxV1.String())
	assert.Equal(t, "lab/v1", KindLabV1.String())
	assert.Equal(t, "task/v1", KindTaskV1.String())
	assert.Equal(t, "flow/v1", KindFlowV1.String())
	assert.Equal(t, "dump/v1", KindDumpV1.String())
}
