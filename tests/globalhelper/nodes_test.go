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
