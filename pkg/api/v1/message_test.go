package v1

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	clientOrigin = "hckctl-0.0.0-os"
	serverOrigin = "hckadm-0.0.0-info"
)

var testBoxes = []string{"box-alpine-123", "box-alpine-456"}

func TestMethods(t *testing.T) {
	assert.Equal(t, 5, len(methods))
	assert.Equal(t, "hck-ping", MethodPing.String())
	assert.Equal(t, "hck-box-create", MethodBoxCreate.String())
	assert.Equal(t, "hck-box-exec", MethodBoxExec.String())
	assert.Equal(t, "hck-box-delete", MethodBoxDelete.String())
	assert.Equal(t, "hck-box-list", MethodBoxList.String())
}

func TestIsValidProtocol(t *testing.T) {
	_, errInvalidProtocol := IsValidProtocol("invalid")
	assert.EqualError(t, errInvalidProtocol, "invalid protocol")

	_, errInvalidMethod := IsValidProtocol("api/v1/todo")
	assert.EqualError(t, errInvalidMethod, "invalid method: method not found todo")

	method, err := IsValidProtocol("api/v1/hck-ping")
	assert.NoError(t, err)
	assert.Equal(t, "hck-ping", method)
}

func TestPingRequest(t *testing.T) {
	message := NewPingMessage(clientOrigin)
	value := `{"kind":"api/v1","origin":"hckctl-0.0.0-os","method":"hck-ping","body":{"value":"ping"}}`

	testMessage[PingBody](t, message, value)
}

func TestPingResponse(t *testing.T) {
	message := NewPongMessage(serverOrigin)
	value := `{"kind":"api/v1","origin":"hckadm-0.0.0-info","method":"hck-ping","body":{"value":"pong"}}`

	testMessage[PongBody](t, message, value)
}

func TestBoxCreateRequest(t *testing.T) {
	message := NewBoxCreateRequest(clientOrigin, "alpine", "s")
	value := `{"kind":"api/v1","origin":"hckctl-0.0.0-os","method":"hck-box-create","body":{"templateName":"alpine","size":"s"}}`

	testMessage[BoxCreateRequestBody](t, message, value)
}

func TestBoxCreateResponse(t *testing.T) {
	message := NewBoxCreateResponse(serverOrigin, testBoxes[0], "m")
	value := `{"kind":"api/v1","origin":"hckadm-0.0.0-info","method":"hck-box-create","body":{"name":"box-alpine-123","size":"m"}}`

	testMessage[BoxCreateResponseBody](t, message, value)
}

func TestBoxExecRequest(t *testing.T) {
	message := NewBoxExecRequest(clientOrigin, "alpine")
	value := `{"kind":"api/v1","origin":"hckctl-0.0.0-os","method":"hck-box-exec","body":{"name":"alpine"}}`

	testMessage[BoxExecRequestBody](t, message, value)
}

func TestBoxDeleteRequest(t *testing.T) {
	message := NewBoxDeleteRequest(clientOrigin, testBoxes)
	value := `{"kind":"api/v1","origin":"hckctl-0.0.0-os","method":"hck-box-delete","body":{"names":["box-alpine-123","box-alpine-456"]}}`

	testMessage[BoxDeleteRequestBody](t, message, value)
}

func TestBoxDeleteResponse(t *testing.T) {
	items := []BoxDeleteItem{
		{Id: "123", Name: testBoxes[0]},
		{Id: "456", Name: testBoxes[1]},
	}
	message := NewBoxDeleteResponse(serverOrigin, items)
	value := `{"kind":"api/v1","origin":"hckadm-0.0.0-info","method":"hck-box-delete","body":{"items":[{"Id":"123","Name":"box-alpine-123"},{"Id":"456","Name":"box-alpine-456"}]}}`

	testMessage[BoxDeleteResponseBody](t, message, value)
}

func TestBoxListRequest(t *testing.T) {
	request := NewBoxListRequest(clientOrigin)
	value := `{"kind":"api/v1","origin":"hckctl-0.0.0-os","method":"hck-box-list","body":{}}`

	testMessage[BoxListRequestBody](t, request, value)
}

func TestBoxListResponse(t *testing.T) {
	items := []BoxListItem{
		{Id: "123", Name: testBoxes[0], Healthy: true},
		{Id: "456", Name: testBoxes[1], Healthy: false},
	}
	request := NewBoxListResponse(serverOrigin, items)
	value := `{"kind":"api/v1","origin":"hckadm-0.0.0-info","method":"hck-box-list","body":{"items":[{"Id":"123","Name":"box-alpine-123","Healthy":true},{"Id":"456","Name":"box-alpine-456","Healthy":false}]}}`

	testMessage[BoxListResponseBody](t, request, value)
}

func testMessage[T body](t *testing.T, message *Message[T], value string) {
	jsonString, err := message.Encode()
	assert.NoError(t, err)
	assert.Equal(t, value, jsonString)

	result, err := Decode[T](value)
	assert.NoError(t, err)
	assert.Equal(t, message, result)

	protocol := fmt.Sprintf("%s/%s", message.Kind, message.Method)
	assert.Equal(t, protocol, message.Protocol())
}
