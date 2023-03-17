package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewCommandOpenBox(t *testing.T) {
	assert.Equal(t, "hck-box-open::my-group/my-name:main", NewCommandOpenBox("my-group/my-name", "main"))
}
