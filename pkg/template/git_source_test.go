package template

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCacheDirName(t *testing.T) {
	opts := &GitSourceOptions{
		RepositoryUrl: "https://github.com/hckops/megalopolis.git",
	}
	assert.Equal(t, "megalopolis", opts.CacheDirName())
}

func TestCachePath(t *testing.T) {
	opts := &GitSourceOptions{
		CacheBaseDir:  "/home/test/.cache/hck",
		RepositoryUrl: "https://github.com/hckops/megalopolis.git",
	}
	assert.Equal(t, "/home/test/.cache/hck/megalopolis", opts.CachePath())
}
