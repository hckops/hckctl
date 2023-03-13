package template

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateTemplateParam(t *testing.T) {

	assert.NoError(t, validateTemplateParam(&TemplateParam{TemplateName: "myname"}))
}

func TestValidateTemplateParamError(t *testing.T) {

	err := validateTemplateParam(&TemplateParam{})
	if assert.Error(t, err) {
		assert.ErrorContains(t, err, "invalid name")
	}
}

func TestBuildPath(t *testing.T) {
	param := &TemplateParam{
		TemplateKind:  "box/v1",
		TemplateName:  "myname",
		Revision:      "main",
		ClientVersion: "myversion",
	}

	result, err := buildPath(param)
	assert.NoError(t, err)
	assert.Equal(t, "https://raw.githubusercontent.com/hckops/megalopolis/main/boxes/official/myname.yml", result)
}

func TestTemplateKindToPath(t *testing.T) {

	result, err := templateKindToPath("box/v1")
	assert.NoError(t, err)
	assert.Equal(t, "boxes", result)
}

func TestTemplateKindToPathError(t *testing.T) {

	result, err := templateKindToPath("")
	if assert.Error(t, err) {
		assert.ErrorContains(t, err, "invalid template kind")
	}
	assert.Empty(t, result)
}
