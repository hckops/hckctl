package event

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMethods(t *testing.T) {
	assert.Equal(t, 7, len(events))
	assert.Equal(t, "debug", LogDebug.String())
	assert.Equal(t, "info", LogInfo.String())
	assert.Equal(t, "warning", LogWarning.String())
	assert.Equal(t, "error", LogError.String())
	assert.Equal(t, "console", PrintConsole.String())
	assert.Equal(t, "update", LoaderUpdate.String())
	assert.Equal(t, "stop", LoaderStop.String())
}
