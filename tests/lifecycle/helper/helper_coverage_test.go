package helper

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/onsi/gomega"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	k8sfake "k8s.io/client-go/kubernetes/fake"
	corev1Typed "k8s.io/client-go/kubernetes/typed/core/v1"

	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/globalhelper"
	tsparams "github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/lifecycle/parameters"
)

var errInjected = errors.New("injected failure")

const testNamespace = "lifecycle-ns"

// withFakeNodes installs a fake clientset holding the given nodes, restoring the
// global client on cleanup.
func withFakeNodes(t *testing.T, nodes ...runtime.Object) *k8sfake.Clientset {
	t.Helper()

	client := k8sfake.NewClientset(nodes...)
	globalhelper.SetTestK8sAPIClient(client)
	t.Cleanup(globalhelper.UnsetTestK8sAPIClient)

	return client
}

// stubWaitForNodesReady overrides the waitForNodesReady seam; the previous value
// is restored on cleanup.
func stubWaitForNodesReady(t *testing.T, err error) {
	t.Helper()

	oldWait := waitForNodesReady

	t.Cleanup(func() { waitForNodesReady = oldWait })

	waitForNodesReady = func(corev1Typed.NodeInterface, time.Duration, time.Duration) error {
		return err
	}
}

// testNode builds a node with an explicit Ready condition and schedulability.
func testNode(name string, ready, schedulable bool) *corev1.Node {
	node := &corev1.Node{
		ObjectMeta: metav1.ObjectMeta{Name: name},
		Spec:       corev1.NodeSpec{Unschedulable: !schedulable},
	}

	if ready {
		node.Status.Conditions = []corev1.NodeCondition{{
			Type:   corev1.NodeReady,
			Status: corev1.ConditionTrue,
		}}
	}

	return node
}

func TestDefineDeployment(t *testing.T) {
	t.Run("builds a deployment with the requested replicas and containers", func(t *testing.T) {
		t.Setenv("ENABLE_INFRASTRUCTURE_TOLERATIONS", "true")

		dep, err := DefineDeployment(3, 2, tsparams.TestDeploymentName, testNamespace)
		require.NoError(t, err)

		assert.Equal(t, tsparams.TestDeploymentName, dep.Name)
		assert.Equal(t, testNamespace, dep.Namespace)
		assert.Equal(t, tsparams.TestTargetLabels, dep.Labels)
		assert.Equal(t, tsparams.TestTargetLabels, dep.Spec.Selector.MatchLabels)
		assert.Equal(t, tsparams.TestTargetLabels, dep.Spec.Template.Labels)
		require.NotNil(t, dep.Spec.Replicas)
		assert.Equal(t, int32(3), *dep.Spec.Replicas)

		require.Len(t, dep.Spec.Template.Spec.Containers, 2)

		for _, container := range dep.Spec.Template.Spec.Containers {
			assert.Equal(t, tsparams.SampleWorkloadImage, container.Image)
			assert.Equal(t, []string{"/bin/bash", "-c", "sleep INF"}, container.Command)
		}

		// Container names come from the utils builders; keeping them distinct is
		// the wrapper's part of the contract.
		assert.NotEqual(t, dep.Spec.Template.Spec.Containers[0].Name,
			dep.Spec.Template.Spec.Containers[1].Name)
		assert.NotEmpty(t, dep.Spec.Template.Spec.Tolerations)
	})
}

func TestDefineDeploymentWithoutInfrastructureTolerations(t *testing.T) {
	t.Run("builds a deployment without tolerations even when they are enabled", func(t *testing.T) {
		t.Setenv("ENABLE_INFRASTRUCTURE_TOLERATIONS", "true")

		dep, err := DefineDeploymentWithoutInfrastructureTolerations(2, 1,
			tsparams.TestDeploymentName, testNamespace)
		require.NoError(t, err)

		assert.Equal(t, tsparams.TestDeploymentName, dep.Name)
		assert.Equal(t, testNamespace, dep.Namespace)
		require.NotNil(t, dep.Spec.Replicas)
		assert.Equal(t, int32(2), *dep.Spec.Replicas)
		require.Len(t, dep.Spec.Template.Spec.Containers, 1)
		assert.Equal(t, tsparams.SampleWorkloadImage, dep.Spec.Template.Spec.Containers[0].Image)
		assert.Empty(t, dep.Spec.Template.Spec.Tolerations)
	})
}

func TestDefineDeploymentRejectsZeroContainers(t *testing.T) {
	constructors := []struct {
		name      string
		construct func(replica int32, containers int, name, namespace string) (*appsv1.Deployment, error)
	}{
		{name: "DefineDeployment", construct: DefineDeployment},
		{name: "DefineDeploymentWithoutInfrastructureTolerations",
			construct: DefineDeploymentWithoutInfrastructureTolerations},
	}

	for _, constructor := range constructors {
		t.Run(constructor.name, func(t *testing.T) {
			_, err := constructor.construct(1, 0, tsparams.TestDeploymentName, testNamespace)
			require.Error(t, err)
			assert.EqualError(t, err, "invalid containers number")
		})
	}
}

func TestDefineReplicaSet(t *testing.T) {
	replicaSet := DefineReplicaSet(tsparams.TestReplicaSetName, testNamespace)

	assert.Equal(t, tsparams.TestReplicaSetName, replicaSet.Name)
	assert.Equal(t, testNamespace, replicaSet.Namespace)
	assert.Equal(t, tsparams.TestTargetLabels, replicaSet.Spec.Selector.MatchLabels)
	assert.Equal(t, tsparams.TestTargetLabels, replicaSet.Spec.Template.Labels)
	require.Len(t, replicaSet.Spec.Template.Spec.Containers, 1)
	assert.Equal(t, tsparams.SampleWorkloadImage, replicaSet.Spec.Template.Spec.Containers[0].Image)
}

func TestDefineStatefulSet(t *testing.T) {
	sts := DefineStatefulSet(tsparams.TestStatefulSetName, testNamespace)

	assert.Equal(t, tsparams.TestStatefulSetName, sts.Name)
	assert.Equal(t, testNamespace, sts.Namespace)
	assert.Equal(t, tsparams.TestTargetLabels, sts.Spec.Selector.MatchLabels)
	assert.Equal(t, tsparams.TestTargetLabels, sts.Spec.Template.Labels)
	require.Len(t, sts.Spec.Template.Spec.Containers, 1)
	assert.Equal(t, tsparams.SampleWorkloadImage, sts.Spec.Template.Spec.Containers[0].Image)
}

func TestDefinePod(t *testing.T) {
	pod := DefinePod(tsparams.TestPodName, testNamespace)

	assert.Equal(t, tsparams.TestPodName, pod.Name)
	assert.Equal(t, testNamespace, pod.Namespace)
	assert.Equal(t, tsparams.TestTargetLabels, pod.Labels)
	require.Len(t, pod.Spec.Containers, 1)
	assert.Equal(t, tsparams.SampleWorkloadImage, pod.Spec.Containers[0].Image)
}

func TestDefineDaemonSetWithImagePullPolicy(t *testing.T) {
	daemonSet := DefineDaemonSetWithImagePullPolicy(
		tsparams.TestDaemonSetName, testNamespace, "quay.io/test/image:1", corev1.PullAlways)

	assert.Equal(t, tsparams.TestDaemonSetName, daemonSet.Name)
	assert.Equal(t, testNamespace, daemonSet.Namespace)
	assert.Equal(t, tsparams.TestTargetLabels, daemonSet.Spec.Selector.MatchLabels)
	require.Len(t, daemonSet.Spec.Template.Spec.Containers, 1)
	assert.Equal(t, "quay.io/test/image:1", daemonSet.Spec.Template.Spec.Containers[0].Image)
	assert.Equal(t, corev1.PullAlways, daemonSet.Spec.Template.Spec.Containers[0].ImagePullPolicy)
}

func TestWaitUntilClusterIsStable(t *testing.T) {
	t.Run("returns nil on a ready, schedulable cluster", func(t *testing.T) {
		gomega.RegisterTestingT(t)
		withFakeNodes(t, testNode("node-1", true, true))

		require.NoError(t, WaitUntilClusterIsStable())
	})

	t.Run("uncordons unschedulable nodes", func(t *testing.T) {
		gomega.RegisterTestingT(t)

		client := withFakeNodes(t, testNode("node-1", true, false))

		require.NoError(t, WaitUntilClusterIsStable())

		nodesList, err := client.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
		require.NoError(t, err)
		require.Len(t, nodesList.Items, 1)
		assert.False(t, nodesList.Items[0].Spec.Unschedulable)
	})

	t.Run("wraps the node readiness wait failure", func(t *testing.T) {
		gomega.RegisterTestingT(t)
		withFakeNodes(t, testNode("node-1", true, true))
		stubWaitForNodesReady(t, errInjected)

		err := WaitUntilClusterIsStable()
		require.ErrorIs(t, err, errInjected)
		assert.Contains(t, err.Error(), "failed to wait for node to become ready")
	})
}
