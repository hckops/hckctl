package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLocalLabels(t *testing.T) {
	labels := NewLocalLabels()
	expected := BoxLabels{
		"com.hckops.schema.kind":    "box/v1",
		"com.hckops.template.local": "true",
	}

	assert.Equal(t, 2, len(expected))
	assert.Equal(t, expected, labels)
}

func TestGitLabels(t *testing.T) {
	labels := NewGitLabels("myName", "myUrl", "myRevision")
	expected := BoxLabels{
		"com.hckops.schema.kind":           "box/v1",
		"com.hckops.template.git":          "true",
		"com.hckops.template.git.name":     "myName",
		"com.hckops.template.git.url":      "myUrl",
		"com.hckops.template.git.revision": "myRevision",
	}

	assert.Equal(t, 5, len(expected))
	assert.Equal(t, expected, labels)
}

func TestAddLabels(t *testing.T) {
	defaultLabels := BoxLabels{
		"com.hckops.schema.kind": "box/v1",
		"com.hckops.test":        "true",
	}
	labels := defaultLabels.AddLabels("myPath", ExtraLarge)
	expected := BoxLabels{
		"com.hckops.schema.kind":          "box/v1",
		"com.hckops.test":                 "true",
		"com.hckops.template.common.path": "myPath",
		"com.hckops.box.size":             "xl",
	}

	assert.Equal(t, 4, len(expected))
	assert.Equal(t, expected, labels)
}
