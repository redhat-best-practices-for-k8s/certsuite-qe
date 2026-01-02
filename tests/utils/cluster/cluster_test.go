package cluster

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/runtime"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8sfake "k8s.io/client-go/kubernetes/fake"
)

func TestIsClusterStable(t *testing.T) {
	generateNode := func(unschedulable bool) *corev1.Node {
		return &corev1.Node{
			ObjectMeta: metav1.ObjectMeta{
				Name: "test-node",
			},
			Spec: corev1.NodeSpec{
				Unschedulable: unschedulable,
			},
		}
	}

	testCases := []struct {
		testUnschedulable bool
	}{
		{testUnschedulable: true}, // attempts an uncordon
		{testUnschedulable: false},
	}

	for _, testCase := range testCases {
		var runtimeObjects []runtime.Object
		runtimeObjects = append(runtimeObjects, generateNode(testCase.testUnschedulable))

		client := k8sfake.NewClientset(runtimeObjects...)
		result, err := IsClusterStable(client.CoreV1().Nodes())
		assert.Nil(t, err)
		assert.True(t, result)
	}
}
