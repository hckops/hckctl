package template

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatFlag(t *testing.T) {
	assert.Equal(t, 2, len(formatIds))
	assert.Equal(t, []string{"yaml", "yml"}, formatIds[yamlFlag])
	assert.Equal(t, []string{"json"}, formatIds[jsonFlag])
}

func TestFormatValues(t *testing.T) {
	values := formatValues()
	expected := []string{"json", "yaml", "yml"}

	assert.Equal(t, 3, len(values))
	assert.Equal(t, expected, values)
}
