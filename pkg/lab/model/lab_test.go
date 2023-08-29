package model

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	boxModel "github.com/hckops/hckctl/pkg/box/model"
)

func TestMerge(t *testing.T) {
	testBox := &boxModel.BoxV1{
		Env: []string{
			"TTYD_USERNAME=username",
			"TTYD_PASSWORD=password",
		},
	}
	testLab := &LabV1{
		Box: LabBox{
			Template: BoxTemplate{
				Env: []string{
					"TTYD_PASSWORD=${password:random}",
				},
			},
		},
	}
	expected := &boxModel.BoxV1{
		Env: []string{
			"TTYD_USERNAME=username",
			"TTYD_PASSWORD=${password:random}", // merge
		},
	}
	result := testLab.Box.Template.Merge(testBox)

	assert.Equal(t, expected, result)
}

func TestExpand(t *testing.T) {
	testLabBox := &LabBox{
		Alias: "${alias:myName}", // expand
		Template: BoxTemplate{
			Name: "myTemplate",
			Env:  []string{"MY_KEY=MY_VALUE"},
		},
		Size:  "xs",
		Vpn:   "${vpn:default}", // expand
		Ports: []string{"port-1", "port-2", "port-3"},
		Dumps: []string{"dump-1", "dump-2", "dump-3"},
	}
	input := map[string]string{
		"alias": "myAlias",
		"vpn":   "myVpn",
	}
	expected := &LabBox{
		Alias: "myAlias",
		Template: BoxTemplate{
			Name: "myTemplate",
			Env:  []string{"MY_KEY=MY_VALUE"},
		},
		Size:  "xs",
		Vpn:   "myVpn",
		Ports: []string{"port-1", "port-2", "port-3"},
		Dumps: []string{"dump-1", "dump-2", "dump-3"},
	}
	result, err := testLabBox.Expand(input)

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestExpandEnv(t *testing.T) {
	testLabEnv := &LabV1{
		Box: LabBox{
			Template: BoxTemplate{
				Env: []string{
					"MY_KEY=${value:myValue}",
					"PASSWORD=${password:random}",
				},
			},
		},
	}
	result, err := testLabEnv.Box.Expand(map[string]string{})

	assert.NoError(t, err)
	assert.Equal(t, result.Template.Env[0], "MY_KEY=myValue")
	assert.True(t, strings.HasPrefix(result.Template.Env[1], "PASSWORD="))
	assert.Len(t, result.Template.Env[1], 19)
}

func TestExpandAliasEmpty(t *testing.T) {
	testLabAlias := &LabBox{
		Alias: "${ \n\t\r  }",
	}
	expected := &LabBox{
		Alias: "",
	}
	result, err := testLabAlias.Expand(map[string]string{})

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestExpandAliasNone(t *testing.T) {
	testLabAlias := &LabBox{
		Alias: "{alias:myName}",
	}
	expected := &LabBox{
		Alias: "{alias:myName}",
	}
	result, err := testLabAlias.Expand(map[string]string{})

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestExpandAliasRequiredSimpleFormat(t *testing.T) {
	testLabAlias := &LabV1{
		Box: LabBox{
			Alias: "$alias",
		},
	}
	_, err := testLabAlias.Box.Expand(map[string]string{})

	assert.EqualError(t, err, "alias required")
}

func TestExpandAliasRequiredTemplateFormat(t *testing.T) {
	testLabAlias := &LabV1{
		Box: LabBox{
			Alias: "${alias}",
		},
	}
	_, err := testLabAlias.Box.Expand(map[string]string{})

	assert.EqualError(t, err, "alias required")
}

func TestExpandAliasInput(t *testing.T) {
	testLabAlias := &LabBox{
		Alias: "${alias:myName}",
	}
	input := map[string]string{
		"alias": "myAlias",
	}
	expected := &LabBox{
		Alias: "myAlias",
	}
	result, err := testLabAlias.Expand(input)

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestExpandAliasDefault(t *testing.T) {
	testLabAlias := &LabBox{
		Alias: "${alias:myName}",
	}
	expected := &LabBox{
		Alias: "myName",
	}
	result, err := testLabAlias.Expand(map[string]string{})

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestExpandAliasRandom(t *testing.T) {
	testLabAlias := &LabBox{
		Alias: "${alias:random}",
	}
	result, err := testLabAlias.Expand(map[string]string{})

	assert.NoError(t, err)
	assert.Len(t, result.Alias, 10)
}

func TestExpandAliasMissing(t *testing.T) {
	expected := &LabBox{
		Alias: "",
	}
	result, err := (&LabBox{}).Expand(map[string]string{})

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}
