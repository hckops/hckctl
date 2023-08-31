package model

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/hckops/hckctl/pkg/box/model"
)

func TestLocalLabels(t *testing.T) {
	labels := NewBoxLabels().AddDefaultLocal()
	expected := Labels{
		"com.hckops.schema.kind":    "box/v1",
		"com.hckops.template.local": "true",
	}

	assert.Equal(t, 2, len(labels))
	assert.Equal(t, expected, labels)
}

func TestGitLabels(t *testing.T) {
	labels := NewBoxLabels().AddDefaultGit("myUrl", "myRevision", "myDir")
	expected := Labels{
		"com.hckops.schema.kind":           "box/v1",
		"com.hckops.template.git":          "true",
		"com.hckops.template.git.url":      "myUrl",
		"com.hckops.template.git.revision": "myRevision",
		"com.hckops.template.git.dir":      "myDir",
	}

	assert.Equal(t, 5, len(labels))
	assert.Equal(t, expected, labels)
}

func TestAboxModelddLocalLabels(t *testing.T) {
	labels := NewBoxLabels().AddBoxSize(model.Small).AddLocal("/tmp/cache")
	expected := Labels{
		"com.hckops.schema.kind":         "box/v1",
		"com.hckops.template.local":      "true",
		"com.hckops.template.cache.path": "/tmp/cache",
		"com.hckops.box.size":            "s",
	}

	assert.Equal(t, 4, len(labels))
	assert.Equal(t, expected, labels)
}

func TestAddLocalLabelsInvalid(t *testing.T) {
	initial := NewBoxLabels().AddDefaultLocal()
	labels := initial.AddGit("myPath", "myCommit")

	assert.Equal(t, len(initial), len(labels))
}

func TestAddGitLabels(t *testing.T) {
	gitLabels := NewBoxLabels().AddDefaultGit("https://github.com/hckops/megalopolis", "main", "megalopolis")

	path := "/home/test/.cache/hck/megalopolis/box/base/arch.yml"
	labels := gitLabels.AddBoxSize(model.Medium).AddGit(path, "myCommit")
	expected := Labels{
		"com.hckops.schema.kind":           "box/v1",
		"com.hckops.template.git":          "true",
		"com.hckops.template.git.url":      "https://github.com/hckops/megalopolis",
		"com.hckops.template.git.revision": "main",
		"com.hckops.template.git.dir":      "megalopolis",
		"com.hckops.template.git.commit":   "myCommit",
		"com.hckops.template.git.name":     "box/base/arch",
		"com.hckops.template.cache.path":   path,
		"com.hckops.box.size":              "m",
	}

	assert.Equal(t, 9, len(labels))
	assert.Equal(t, expected, labels)
}

func TestAddGitLabelsInvalid(t *testing.T) {
	initial := NewBoxLabels().AddDefaultGit("https://github.com/hckops/megalopolis", "main", "megalopolis")
	labels := initial.AddLocal("/tmp/cache")

	assert.Equal(t, len(initial), len(labels))
}

func TestExist(t *testing.T) {
	labels := Labels{
		"my.name": "myValue",
	}
	value, err := labels.exist("my.name")

	assert.NoError(t, err)
	assert.Equal(t, "myValue", value)
}

func TestExistError(t *testing.T) {
	_, err := Labels{}.exist("my.name")
	assert.EqualError(t, err, "label my.name not found")
}

func TestToSize(t *testing.T) {
	labels := Labels{
		"com.hckops.box.size": "xl",
	}
	size, err := labels.ToBoxSize()

	assert.NoError(t, err)
	assert.Equal(t, model.ExtraLarge, size)
}

func TestToSizeError(t *testing.T) {
	_, errLabel := Labels{}.ToBoxSize()
	assert.EqualError(t, errLabel, "label com.hckops.box.size not found")

	_, errSize := Labels{"com.hckops.box.size": "abc"}.ToBoxSize()
	assert.EqualError(t, errSize, "invalid resource size value=abc")
}

func TestToCachedTemplateInfo(t *testing.T) {
	info := NewBoxLabels().AddDefaultLocal().
		AddBoxSize(model.Small).
		AddLocal("/tmp/cache").
		ToCachedTemplateInfo()

	expected := &CachedTemplateInfo{
		Path: "/tmp/cache",
	}

	assert.Equal(t, expected, info)
}

func TestToBoxTemplateInfo(t *testing.T) {
	info := NewBoxLabels().AddDefaultGit("myUrl", "myRevision", "myDir").
		AddBoxSize(model.Medium).
		AddGit("myDir/myName", "myCommit").
		ToGitTemplateInfo()

	expected := &GitTemplateInfo{
		Url:      "myUrl",
		Revision: "myRevision",
		Commit:   "myCommit",
		Name:     "myName",
	}

	assert.Equal(t, expected, info)
}
