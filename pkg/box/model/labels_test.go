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
	labels := NewGitLabels("myUrl", "myRevision", "myDir")
	expected := BoxLabels{
		"com.hckops.schema.kind":           "box/v1",
		"com.hckops.template.git":          "true",
		"com.hckops.template.git.url":      "myUrl",
		"com.hckops.template.git.revision": "myRevision",
		"com.hckops.template.git.dir":      "myDir",
	}

	assert.Equal(t, 5, len(expected))
	assert.Equal(t, expected, labels)
}

func TestAddLabels(t *testing.T) {
	defaultLabels := BoxLabels{
		"com.hckops.schema.kind": "box/v1",
		"com.hckops.test":        "true",
	}
	labels := defaultLabels.AddLabels("myPath", "skipped", ExtraLarge)
	expected := BoxLabels{
		"com.hckops.schema.kind":          "box/v1",
		"com.hckops.test":                 "true",
		"com.hckops.template.common.path": "myPath",
		"com.hckops.box.size":             "xl",
	}

	assert.Equal(t, 4, len(expected))
	assert.Equal(t, expected, labels)
}

func TestAddGitLabels(t *testing.T) {
	gitLabels := NewGitLabels("https://github.com/hckops/megalopolis", "main", "megalopolis")

	path := "/home/test/.cache/hck/megalopolis/box/base/arch.yml"
	labels := gitLabels.AddLabels(path, "myCommit", Medium)
	expected := BoxLabels{
		"com.hckops.schema.kind":           "box/v1",
		"com.hckops.template.git":          "true",
		"com.hckops.template.git.url":      "https://github.com/hckops/megalopolis",
		"com.hckops.template.git.revision": "main",
		"com.hckops.template.git.dir":      "megalopolis",
		"com.hckops.template.git.commit":   "myCommit",
		"com.hckops.template.git.name":     "box/base/arch",
		"com.hckops.template.common.path":  path,
		"com.hckops.box.size":              "m",
	}

	assert.Equal(t, 9, len(expected))
	assert.Equal(t, expected, labels)
}

func TestExist(t *testing.T) {
	labels := BoxLabels{
		"my.name": "myValue",
	}
	value, err := labels.exist("my.name")

	assert.NoError(t, err)
	assert.Equal(t, "myValue", value)
}

func TestExistError(t *testing.T) {
	_, err := BoxLabels{}.exist("my.name")
	assert.EqualError(t, err, "label my.name not found")
}

func TestToSize(t *testing.T) {
	labels := BoxLabels{
		"com.hckops.box.size": "xl",
	}
	size, err := labels.ToSize()

	assert.NoError(t, err)
	assert.Equal(t, ExtraLarge, size)
}

func TestToSizeError(t *testing.T) {
	_, errLabel := BoxLabels{}.ToSize()
	assert.EqualError(t, errLabel, "label com.hckops.box.size not found")

	_, errSize := BoxLabels{"com.hckops.box.size": "invalid"}.ToSize()
	assert.EqualError(t, errSize, "invalid resource size")
}

func TestToBoxTemplateInfo(t *testing.T) {
	labels := BoxLabels{
		"com.hckops.template.git":          "true",
		"com.hckops.template.git.url":      "myUrl",
		"com.hckops.template.git.revision": "myRevision",
		"com.hckops.template.git.commit":   "myCommit",
		"com.hckops.template.git.name":     "myName",
	}
	expected := &BoxTemplateInfo{
		Url:      "myUrl",
		Revision: "myRevision",
		Commit:   "myCommit",
		Name:     "myName",
	}
	result, err := labels.ToBoxTemplateInfo()

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestToBoxTemplateInfoError(t *testing.T) {
	_, err := BoxLabels{}.ToBoxTemplateInfo()
	assert.EqualError(t, err, "label com.hckops.template.git not found")
}
