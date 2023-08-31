package model

import (
	"testing"

	"github.com/stretchr/testify/assert"

	commonModel "github.com/hckops/hckctl/pkg/common/model"
)

func TestLocalLabels(t *testing.T) {
	labels := NewBoxLabels().AddDefaultLocal()
	expected := commonModel.Labels{
		"com.hckops.schema.kind":    "box/v1",
		"com.hckops.template.local": "true",
	}

	assert.Equal(t, 2, len(labels))
	assert.Equal(t, expected, labels)
}

func TestGitLabels(t *testing.T) {
	labels := NewBoxLabels().AddDefaultGit("myUrl", "myRevision", "myDir")
	expected := commonModel.Labels{
		"com.hckops.schema.kind":           "box/v1",
		"com.hckops.template.git":          "true",
		"com.hckops.template.git.url":      "myUrl",
		"com.hckops.template.git.revision": "myRevision",
		"com.hckops.template.git.dir":      "myDir",
	}

	assert.Equal(t, 5, len(labels))
	assert.Equal(t, expected, labels)
}

func TestAddLocalLabels(t *testing.T) {
	labels := AddBoxSize(NewBoxLabels().AddDefaultLocal(), Small).AddLocal("/tmp/cache")
	expected := commonModel.Labels{
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
	labels := AddBoxSize(gitLabels, Medium).AddGit(path, "myCommit")
	expected := commonModel.Labels{
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
	labels := commonModel.Labels{
		"my.name": "myValue",
	}
	value, err := labels.Exist("my.name")

	assert.NoError(t, err)
	assert.Equal(t, "myValue", value)
}

func TestExistError(t *testing.T) {
	_, err := commonModel.Labels{}.Exist("my.name")
	assert.EqualError(t, err, "label my.name not found")
}

func TestToSize(t *testing.T) {
	labels := commonModel.Labels{
		"com.hckops.box.size": "xl",
	}
	size, err := ToBoxSize(labels)

	assert.NoError(t, err)
	assert.Equal(t, ExtraLarge, size)
}

func TestToSizeError(t *testing.T) {
	_, errLabel := ToBoxSize(commonModel.Labels{})
	assert.EqualError(t, errLabel, "label com.hckops.box.size not found")

	_, errSize := ToBoxSize(commonModel.Labels{"com.hckops.box.size": "abc"})
	assert.EqualError(t, errSize, "invalid resource size value=abc")
}

func TestToCachedTemplateInfo(t *testing.T) {
	info := AddBoxSize(NewBoxLabels().AddDefaultLocal(), Small).
		AddLocal("/tmp/cache").
		ToCachedTemplateInfo()

	expected := &commonModel.CachedTemplateInfo{
		Path: "/tmp/cache",
	}

	assert.Equal(t, expected, info)
}

func TestToBoxTemplateInfo(t *testing.T) {
	info := AddBoxSize(NewBoxLabels().AddDefaultGit("myUrl", "myRevision", "myDir"), Medium).
		AddGit("myDir/myName", "myCommit").
		ToGitTemplateInfo()

	expected := &commonModel.GitTemplateInfo{
		Url:      "myUrl",
		Revision: "myRevision",
		Commit:   "myCommit",
		Name:     "myName",
	}

	assert.Equal(t, expected, info)
}
