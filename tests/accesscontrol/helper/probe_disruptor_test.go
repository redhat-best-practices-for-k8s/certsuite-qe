package helper

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProbeMemoryPressureCommand(t *testing.T) {
	t.Parallel()

	cmd := probeMemoryPressureCommand()
	assert.Len(t, cmd, 3)
	assert.Equal(t, "/bin/sh", cmd[0])
	assert.Equal(t, "-c", cmd[1])
	assert.Contains(t, cmd[2], "kill -9 1")
	assert.Contains(t, cmd[2], "1024*1024")
}
