package replicaset

import (
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
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

func TestRedefineWithCustomTolerations(t *testing.T) {
	testRS := DefineReplicaSet("testRS", "testNamespace", "testImage", map[string]string{"app": "test"})
	assert.NotNil(t, testRS)

	// Replace whatever DefineReplicaSet auto-added (depends on ENABLE_INFRASTRUCTURE_TOLERATIONS)
	// with one known toleration, so the assertion is deterministic and pins that
	// RedefineWithCustomTolerations appends rather than replaces.
	existingToleration := corev1.Toleration{Key: "existing-taint", Operator: corev1.TolerationOpExists}
	testRS.Spec.Template.Spec.Tolerations = []corev1.Toleration{existingToleration}

	customToleration := corev1.Toleration{
		Key:      "custom-taint",
		Operator: corev1.TolerationOpEqual,
		Value:    "custom-value",
		Effect:   corev1.TaintEffectNoExecute,
	}

	RedefineWithCustomTolerations(testRS, []corev1.Toleration{customToleration})
	assert.Equal(t, []corev1.Toleration{existingToleration, customToleration}, testRS.Spec.Template.Spec.Tolerations)
}
