package template

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCacheDirName(t *testing.T) {
	opts := &SourceOptions{
		SourceCacheDir: "/home/test/.cache/hck/megalopolis",
	}
	assert.Equal(t, "megalopolis", opts.CacheDirName())
}
