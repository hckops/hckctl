package schema

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidBoxV1(t *testing.T) {
	data :=
		`{
			"kind": "box/v1",
			"name": "my-test",
			"tags": [
				"my-test"
			],
			"image": {
				"repository": "hckops/my-test",
				"version": ""
			},
			"network": {
				"ports": [
					"aaa:123",
					"bbb:456:789"
				]
			}
		}`
	assert.NoError(t, ValidateBoxV1(data))
}

func TestRequiredBoxV1(t *testing.T) {
	expected := fmt.Errorf("validation error: jsonschema: '' does not validate with https://schema.hckops.com/box-v1.json#/required: missing properties: 'kind', 'name', 'tags', 'image'")

	err := ValidateBoxV1("{}")
	assert.Error(t, err)
	assert.Equal(t, expected, err)
}

func TestNullBoxV1(t *testing.T) {
	expected := fmt.Errorf("validation error: jsonschema: '' does not validate with https://schema.hckops.com/box-v1.json#/type: expected object, but got null")

	err := ValidateBoxV1("")
	assert.Error(t, err)
	assert.Equal(t, expected, err)
}
