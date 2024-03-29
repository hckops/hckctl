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
			"name": "my-name",
			"tags": [
				"my-tag"
			],
			"image": {
				"repository": "hckops/my-image",
				"version": "latest"
			},
			"shell": "/bin/bash",
			"env": [
				"MY_KEY_1=my-value-1",
				"MY_KEY_2=my-value-2"
			],
			"network": {
				"ports": [
					"aaa:123",
					"bbb:456:789"
				]
			}
		}`
	assert.NoError(t, ValidateBoxV1(data))
}

func TestBoxRequired(t *testing.T) {
	data :=
		`{
			"kind": "box/v1",
			"name": "my-name",
			"tags": [
				"my-tag"
			],
			"image": {
				"repository": "hckops/my-image"
			}
		}`
	assert.NoError(t, ValidateBoxV1(data))
}

func TestBoxMissingRequired(t *testing.T) {
	err := ValidateBoxV1("{}")
	expected := fmt.Errorf("validation error: jsonschema: '' does not validate with https://schema.hckops.com/box-v1.json#/required: missing properties: 'kind', 'name', 'tags', 'image'")
	assert.Error(t, err)
	assert.Equal(t, expected, err)
}

func TestBoxNull(t *testing.T) {
	err := ValidateBoxV1("")
	expected := fmt.Errorf("validation error: jsonschema: '' does not validate with https://schema.hckops.com/box-v1.json#/type: expected object, but got null")
	assert.Error(t, err)
	assert.Equal(t, expected, err)
}

func TestBoxKind(t *testing.T) {
	data :=
		`{
			"kind": "foo",
			"name": "my-name",
			"tags": [
				"my-tag"
			],
			"image": {
				"repository": "hckops/my-image"
			}
		}`
	err := ValidateBoxV1(data)
	expected := fmt.Errorf("validation error: jsonschema: '/kind' does not validate with https://schema.hckops.com/box-v1.json#/properties/kind/const: value must be \"box/v1\"")
	assert.Equal(t, expected, err)
}

func TestBoxTagsMinItems(t *testing.T) {
	data :=
		`{
			"kind": "box/v1",
			"name": "my-name",
			"tags": [],
			"image": {
				"repository": "hckops/my-image"
			}
		}`
	err := ValidateBoxV1(data)
	expected := fmt.Errorf("validation error: jsonschema: '/tags' does not validate with https://schema.hckops.com/box-v1.json#/properties/tags/minItems: minimum 1 items required, but found 0 items")
	assert.Equal(t, expected, err)
}

func TestBoxTagsUniqueItems(t *testing.T) {
	data :=
		`{
			"kind": "box/v1",
			"name": "my-name",
			"tags": [
				"my-tag",
				"my-tag"
			],
			"image": {
				"repository": "hckops/my-image"
			}
		}`
	err := ValidateBoxV1(data)
	expected := fmt.Errorf("validation error: jsonschema: '/tags' does not validate with https://schema.hckops.com/box-v1.json#/properties/tags/uniqueItems: items at index 0 and 1 are equal")
	assert.Equal(t, expected, err)
}

func TestBoxEnvMinItems(t *testing.T) {
	data :=
		`{
			"kind": "box/v1",
			"name": "my-name",
			"tags": [
				"my-tag"
			],
			"image": {
				"repository": "hckops/my-image",
			},
			"env": []
		}`
	err := ValidateBoxV1(data)
	expected := fmt.Errorf("validation error: jsonschema: '/env' does not validate with https://schema.hckops.com/box-v1.json#/properties/env/minItems: minimum 1 items required, but found 0 items")
	assert.Equal(t, expected, err)
}

func TestBoxEnvUniqueItems(t *testing.T) {
	data :=
		`{
			"kind": "box/v1",
			"name": "my-name",
			"tags": [
				"my-tag"
			],
			"image": {
				"repository": "hckops/my-image",
			},
			"env": [
				"MY_KEY_1=my-value-1",
				"MY_KEY_1=my-value-1"
			],
		}`
	err := ValidateBoxV1(data)
	expected := fmt.Errorf("validation error: jsonschema: '/env' does not validate with https://schema.hckops.com/box-v1.json#/properties/env/uniqueItems: items at index 0 and 1 are equal")
	assert.Equal(t, expected, err)
}

func TestBoxNetworkPortRequired(t *testing.T) {
	data :=
		`{
			"kind": "box/v1",
			"name": "my-name",
			"tags": [
				"my-tag"
			],
			"image": {
				"repository": "hckops/my-image",
				"version": ""
			},
			"network": {}
		}`
	err := ValidateBoxV1(data)
	expected := fmt.Errorf("validation error: jsonschema: '/network' does not validate with https://schema.hckops.com/box-v1.json#/properties/network/required: missing properties: 'ports'")
	assert.Equal(t, expected, err)
}

// TODO bad validation, it should fail
func TestBoxNetworkPortMinItems(t *testing.T) {
	data :=
		`{
			"kind": "box/v1",
			"name": "my-name",
			"tags": [
				"my-tag"
			],
			"image": {
				"repository": "hckops/my-image",
				"version": ""
			},
			"network": {
				"ports": []
			}
		}`
	assert.NoError(t, ValidateBoxV1(data))
}

// TODO bad validation, it should fail
func TestBoxNetworkPortUniqueItems(t *testing.T) {
	data :=
		`{
			"kind": "box/v1",
			"name": "my-name",
			"tags": [
				"my-tag"
			],
			"image": {
				"repository": "hckops/my-image",
				"version": ""
			},
			"network": {
				"ports": [
					"aaa:123",
					"aaa:123"
				]
			}
		}`
	assert.NoError(t, ValidateBoxV1(data))
}

func TestValidLabV1(t *testing.T) {
	data :=
		`{
			"kind": "lab/v1",
			"name": "my-name",
			"tags": [
				"my-tag"
			],
			"box": {
				"alias": "my-alias",
				"template": {
					"name": "my-template",
					"env": [
						"MY_KEY_1=my-value-1",
						"MY_KEY_2=my-value-2"
					],
				},
				"size": "M",
				"vpn": "my-vpn",
				"ports": [
					"tty"
				],
				"dumps": [
					"seclists"
				]
			}
		}`
	assert.NoError(t, ValidateLabV1(data))
}

// TODO bad validation, box.template.name is required and it should fail
// TODO "box" is temporary required, change it to mutually exclusive
func TestLabRequired(t *testing.T) {
	data :=
		`{
			"kind": "lab/v1",
			"name": "my-name",
			"tags": [
				"my-tag"
			],
			"box": {
				"template": {
					"name": "my-template"
				},
				"size": "M"
			}
		}`
	assert.NoError(t, ValidateLabV1(data))
}

func TestValidDumpV1(t *testing.T) {
	data :=
		`{
			"kind": "dump/v1",
			"name": "my-name",
			"tags": [
				"my-tag"
			],
			"group": "my-group",
			"git": {
				"repository": "my-url",
				"branch": "my-branch"
			}
		}`
	assert.NoError(t, ValidateDumpV1(data))
}

func TestValidTaskV1(t *testing.T) {
	data :=
		`{
			"kind": "task/v1",
			"name": "my-name",
			"tags": [
				"my-tag"
			]
		}`
	assert.NoError(t, ValidateTaskV1(data))
}
