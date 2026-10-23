package cluster

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	k8sfake "k8s.io/client-go/kubernetes/fake"
	k8stesting "k8s.io/client-go/testing"
)

var errInjected = errors.New("injected failure")

func generateNode(name string, unschedulable bool) *corev1.Node {
	return &corev1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Spec: corev1.NodeSpec{
			Unschedulable: unschedulable,
		},
	}
}

// newFailingNodesClient builds a fake clientset whose node calls for the given verb always fail.
func newFailingNodesClient(verb string, objects ...runtime.Object) *k8sfake.Clientset {
	client := k8sfake.NewClientset(objects...)
	client.PrependReactor(verb, "nodes", func(k8stesting.Action) (bool, runtime.Object, error) {
		return true, nil, errInjected
	})

	return client
}

func TestIsClusterStable(t *testing.T) {
	testCases := []struct {
		testUnschedulable bool
	}{
		{testUnschedulable: true}, // attempts an uncordon
		{testUnschedulable: false},
	}

	for _, testCase := range testCases {
		client := k8sfake.NewClientset(generateNode("test-node", testCase.testUnschedulable))

		result, err := IsClusterStable(client.CoreV1().Nodes())
		assert.Nil(t, err)
		assert.True(t, result)
	}
}

func TestIsClusterStableWithEmptyNodeList(t *testing.T) {
	client := k8sfake.NewClientset()

	result, err := IsClusterStable(client.CoreV1().Nodes())
	assert.Nil(t, err)
	assert.True(t, result)
}

func TestIsClusterStableWithMultipleNodes(t *testing.T) {
	t.Run("uncordons every cordoned node and leaves schedulable ones alone", func(t *testing.T) {
		client := k8sfake.NewClientset(
			generateNode("schedulable-node", false),
			generateNode("cordoned-node-1", true),
			generateNode("cordoned-node-2", true))

		var updates int

		client.PrependReactor("update", "nodes", func(k8stesting.Action) (bool, runtime.Object, error) {
			updates++

			return false, nil, nil
		})

		result, err := IsClusterStable(client.CoreV1().Nodes())
		assert.Nil(t, err)
		assert.True(t, result)
		assert.Equal(t, 2, updates, "only the cordoned nodes should be updated")

		for _, nodeName := range []string{"cordoned-node-1", "cordoned-node-2"} {
			node, err := client.CoreV1().Nodes().Get(t.Context(), nodeName, metav1.GetOptions{})
			require.Nil(t, err)
			assert.False(t, node.Spec.Unschedulable, "node %s should have been uncordoned", nodeName)
		}
	})
}

func TestIsClusterStableFailures(t *testing.T) {
	t.Run("list failure", func(t *testing.T) {
		client := newFailingNodesClient("list")

		result, err := IsClusterStable(client.CoreV1().Nodes())
		assert.False(t, result)
		assert.ErrorIs(t, err, errInjected)
	})

	t.Run("uncordon update failure", func(t *testing.T) {
		client := newFailingNodesClient("update", generateNode("test-node", true))

		result, err := IsClusterStable(client.CoreV1().Nodes())
		assert.False(t, result)
		assert.ErrorIs(t, err, errInjected)
	})

	t.Run("get failure after uncordon", func(t *testing.T) {
		client := k8sfake.NewClientset(generateNode("test-node", true))

		var updated bool

		client.PrependReactor("update", "nodes", func(k8stesting.Action) (bool, runtime.Object, error) {
			updated = true

			return false, nil, nil
		})

		client.PrependReactor("get", "nodes", func(k8stesting.Action) (bool, runtime.Object, error) {
			if updated {
				return true, nil, errInjected
			}

			return false, nil, nil
		})

		result, err := IsClusterStable(client.CoreV1().Nodes())
		assert.False(t, result)
		assert.ErrorIs(t, err, errInjected)
	})

	t.Run("node remains unschedulable after uncordon", func(t *testing.T) {
		client := k8sfake.NewClientset(generateNode("test-node", true))

		var updated bool

		client.PrependReactor("update", "nodes", func(k8stesting.Action) (bool, runtime.Object, error) {
			updated = true

			return false, nil, nil
		})

		client.PrependReactor("get", "nodes", func(k8stesting.Action) (bool, runtime.Object, error) {
			// Simulates the node being re-cordoned after the uncordon updated it.
			if updated {
				return true, generateNode("test-node", true), nil
			}

			return false, nil, nil
		})

		result, err := IsClusterStable(client.CoreV1().Nodes())
		assert.False(t, result)
		assert.Nil(t, err)
	})
}
