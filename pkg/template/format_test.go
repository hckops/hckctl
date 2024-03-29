package template

import (
	"testing"

	"github.com/MakeNowJust/heredoc"
	"github.com/stretchr/testify/assert"

	"github.com/hckops/hckctl/pkg/box/model"
	"github.com/hckops/hckctl/pkg/schema"
)

var exampleYaml = heredoc.Doc(`
	kind: box/v1
	name: my-box
	# examples
	tags:
	  - 'test'
	  - "official"
	image:
	  repository: "hckops/my-image"
	  # sha or tag
	  version: latest
	env:
	  - TTYD_USERNAME=username
	  - TTYD_PASSWORD=password
	network:
	  # name:local[:remote]
	  ports:
	    - aaa:123
	    - 'bbb:456:789'
`)

func TestConvertFromYamlToYamlEmpty(t *testing.T) {
	expected := heredoc.Doc(`
		kind: ""
		name: ""
		tags: []
		image:
		  repository: ""
		  version: ""
		shell: ""
		env: []
		network:
		  ports: []
	`)

	result, err := convertFromYamlToYaml(schema.KindBoxV1, "")
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestConvertFromYamlToYaml(t *testing.T) {
	expected := heredoc.Doc(`
		kind: box/v1
		name: my-box
		tags:
		- test
		- official
		image:
		  repository: hckops/my-image
		  version: latest
		shell: ""
		env:
		- TTYD_USERNAME=username
		- TTYD_PASSWORD=password
		network:
		  ports:
		  - aaa:123
		  - bbb:456:789
	`)

	result, err := convertFromYamlToYaml(schema.KindBoxV1, exampleYaml)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestConvertFromYamlToJsonEmpty(t *testing.T) {
	expected := heredoc.Doc(`{
      "Kind": "",
      "Name": "",
      "Tags": null,
      "Image": {
        "Repository": "",
        "Version": ""
      },
      "Shell": "",
      "Env": null,
      "Network": {
        "Ports": null
      }
    }`)

	result, err := convertFromYamlToJson(schema.KindBoxV1, "")
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestConvertFromYamlToJson(t *testing.T) {
	expected := heredoc.Doc(`{
      "Kind": "box/v1",
      "Name": "my-box",
      "Tags": [
        "test",
        "official"
      ],
      "Image": {
        "Repository": "hckops/my-image",
        "Version": "latest"
      },
      "Shell": "",
      "Env": [
        "TTYD_USERNAME=username",
        "TTYD_PASSWORD=password"
      ],
      "Network": {
        "Ports": [
          "aaa:123",
          "bbb:456:789"
        ]
      }
    }`)

	result, err := convertFromYamlToJson(schema.KindBoxV1, exampleYaml)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestDecodeFromYaml(t *testing.T) {
	expected := &model.BoxV1{
		Kind: "box/v1",
		Name: "my-box",
		Tags: []string{"test", "official"},
		Image: struct {
			Repository string
			Version    string
		}{
			Repository: "hckops/my-image",
			Version:    "latest",
		},
		Shell: "",
		Env: []string{
			"TTYD_USERNAME=username",
			"TTYD_PASSWORD=password",
		},
		Network: struct{ Ports []string }{Ports: []string{
			"aaa:123",
			"bbb:456:789",
		}},
	}

	result, err := decodeFromYaml[model.BoxV1](exampleYaml)
	assert.NoError(t, err)
	assert.Equal(t, expected, &result)
}
