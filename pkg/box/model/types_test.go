package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToBoxNames(t *testing.T) {
	boxes := []BoxInfo{
		{Id: "id1", Name: "name1", Healthy: true},
		{Id: "id2", Name: "name2", Healthy: false},
		{Id: "id3", Name: "name3", Healthy: true},
	}
	expected := []string{"name1", "name2", "name3"}

	assert.Equal(t, 3, len(boxes))
	assert.Equal(t, expected, ToBoxNames(boxes))
}
