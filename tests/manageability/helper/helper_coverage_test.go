package helper

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"

	tsparams "github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/manageability/parameters"
)

func TestDefineManageabilityPod(t *testing.T) {
	labels := map[string]string{"app": "test"}

	pod := DefineManageabilityPod(tsparams.TestPodName, tsparams.ManageabilityNamespace,
		"quay.io/test/image:1", labels)

	assert.Equal(t, tsparams.TestPodName, pod.Name)
	assert.Equal(t, tsparams.ManageabilityNamespace, pod.Namespace)
	assert.Equal(t, labels, pod.Labels)
	require.NotNil(t, pod.Spec.TerminationGracePeriodSeconds)
	assert.Equal(t, int64(0), *pod.Spec.TerminationGracePeriodSeconds)

	require.Len(t, pod.Spec.Containers, 1)

	container := pod.Spec.Containers[0]
	assert.Equal(t, "httpd", container.Name)
	assert.Equal(t, "quay.io/test/image:1", container.Image)
	require.Len(t, container.Ports, 1)
	assert.Equal(t, "http", container.Ports[0].Name)
	assert.Equal(t, int32(80), container.Ports[0].ContainerPort)
}

func TestRedefinePodWithContainerPort(t *testing.T) {
	newPod := func() *corev1.Pod {
		return DefineManageabilityPod(tsparams.TestPodName, tsparams.ManageabilityNamespace,
			"quay.io/test/image:1", nil)
	}

	t.Run("renames the port of the given container", func(t *testing.T) {
		pod := newPod()

		RedefinePodWithContainerPort(pod, 0, "renamed-port")
		assert.Equal(t, "renamed-port", pod.Spec.Containers[0].Ports[0].Name)
	})

	t.Run("out-of-range indices leave the pod untouched", func(t *testing.T) {
		pod := newPod()

		RedefinePodWithContainerPort(pod, -1, "renamed-port")
		RedefinePodWithContainerPort(pod, 1, "renamed-port")
		assert.Equal(t, "http", pod.Spec.Containers[0].Ports[0].Name)
	})
}
