package replicaset

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefineReplicaSet(t *testing.T) {
	testRS := DefineReplicaSet("testRS", "testNamespace", "testImage", map[string]string{"app": "test"})
	assert.NotNil(t, testRS)
	assert.Equal(t, "testRS", testRS.Name)
	assert.Equal(t, "testNamespace", testRS.Namespace)
	assert.Equal(t, "testImage", testRS.Spec.Template.Spec.Containers[0].Image)
	assert.Equal(t, "test", testRS.Spec.Template.Labels["app"])
}

func TestRedefineWithReplicaNumber(t *testing.T) {
	testRS := DefineReplicaSet("testRS", "testNamespace", "testImage", map[string]string{"app": "test"})
	assert.NotNil(t, testRS)
	assert.Equal(t, "testRS", testRS.Name)
	assert.Equal(t, int32(1), *testRS.Spec.Replicas)
	RedefineWithReplicaNumber(testRS, 2)
	assert.Equal(t, int32(2), *testRS.Spec.Replicas)
}

func TestRedefineWithPVC(t *testing.T) {
	testRS := DefineReplicaSet("testRS", "testNamespace", "testImage", map[string]string{"app": "test"})
	assert.NotNil(t, testRS)
	assert.Equal(t, "testRS", testRS.Name)
	RedefineWithPVC(testRS, "testVolume", "testPVCName")
	assert.Equal(t, "testVolume", testRS.Spec.Template.Spec.Volumes[0].Name)
	assert.Equal(t, "testPVCName", testRS.Spec.Template.Spec.Volumes[0].PersistentVolumeClaim.ClaimName)
}
