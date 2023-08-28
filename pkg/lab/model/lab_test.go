package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExpandAliasEmpty(t *testing.T) {
	var testLabAlias = &LabV1{
		Box: LabBox{
			Alias: "${ \n\t\r  }",
		},
	}
	expected := &LabV1{
		Box: LabBox{
			Alias: "",
		},
	}
	result, err := testLabAlias.ExpandBox(map[string]string{})

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestExpandAliasNone(t *testing.T) {
	var testLabAlias = &LabV1{
		Box: LabBox{
			Alias: "{alias:myName}",
		},
	}
	expected := &LabV1{
		Box: LabBox{
			Alias: "{alias:myName}",
		},
	}
	result, err := testLabAlias.ExpandBox(map[string]string{})

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestExpandAliasRequiredSimpleFormat(t *testing.T) {
	var testLabAlias = &LabV1{
		Box: LabBox{
			Alias: "$alias",
		},
	}
	_, err := testLabAlias.ExpandBox(map[string]string{})

	assert.EqualError(t, err, "alias required")
}

func TestExpandAliasRequiredTemplateFormat(t *testing.T) {
	var testLabAlias = &LabV1{
		Box: LabBox{
			Alias: "${alias}",
		},
	}
	_, err := testLabAlias.ExpandBox(map[string]string{})

	assert.EqualError(t, err, "alias required")
}

func TestExpandAliasInput(t *testing.T) {
	var testLabAlias = &LabV1{
		Box: LabBox{
			Alias: "${alias:myName}",
		},
	}
	input := map[string]string{
		"alias": "myAlias",
	}
	expected := &LabV1{
		Box: LabBox{
			Alias: "myAlias",
		},
	}
	result, err := testLabAlias.ExpandBox(input)

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestExpandAliasDefault(t *testing.T) {
	var testLabAlias = &LabV1{
		Box: LabBox{
			Alias: "${alias:myName}",
		},
	}
	expected := &LabV1{
		Box: LabBox{
			Alias: "myName",
		},
	}
	result, err := testLabAlias.ExpandBox(map[string]string{})

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestExpandAliasRandom(t *testing.T) {
	var testLabAlias = &LabV1{
		Box: LabBox{
			Alias: "${alias:random}",
		},
	}
	result, err := testLabAlias.ExpandBox(map[string]string{})

	assert.NoError(t, err)
	assert.Len(t, result.Box.Alias, 10)
}

func TestExpandAliasMissing(t *testing.T) {
	expected := &LabV1{
		Box: LabBox{
			Alias: "",
		},
	}
	result, err := (&LabV1{Box: LabBox{}}).ExpandBox(map[string]string{})

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}
