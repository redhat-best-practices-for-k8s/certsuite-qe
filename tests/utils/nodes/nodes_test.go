package nodes

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	k8sfake "k8s.io/client-go/kubernetes/fake"
)

func TestEnsureAllNodesAreLabeled(t *testing.T) {
	// generate a k8s node object with labels
	generateNode := func(labelValue string) *corev1.Node {
		return &corev1.Node{
			ObjectMeta: metav1.ObjectMeta{
				Name: "nodeName",
				Labels: map[string]string{
					"test":      "test",
					"testValue": labelValue,
				},
			},
		}
	}

	testCases := []struct {
		nodeName       string
		testLabelValue string
	}{
		{
			testLabelValue: "test-value-1",
		},
		{
			testLabelValue: "",
		},
	}

	// Set UNIT_TEST env variable to true
	t.Setenv("UNIT_TEST", "true")

	for _, testCase := range testCases {
		var runtimeObjects []runtime.Object
		runtimeObjects = append(runtimeObjects, generateNode(testCase.testLabelValue))
		client := k8sfake.NewSimpleClientset(runtimeObjects...)

		// Label all the nodes accordingly
		err := EnsureAllNodesAreLabeled(client.CoreV1(), "testValue")
		assert.Nil(t, err)

		// Without generating a runtimeObject
		err = EnsureAllNodesAreLabeled(client.CoreV1(), "worker-cnf")
		assert.Nil(t, err)

		err = EnsureAllNodesAreLabeled(client.CoreV1(), "node-role.kubernetes.io/worker-cnf")
		assert.Nil(t, err)

		// Get all of the nodes from the fake client and test their labels
		nodes, err := client.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
		assert.Nil(t, err)
		assert.Equal(t, 1, len(nodes.Items))

		for _, node := range nodes.Items {
			assert.Equal(t, "test", node.Labels["test"])
			assert.Equal(t, testCase.testLabelValue, node.Labels["testValue"])
			assert.Equal(t, "", node.Labels["worker-cnf"])
			assert.Equal(t, "", node.Labels["node-role.kubernetes.io/worker-cnf"])
		}
	}
}
