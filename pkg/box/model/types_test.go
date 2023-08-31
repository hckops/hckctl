package model

import (
	"testing"

	"github.com/stretchr/testify/assert"

	commonModel "github.com/hckops/hckctl/pkg/common/model"
)

func TestIsCached(t *testing.T) {
	info := &BoxTemplateInfo{
		CachedTemplate: &commonModel.CachedTemplateInfo{
			Path: "/tmp/cache",
		},
	}
	assert.True(t, info.IsCached())

	assert.False(t, (&BoxTemplateInfo{}).IsCached())
}
