package persistentvolumeclaim

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefinePersistentVolumeClaim(t *testing.T) {
	testPVC := DefinePersistentVolumeClaim("testPVC", "testNamespace")
	assert.NotNil(t, testPVC)
	assert.Equal(t, "testPVC", testPVC.Name)
	assert.Equal(t, "testNamespace", testPVC.Namespace)
	assert.Equal(t, "3Gi", testPVC.Spec.Resources.Requests.Storage().String())
}

func TestRedefineWithStorageClass(t *testing.T) {
	testPVC := DefinePersistentVolumeClaim("testPVC", "testNamespace")
	assert.NotNil(t, testPVC)
	assert.Equal(t, "testPVC", testPVC.Name)
	assert.Nil(t, testPVC.Spec.StorageClassName)
	RedefineWithStorageClass(testPVC, "testStorageClass")
	assert.Equal(t, "testStorageClass", *testPVC.Spec.StorageClassName)
}
