package model

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateLogFileName(t *testing.T) {

	opts := &RunOptions{
		LogDir: "/tmp/demo",
	}
	logFileName := opts.GenerateLogFileName(Docker, "task-my-name")

	assert.True(t, strings.HasPrefix(logFileName, "/tmp/demo/docker-"))
	assert.True(t, strings.HasSuffix(logFileName, "-task-my-name"))
	assert.Equal(t, 49, len(logFileName))
}
