package persistentvolume

import (
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
)

func TestDefinePersistentVolume(t *testing.T) {
	testPV := DefinePersistentVolume("testPV", "testPVC", "testNamespace")
	assert.NotNil(t, testPV)
	assert.Equal(t, "testPV", testPV.Name)
	assert.Equal(t, "testPVC", testPV.Spec.ClaimRef.Name)
	assert.Equal(t, "testNamespace", testPV.Spec.ClaimRef.Namespace)
	assert.Equal(t, "/tmp", testPV.Spec.Local.Path)
	assert.Equal(t, "10Gi", testPV.Spec.Capacity.Storage().String())
}

func TestRedefineWithPVReclaimPolicy(t *testing.T) {
	testPV := DefinePersistentVolume("testPV", "testPVC", "testNamespace")
	assert.NotNil(t, testPV)
	assert.Equal(t, "testPV", testPV.Name)
	assert.Equal(t, corev1.PersistentVolumeReclaimRetain, testPV.Spec.PersistentVolumeReclaimPolicy)
	RedefineWithPVReclaimPolicy(testPV, corev1.PersistentVolumeReclaimDelete)
	assert.Equal(t, corev1.PersistentVolumeReclaimDelete, testPV.Spec.PersistentVolumeReclaimPolicy)
}

func TestRedefineWithStorageClass(t *testing.T) {
	testPV := DefinePersistentVolume("testPV", "testPVC", "testNamespace")
	assert.NotNil(t, testPV)
	assert.Equal(t, "testPV", testPV.Name)
	assert.Equal(t, "", testPV.Spec.StorageClassName)
	RedefineWithStorageClass(testPV, "testStorageClass")
	assert.Equal(t, "testStorageClass", testPV.Spec.StorageClassName)
}
