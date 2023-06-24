package v1

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPingRequest(t *testing.T) {
	request := NewPingRequest()
	value := `{"kind":"api/v1","name":"hck-ping","body":{}}`

	testRequest[PingBody](t, request, value)
}

func TestBoxCreateRequest(t *testing.T) {
	request := NewBoxCreateRequest("alpine")
	value := `{"kind":"api/v1","name":"hck-box-create","body":{"templateName":"alpine"}}`

	testRequest[BoxCreateBody](t, request, value)
}

func TestBoxDeleteRequest(t *testing.T) {
	request := NewBoxDeleteRequest("box-alpine-12345")
	value := `{"kind":"api/v1","name":"hck-box-delete","body":{"name":"box-alpine-12345"}}`

	testRequest[BoxDeleteBody](t, request, value)
}

func TestBoxListRequest(t *testing.T) {
	request := NewBoxListRequest()
	value := `{"kind":"api/v1","name":"hck-box-list","body":{}}`

	testRequest[BoxListBody](t, request, value)
}

func testRequest[T body](t *testing.T, request *Request[T], value string) {
	jsonString, err := request.Encode()
	assert.NoError(t, err)
	assert.Equal(t, value, jsonString)

	result, err := Decode[T](value)
	assert.NoError(t, err)
	assert.Equal(t, request, result)
}
