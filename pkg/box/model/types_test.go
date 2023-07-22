package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsCached(t *testing.T) {
	info := &BoxTemplateInfo{
		CachedTemplate: &CachedTemplateInfo{
			Path: "/tmp/cache",
		},
	}
	assert.True(t, info.IsCached())

	assert.False(t, (&BoxTemplateInfo{}).IsCached())
}
