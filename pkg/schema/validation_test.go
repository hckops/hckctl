package schema

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseValidBoxV1(t *testing.T) {
	data :=
		`{
			"kind": "box/v1",
			"name": "mybox",
			"tags": [
				"my-test"
			],
			"image": {
				"repository": "hckops/box-mybox",
				"version": ""
			},
			"network": {
				"ports": [
					"aaa:123",
					"bbb:456:789"
				]
			}
		}`
	expected := &BoxV1{
		Kind: "box/v1",
		Name: "mybox",
		Tags: []string{"my-test"},
		Image: struct {
			Repository string
			Version    string
		}{
			Repository: "hckops/box-mybox",
		},
		Network: struct{ Ports []string }{Ports: []string{
			"aaa:123",
			"bbb:456:789",
		}},
	}

	result, err := ParseValidBoxV1(data)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestParseInvalidBoxV1(t *testing.T) {
	expectedError := fmt.Errorf("validation error: jsonschema: '' does not validate with https://schema.hckops.com/box-v1.json#/type: expected object, but got null")

	result, err := ParseValidBoxV1("")
	if assert.Error(t, err) {
		assert.Equal(t, expectedError, err)
	}
	assert.Nil(t, result)
}
