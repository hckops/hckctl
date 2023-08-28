package model

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExpandBox(t *testing.T) {
	testLab := &LabV1{
		Kind: "lab/v1",
		Name: "my-name",
		Tags: []string{"my-tag"},
		Box: LabBox{
			Alias:    "${alias:myName}",
			Template: "myTemplate",
			Env:      []string{"MY_KEY=MY_VALUE"},
			Size:     "xs",
			Vpn:      "${vpn:default}",
			Ports:    []string{"port-1", "port-2", "port-3"},
			Dumps:    []string{"dump-1", "dump-2", "dump-3"},
		},
	}
	input := map[string]string{
		"alias": "myAlias",
		"vpn":   "myVpn",
	}
	expected := &LabV1{
		Kind: "lab/v1",
		Name: "my-name",
		Tags: []string{"my-tag"},
		Box: LabBox{
			Alias:    "myAlias",
			Template: "myTemplate",
			Env:      []string{"MY_KEY=MY_VALUE"},
			Size:     "xs",
			Vpn:      "myVpn",
			Ports:    []string{"port-1", "port-2", "port-3"},
			Dumps:    []string{"dump-1", "dump-2", "dump-3"},
		},
	}
	result, err := testLab.ExpandBox(input)

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestExpandEnv(t *testing.T) {
	testLabEnv := &LabV1{
		Box: LabBox{
			Env: []string{
				"MY_KEY=${value:myValue}",
				"PASSWORD=${password:random}",
			},
		},
	}
	result, err := testLabEnv.ExpandBox(map[string]string{})

	assert.NoError(t, err)
	assert.Equal(t, result.Box.Env[0], "MY_KEY=myValue")
	assert.True(t, strings.HasPrefix(result.Box.Env[1], "PASSWORD="))
	assert.Len(t, result.Box.Env[1], 19)
}

func TestExpandAliasEmpty(t *testing.T) {
	testLabAlias := &LabV1{
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
	testLabAlias := &LabV1{
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
	testLabAlias := &LabV1{
		Box: LabBox{
			Alias: "$alias",
		},
	}
	_, err := testLabAlias.ExpandBox(map[string]string{})

	assert.EqualError(t, err, "alias required")
}

func TestExpandAliasRequiredTemplateFormat(t *testing.T) {
	testLabAlias := &LabV1{
		Box: LabBox{
			Alias: "${alias}",
		},
	}
	_, err := testLabAlias.ExpandBox(map[string]string{})

	assert.EqualError(t, err, "alias required")
}

func TestExpandAliasInput(t *testing.T) {
	testLabAlias := &LabV1{
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
	testLabAlias := &LabV1{
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
	testLabAlias := &LabV1{
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
