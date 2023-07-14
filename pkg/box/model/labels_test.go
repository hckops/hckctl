package model

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLocalLabels(t *testing.T) {
	labels := NewLocalLabels()
	expected := map[string]string{
		"com.hckops.schema.kind":    "box/v1",
		"com.hckops.template.local": "true",
	}

	assert.Equal(t, 2, len(expected))
	assert.Equal(t, expected, labels)
}

func TestGitLabels(t *testing.T) {
	labels := NewGitLabels("myName", "myUrl", "myRevision")
	expected := map[string]string{
		"com.hckops.schema.kind":           "box/v1",
		"com.hckops.template.git":          "true",
		"com.hckops.template.git.name":     "myName",
		"com.hckops.template.git.url":      "myUrl",
		"com.hckops.template.git.revision": "myRevision",
	}

	assert.Equal(t, 5, len(expected))
	assert.Equal(t, expected, labels)
}
