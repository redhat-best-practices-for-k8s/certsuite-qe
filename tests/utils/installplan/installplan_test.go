package installplan

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefineInstallPlan(t *testing.T) {
	ip := DefineInstallPlan("test", "default")
	assert.NotNil(t, ip)
	assert.Equal(t, "test", ip.Name)
	assert.Equal(t, "default", ip.Namespace)
}
