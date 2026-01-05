package globalhelper

import (
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	k8sfake "k8s.io/client-go/kubernetes/fake"
)

func TestAddControlPlaneTaint(t *testing.T) {
	// define test node object
	defineNode := func() *corev1.Node {
		testNode := &corev1.Node{

			ObjectMeta: metav1.ObjectMeta{
				Name: "testNode",
				Labels: map[string]string{
					"test": "test",
				},
			},
			Spec: corev1.NodeSpec{
				Taints: []corev1.Taint{},
			},
		}

		return testNode
	}

	// Create a fake clientset
	var runtimeObjects []runtime.Object
	runtimeObjects = append(runtimeObjects, defineNode())
	client := k8sfake.NewClientset(runtimeObjects...)

	// Add the taint to the node
	err := addControlPlaneTaint(client.CoreV1(), defineNode())
	assert.Nil(t, err)

	// Check that the taint was added
	nodes, err := client.CoreV1().Nodes().List(t.Context(), metav1.ListOptions{})
	assert.Nil(t, err)

	for _, node := range nodes.Items {
		assert.Equal(t, 1, len(node.Spec.Taints))
	}
}

func TestRemoveControlPlaneTaint(t *testing.T) {
	// define test node object
	defineNode := func() *corev1.Node {
		testNode := &corev1.Node{

			ObjectMeta: metav1.ObjectMeta{
				Name: "testNode",
				Labels: map[string]string{
					"test": "test",
				},
			},
			Spec: corev1.NodeSpec{
				Taints: []corev1.Taint{},
			},
		}

		return testNode
	}

	// Create a fake clientset
	var runtimeObjects []runtime.Object
	runtimeObjects = append(runtimeObjects, defineNode())
	client := k8sfake.NewClientset(runtimeObjects...)

	// Add the taint to the node
	err := removeControlPlaneTaint(client.CoreV1(), defineNode())
	assert.Nil(t, err)

	// Check that the taint was added
	nodes, err := client.CoreV1().Nodes().List(t.Context(), metav1.ListOptions{})
	assert.Nil(t, err)

	for _, node := range nodes.Items {
		assert.Equal(t, 0, len(node.Spec.Taints))
	}
}

func TestNodeHasHugePagesEnabled(t *testing.T) {
	// define test node object
	defineNode := func(oneG, twoM string) *corev1.Node {
		testNode := &corev1.Node{

			ObjectMeta: metav1.ObjectMeta{
				Name: "testNode",
				Labels: map[string]string{
					"test": "test",
				},
			},
			Spec: corev1.NodeSpec{
				Taints: []corev1.Taint{},
			},
			Status: corev1.NodeStatus{
				Capacity: corev1.ResourceList{
					corev1.ResourceName("hugepages-1Gi"): resource.MustParse(oneG),
					corev1.ResourceName("hugepages-2Mi"): resource.MustParse(twoM),
				},
			},
		}

		return testNode
	}

	testCases := []struct {
		oneGValue                string
		twoMValue                string
		resourceName             string
		expectedHugePagesEnabled bool
	}{
		{
			resourceName:             "1Gi",
			oneGValue:                "0",
			twoMValue:                "0",
			expectedHugePagesEnabled: false,
		},
		{
			resourceName:             "1Gi",
			oneGValue:                "1Gi",
			twoMValue:                "0",
			expectedHugePagesEnabled: true,
		},
		{
			resourceName:             "1Gi",
			oneGValue:                "0",
			twoMValue:                "2Mi",
			expectedHugePagesEnabled: false,
		},
		{
			resourceName:             "2Mi",
			oneGValue:                "0",
			twoMValue:                "2Mi",
			expectedHugePagesEnabled: true,
		},
	}

	for _, tc := range testCases {
		node := defineNode(tc.oneGValue, tc.twoMValue)
		assert.Equal(t, tc.expectedHugePagesEnabled, NodeHasHugePagesEnabled(node, tc.resourceName))
	}
}
