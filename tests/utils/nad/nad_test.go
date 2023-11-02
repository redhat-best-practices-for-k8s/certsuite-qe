package nad

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefineNad(t *testing.T) {
	testNad := DefineNad("test", "default")
	assert.NotNil(t, testNad)
	assert.Equal(t, "test", testNad.Name)
	assert.Equal(t, "default", testNad.Namespace)
	assert.Equal(t, `{"cniVersion": "0.4.0", "name": "test", "type": "macvlan", "mode": "bridge"}`, testNad.Spec.Config)
}

func TestRedefineNadWithWhereaboutsIpam(t *testing.T) {
	testNad := DefineNad("test", "default")
	RedefineNadWithWhereaboutsIpam(testNad, "testnetwork")
	//nolint:lll
	assert.Equal(t, `{"cniVersion": "0.4.0", "name": "test", "type": "macvlan", "mode": "bridge", "ipam":{ "type": "whereabouts", "range": "testnetwork"}}`,
		testNad.Spec.Config)
}
