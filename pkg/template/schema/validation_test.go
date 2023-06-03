package schema

//import (
//	"fmt"
//	"testing"
//
//	"github.com/stretchr/testify/assert"
//
//	"github.com/hckops/hckctl/pkg/template"
//)
//
//func TestParseValidBoxV1(t *testing.T) {
//	data :=
//		`{
//			"kind": "box/v1",
//			"name": "my-box",
//			"tags": [
//				"my-test"
//			],
//			"image": {
//				"repository": "hckops/my-box",
//				"version": ""
//			},
//			"network": {
//				"ports": [
//					"aaa:123",
//					"bbb:456:789"
//				]
//			}
//		}`
//	expected := &template.BoxV1{
//		Kind: "box/v1",
//		Name: "my-box",
//		Tags: []string{"my-test"},
//		Image: struct {
//			Repository string
//			Version    string
//		}{
//			Repository: "hckops/my-box",
//		},
//		Network: struct{ Ports []string }{Ports: []string{
//			"aaa:123",
//			"bbb:456:789",
//		}},
//	}
//
//	result, err := ParseValidBoxV1(data)
//	assert.NoError(t, err)
//	assert.Equal(t, expected, result)
//}
//
//func TestParseInvalidBoxV1(t *testing.T) {
//	expectedError := fmt.Errorf("validation error: jsonschema: '' does not validate with https://schema.hckops.com/box-v1.json#/type: expected object, but got null")
//
//	result, err := ParseValidBoxV1("")
//	if assert.Error(t, err) {
//		assert.Equal(t, expectedError, err)
//	}
//	assert.Nil(t, result)
//}
