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
	assert.Equal(t, 7, len(methods))
	assert.Equal(t, "hck-ping", methods[MethodPing])
	assert.Equal(t, "hck-box-create", methods[MethodBoxCreate])
	assert.Equal(t, "hck-box-delete", methods[MethodBoxDelete])
	assert.Equal(t, "hck-box-describe", methods[MethodBoxDescribe])
	assert.Equal(t, "hck-box-exec", methods[MethodBoxExec])
	assert.Equal(t, "hck-box-list", methods[MethodBoxList])
	assert.Equal(t, "hck-lab-create", methods[MethodLabCreate])
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

func TestBoxDeleteRequest(t *testing.T) {
	message := NewBoxDeleteRequest(clientOrigin, testBoxes)
	value := `{"kind":"api/v1","origin":"hckctl-0.0.0-os","method":"hck-box-delete","body":{"names":["box-alpine-123","box-alpine-456"]}}`

	testMessage[BoxDeleteRequestBody](t, message, value)
}

func TestBoxDeleteResponse(t *testing.T) {
	message := NewBoxDeleteResponse(serverOrigin, testBoxes)
	value := `{"kind":"api/v1","origin":"hckadm-0.0.0-info","method":"hck-box-delete","body":{"names":["box-alpine-123","box-alpine-456"]}}`

	testMessage[BoxDeleteResponseBody](t, message, value)
}

func TestBoxDescribeRequest(t *testing.T) {
	message := NewBoxDescribeRequest(clientOrigin, testBoxes[0])
	value := `{"kind":"api/v1","origin":"hckctl-0.0.0-os","method":"hck-box-describe","body":{"name":"box-alpine-123"}}`

	testMessage[BoxDescribeRequestBody](t, message, value)
}

func TestBoxDescribeResponse(t *testing.T) {
	message := NewBoxDescribeResponse(serverOrigin, BoxDescribeResponseBody{
		Id:      "myId",
		Name:    "myName",
		Created: "myCreated",
		Healthy: true,
		Size:    "mySize",
		Template: &BoxDescribeTemplateInfo{
			Public:   true,
			Url:      "infoUrl",
			Revision: "infoRevision",
			Commit:   "infoCommit",
			Name:     "infoName",
		},
		Env:   []string{"KEY_1=VALUE_1", "KEY_2=VALUE_2"},
		Ports: []string{"alias-1/123", "alias-2/456"},
	})
	value := `{"kind":"api/v1","origin":"hckadm-0.0.0-info","method":"hck-box-describe","body":{"id":"myId","name":"myName","created":"myCreated","healthy":true,"size":"mySize","template":{"public":true,"url":"infoUrl","revision":"infoRevision","commit":"infoCommit","name":"infoName"},"env":["KEY_1=VALUE_1","KEY_2=VALUE_2"],"ports":["alias-1/123","alias-2/456"]}}`

	testMessage[BoxDescribeResponseBody](t, message, value)
}

func TestBoxListRequest(t *testing.T) {
	message := NewBoxListRequest(clientOrigin)
	value := `{"kind":"api/v1","origin":"hckctl-0.0.0-os","method":"hck-box-list","body":{}}`

	testMessage[BoxListRequestBody](t, message, value)
}

func TestBoxListResponse(t *testing.T) {
	items := []BoxListItem{
		{Id: "123", Name: testBoxes[0], Healthy: true},
		{Id: "456", Name: testBoxes[1], Healthy: false},
	}
	message := NewBoxListResponse(serverOrigin, items)
	value := `{"kind":"api/v1","origin":"hckadm-0.0.0-info","method":"hck-box-list","body":{"items":[{"Id":"123","Name":"box-alpine-123","Healthy":true},{"Id":"456","Name":"box-alpine-456","Healthy":false}]}}`

	testMessage[BoxListResponseBody](t, message, value)
}

func TestBoxExecSession(t *testing.T) {
	message := NewBoxExecSession(clientOrigin, "alpine")
	value := `{"kind":"api/v1","origin":"hckctl-0.0.0-os","method":"hck-box-exec","body":{"name":"alpine"}}`

	testMessage[BoxExecSessionBody](t, message, value)
}

func TestLabCreateRequest(t *testing.T) {
	parameters := map[string]string{
		"password":        "changeme",
		"vpn.default.ref": "OPENVPN_REMOTE_REF",
	}
	message := NewLabCreateRequest(clientOrigin, "ctf-vpn", parameters)
	value := `{"kind":"api/v1","origin":"hckctl-0.0.0-os","method":"hck-lab-create","body":{"templateName":"ctf-vpn","params":{"password":"changeme","vpn.default.ref":"OPENVPN_REMOTE_REF"}}}`

	testMessage[LabCreateRequestBody](t, message, value)
}

func TestLabCreateResponse(t *testing.T) {
	message := NewLabCreateResponse(serverOrigin, "lab-ctf-vpn-123")
	value := `{"kind":"api/v1","origin":"hckadm-0.0.0-info","method":"hck-lab-create","body":{"name":"lab-ctf-vpn-123"}}`

	testMessage[LabCreateResponseBody](t, message, value)
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
