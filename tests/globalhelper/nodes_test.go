package globalhelper

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	k8sfake "k8s.io/client-go/kubernetes/fake"
)

func TestEnableMasterScheduling(t *testing.T) {
	// define test node object
	defineNode := func(schedulable bool) *corev1.Node {
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

		if !schedulable {
			testNode.Spec.Taints = append(testNode.Spec.Taints, corev1.Taint{
				Key:    "node-role.kubernetes.io/control-plane",
				Effect: corev1.TaintEffectNoSchedule,
			})
			testNode.Labels["node-role.kubernetes.io/control-plane"] = ""
		}

		return testNode
	}

	testCases := []struct {
		defaultSchedulable bool
		schedulable        bool
	}{}

	for _, testCase := range testCases {
		// Create a fake clientset
		var runtimeObjects []runtime.Object
		runtimeObjects = append(runtimeObjects, defineNode(testCase.defaultSchedulable))
		client := k8sfake.NewSimpleClientset(runtimeObjects...)

		if testCase.schedulable {
			err := EnableMasterScheduling(client.CoreV1(), true)
			assert.Nil(t, err)

			// Check that all of the nodes are schedulable
			nodes, err := client.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
			assert.Nil(t, err)

			for _, node := range nodes.Items {
				assert.Equal(t, 0, len(node.Spec.Taints))
			}
		} else {
			err := EnableMasterScheduling(client.CoreV1(), false)
			assert.Nil(t, err)

			// Check that all of the nodes are not schedulable
			nodes, err := client.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
			assert.Nil(t, err)

			for _, node := range nodes.Items {
				assert.Equal(t, 1, len(node.Spec.Taints))
			}
		}
	}
}

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
	client := k8sfake.NewSimpleClientset(runtimeObjects...)

	// Add the taint to the node
	err := addControlPlaneTaint(client.CoreV1(), defineNode())
	assert.Nil(t, err)

	// Check that the taint was added
	nodes, err := client.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
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
	client := k8sfake.NewSimpleClientset(runtimeObjects...)

	// Add the taint to the node
	err := removeControlPlaneTaint(client.CoreV1(), defineNode())
	assert.Nil(t, err)

	// Check that the taint was added
	nodes, err := client.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	assert.Nil(t, err)

	for _, node := range nodes.Items {
		assert.Equal(t, 0, len(node.Spec.Taints))
	}
}
