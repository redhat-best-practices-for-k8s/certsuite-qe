package statefulset

import (
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
)

func TestDefineStatefulSet(t *testing.T) {
	testStatefulSet := DefineStatefulSet("testStatefulSet", "testNamespace", "testImage", map[string]string{"app": "test"})
	assert.NotNil(t, testStatefulSet)
	assert.Equal(t, "testStatefulSet", testStatefulSet.Name)
	assert.Equal(t, "testNamespace", testStatefulSet.Namespace)
	assert.Equal(t, int32(1), *testStatefulSet.Spec.Replicas)
	assert.Equal(t, "test", testStatefulSet.Spec.Selector.MatchLabels["app"])
	assert.Equal(t, "test", testStatefulSet.Spec.Template.Labels["app"])
	assert.Equal(t, "testImage", testStatefulSet.Spec.Template.Spec.Containers[0].Image)
	assert.Equal(t, int64(0), *testStatefulSet.Spec.Template.Spec.TerminationGracePeriodSeconds)
	assert.Equal(t, int32(3), testStatefulSet.Spec.MinReadySeconds)
}

func TestRedefineWithReadinessProbe(t *testing.T) {
	testStatefulSet := DefineStatefulSet("testStatefulSet", "testNamespace", "testImage", map[string]string{"app": "test"})
	RedefineWithReadinessProbe(testStatefulSet)
	assert.NotNil(t, testStatefulSet.Spec.Template.Spec.Containers[0].ReadinessProbe)
}

func TestRedefineWithLivenessProbe(t *testing.T) {
	testStatefulSet := DefineStatefulSet("testStatefulSet", "testNamespace", "testImage", map[string]string{"app": "test"})
	RedefineWithLivenessProbe(testStatefulSet)
	assert.NotNil(t, testStatefulSet.Spec.Template.Spec.Containers[0].LivenessProbe)
}

func TestRedefineWithStartUpProbe(t *testing.T) {
	testStatefulSet := DefineStatefulSet("testStatefulSet", "testNamespace", "testImage", map[string]string{"app": "test"})
	RedefineWithStartUpProbe(testStatefulSet)
	assert.NotNil(t, testStatefulSet.Spec.Template.Spec.Containers[0].StartupProbe)
}

func TestRedefineWithContainerSpecs(t *testing.T) {
	testStatefulSet := DefineStatefulSet("testStatefulSet", "testNamespace", "testImage", map[string]string{"app": "test"})
	RedefineWithContainerSpecs(testStatefulSet, []corev1.Container{
		{
			Name: "test",
		},
	})
	assert.Equal(t, "test", testStatefulSet.Spec.Template.Spec.Containers[0].Name)
}

func TestRedefineWithReplicaNumber(t *testing.T) {
	testStatefulSet := DefineStatefulSet("testStatefulSet", "testNamespace", "testImage", map[string]string{"app": "test"})
	RedefineWithReplicaNumber(testStatefulSet, 2)
	assert.Equal(t, int32(2), *testStatefulSet.Spec.Replicas)
}

func TestRedefineWithPrivilegedContainer(t *testing.T) {
	testStatefulSet := DefineStatefulSet("testStatefulSet", "testNamespace", "testImage", map[string]string{"app": "test"})
	RedefineWithPrivilegedContainer(testStatefulSet)
	assert.Equal(t, true, *testStatefulSet.Spec.Template.Spec.Containers[0].SecurityContext.Privileged)
}

func TestRedefineWithPostStart(t *testing.T) {
	testStatefulSet := DefineStatefulSet("testStatefulSet", "testNamespace", "testImage", map[string]string{"app": "test"})
	RedefineWithPostStart(testStatefulSet)
	assert.Equal(t, "ls", testStatefulSet.Spec.Template.Spec.Containers[0].Lifecycle.PostStart.Exec.Command[0])
}
