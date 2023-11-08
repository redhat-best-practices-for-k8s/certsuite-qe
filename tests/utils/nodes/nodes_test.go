package nodes

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	k8sfake "k8s.io/client-go/kubernetes/fake"
)

// generate a k8s node object with labels.
func generateNode(labelValue string) *corev1.Node {
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

func TestEnsureAllNodesAreLabeled(t *testing.T) {
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

		// Get all of the nodes from the fake client and test their labels
		nodes, err := client.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
		assert.Nil(t, err)
		assert.Equal(t, 1, len(nodes.Items))

		// Label all the nodes accordingly
		err = ensureAllNodesAreLabeled(client.CoreV1(), "testValue")
		assert.Nil(t, err)

		// Without generating a runtimeObject
		err = ensureAllNodesAreLabeled(client.CoreV1(), "worker-cnf")
		assert.Nil(t, err)

		err = ensureAllNodesAreLabeled(client.CoreV1(), "node-role.kubernetes.io/worker-cnf")
		assert.Nil(t, err)

		for _, node := range nodes.Items {
			assert.Equal(t, "test", node.Labels["test"])
			assert.Equal(t, testCase.testLabelValue, node.Labels["testValue"])
			assert.Equal(t, "", node.Labels["worker-cnf"])
			assert.Equal(t, "", node.Labels["node-role.kubernetes.io/worker-cnf"])
		}
	}
}

func TestAddControlPlaneTaint(t *testing.T) {
	testCases := []struct {
		taintAlreadyExists bool
	}{
		{
			taintAlreadyExists: false,
		},
		{
			taintAlreadyExists: true,
		},
	}

	for _, testCase := range testCases {
		testNode := generateNode("test-value-1")
		if testCase.taintAlreadyExists {
			testNode.Spec.Taints = append(testNode.Spec.Taints, corev1.Taint{
				Key:    "node-role.kubernetes.io/control-plane",
				Effect: corev1.TaintEffectNoSchedule,
			})
		}

		var runtimeObjects []runtime.Object
		runtimeObjects = append(runtimeObjects, testNode)

		client := k8sfake.NewSimpleClientset(runtimeObjects...)
		assert.Nil(t, addControlPlaneTaint(client.CoreV1().Nodes(), testNode))

		// Get all of the nodes from the fake client and test their labels
		nodes, err := client.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
		assert.Nil(t, err)
		assert.Equal(t, 1, len(nodes.Items))

		for _, node := range nodes.Items {
			assert.Equal(t, node.Spec.Taints[0].Key, "node-role.kubernetes.io/control-plane")
			assert.Equal(t, node.Spec.Taints[0].Effect, corev1.TaintEffectNoSchedule)
		}
	}
}

func TestRemoveControlPlaneTaint(t *testing.T) {
	testCases := []struct {
		taintAlreadyExists bool
	}{
		{
			taintAlreadyExists: true,
		},
		{
			taintAlreadyExists: false,
		},
	}

	for _, testCase := range testCases {
		testNode := generateNode("test-value-1")
		if testCase.taintAlreadyExists {
			testNode.Spec.Taints = append(testNode.Spec.Taints, corev1.Taint{
				Key:    "node-role.kubernetes.io/control-plane",
				Effect: corev1.TaintEffectNoSchedule,
			})
		}

		var runtimeObjects []runtime.Object
		runtimeObjects = append(runtimeObjects, testNode)

		client := k8sfake.NewSimpleClientset(runtimeObjects...)
		assert.Nil(t, removeControlPlaneTaint(client.CoreV1().Nodes(), testNode))

		// Get all of the nodes from the fake client and test their labels
		nodes, err := client.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
		assert.Nil(t, err)
		assert.Equal(t, 1, len(nodes.Items))

		for _, node := range nodes.Items {
			assert.Equal(t, 0, len(node.Spec.Taints))
		}
	}
}

func TestIsMasterNode(t *testing.T) {
	testCases := []struct {
		isMasterNode       bool
		expectedMasterNode bool
	}{
		{
			isMasterNode: false,
		},
		{
			isMasterNode: true,
		},
	}

	for _, testCase := range testCases {
		testNode := generateNode("test-value-1")
		if testCase.isMasterNode {
			testNode.Labels["node-role.kubernetes.io/control-plane"] = "true"
			result, err := IsNodeMaster(testNode, k8sfake.NewSimpleClientset().CoreV1().Nodes())
			assert.Nil(t, err)
			assert.True(t, result)
		} else {
			result, err := IsNodeMaster(testNode, k8sfake.NewSimpleClientset().CoreV1().Nodes())
			assert.Nil(t, err)
			assert.False(t, result)
		}
	}
}

func TestIsNodeInCondition(t *testing.T) {
	generateNodeWithCondition := func(conditionType corev1.NodeConditionType, status corev1.ConditionStatus) *corev1.Node {
		return &corev1.Node{
			ObjectMeta: metav1.ObjectMeta{
				Name: "nodeName",
			},
			Status: corev1.NodeStatus{
				Conditions: []corev1.NodeCondition{
					{
						Type:   conditionType,
						Status: status,
					},
				},
			},
		}
	}

	testCases := []struct {
		conditionType corev1.NodeConditionType
		status        corev1.ConditionStatus
		expected      bool
	}{
		{
			conditionType: corev1.NodeReady,
			status:        corev1.ConditionTrue,
			expected:      true,
		},
		{
			conditionType: corev1.NodeReady,
			status:        corev1.ConditionFalse,
			expected:      false,
		},
		{
			conditionType: corev1.NodeReady,
			status:        corev1.ConditionUnknown,
			expected:      false,
		},
		{
			conditionType: corev1.NodeMemoryPressure,
			status:        corev1.ConditionTrue,
			expected:      true,
		},
	}

	for _, testCase := range testCases {
		testNode := generateNodeWithCondition(testCase.conditionType, testCase.status)
		result := IsNodeInCondition(testNode, testCase.conditionType)
		assert.Equal(t, testCase.expected, result)
	}
}

func TestGetNumOfReadyNodesInCluster(t *testing.T) {
	generateNodeWithName := func(name string) *corev1.Node {
		return &corev1.Node{
			ObjectMeta: metav1.ObjectMeta{
				Name: name,
			},
			Status: corev1.NodeStatus{
				Conditions: []corev1.NodeCondition{
					{
						Type:   corev1.NodeReady,
						Status: corev1.ConditionTrue,
					},
				},
			},
		}
	}

	testCases := []struct {
		readyNodes int
	}{
		{
			readyNodes: 1,
		},
		{
			readyNodes: 2,
		},
	}

	for _, testCase := range testCases {
		var runtimeObjects []runtime.Object
		for i := 0; i < testCase.readyNodes; i++ {
			runtimeObjects = append(runtimeObjects, generateNodeWithName(fmt.Sprintf("node-%d", i)))
		}

		client := k8sfake.NewSimpleClientset(runtimeObjects...)
		result, err := GetNumOfReadyNodesInCluster(client.CoreV1().Nodes())
		assert.Nil(t, err)
		assert.Equal(t, testCase.readyNodes, int(result))
	}
}

func TestUnCordon(t *testing.T) {
	testNode := generateNode("test-value-1")
	testNode.Spec.Unschedulable = true

	var runtimeObjects []runtime.Object
	runtimeObjects = append(runtimeObjects, testNode)

	client := k8sfake.NewSimpleClientset(runtimeObjects...)
	assert.Nil(t, UnCordon(client.CoreV1().Nodes(), testNode.Name))

	// Get all of the nodes from the fake client and test their labels
	nodes, err := client.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	assert.Nil(t, err)
	assert.Equal(t, 1, len(nodes.Items))

	for _, node := range nodes.Items {
		assert.False(t, node.Spec.Unschedulable)
	}
}

func TestEnableMasterScheduling(t *testing.T) {
	testNode := generateNode("test-value-1")
	// mark node as a control plane node
	testNode.Labels["node-role.kubernetes.io/control-plane"] = "true"

	var runtimeObjects []runtime.Object
	runtimeObjects = append(runtimeObjects, testNode)

	client := k8sfake.NewSimpleClientset(runtimeObjects...)
	assert.Nil(t, EnableMasterScheduling(client.CoreV1().Nodes(), true))

	// Get all of the nodes from the fake client and test their labels
	nodes, err := client.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	assert.Nil(t, err)
	assert.Equal(t, 1, len(nodes.Items))

	// Assert the node has correct taint
	assert.Zero(t, nodes.Items[0].Spec.Taints)

	// Disable master scheduling
	assert.Nil(t, EnableMasterScheduling(client.CoreV1().Nodes(), false))
	nodes, err = client.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	assert.Nil(t, err)
	assert.Equal(t, 1, len(nodes.Items))
	assert.NotZero(t, nodes.Items[0].Spec.Taints)
}
