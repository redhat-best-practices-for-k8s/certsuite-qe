package nodes

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/globalhelper"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	k8sfake "k8s.io/client-go/kubernetes/fake"
	k8stesting "k8s.io/client-go/testing"
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

var errInjected = errors.New("injected failure")

// generateNodeWithNameAndCondition generates a node with the given name and a single status condition.
func generateNodeWithNameAndCondition(name string, conditionType corev1.NodeConditionType,
	status corev1.ConditionStatus) *corev1.Node {
	return &corev1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
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

// generateMasterNode generates a node labeled as a control plane node.
func generateMasterNode(labelValue string) *corev1.Node {
	node := generateNode(labelValue)
	node.Labels["node-role.kubernetes.io/control-plane"] = "true"

	return node
}

// newFailingNodesClient builds a fake clientset whose node calls for the given verb always fail.
func newFailingNodesClient(verb string, objects ...runtime.Object) *k8sfake.Clientset {
	client := k8sfake.NewClientset(objects...)
	client.PrependReactor(verb, "nodes", func(k8stesting.Action) (bool, runtime.Object, error) {
		return true, nil, errInjected
	})

	return client
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
		client := k8sfake.NewClientset(runtimeObjects...)

		// Get all of the nodes from the fake client and test their labels
		nodes, err := client.CoreV1().Nodes().List(t.Context(), metav1.ListOptions{})
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

func TestEnsureAllNodesAreLabeledViaAPIClient(t *testing.T) {
	apiClient := k8sfake.NewClientset(generateNode("test-value-1"))
	globalhelper.SetTestK8sAPIClient(apiClient)
	t.Cleanup(globalhelper.UnsetTestK8sAPIClient)

	assert.Nil(t, EnsureAllNodesAreLabeled("worker-cnf"))

	labeledNodes, err := apiClient.CoreV1().Nodes().List(t.Context(), metav1.ListOptions{})
	assert.Nil(t, err)
	assert.Equal(t, 1, len(labeledNodes.Items))
	assert.Contains(t, labeledNodes.Items[0].Labels, "worker-cnf")
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

		client := k8sfake.NewClientset(runtimeObjects...)
		assert.Nil(t, addControlPlaneTaint(client.CoreV1().Nodes(), testNode))

		// Get all of the nodes from the fake client and test their labels
		nodes, err := client.CoreV1().Nodes().List(t.Context(), metav1.ListOptions{})
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

		client := k8sfake.NewClientset(runtimeObjects...)
		assert.Nil(t, removeControlPlaneTaint(client.CoreV1().Nodes(), testNode))

		// Get all of the nodes from the fake client and test their labels
		nodes, err := client.CoreV1().Nodes().List(t.Context(), metav1.ListOptions{})
		assert.Nil(t, err)
		assert.Equal(t, 1, len(nodes.Items))

		for _, node := range nodes.Items {
			assert.Equal(t, 0, len(node.Spec.Taints))
		}
	}
}

func TestIsMasterNode(t *testing.T) {
	testCases := []struct {
		masterLabel string
		expected    bool
	}{
		{
			expected: false,
		},
		{
			masterLabel: "node-role.kubernetes.io/control-plane",
			expected:    true,
		},
		{
			masterLabel: "node-role.kubernetes.io/master",
			expected:    true,
		},
	}

	nodesClient := k8sfake.NewClientset().CoreV1().Nodes()

	for _, testCase := range testCases {
		testNode := generateNode("test-value-1")
		if testCase.masterLabel != "" {
			testNode.Labels[testCase.masterLabel] = "true"
		}

		result, err := IsNodeMaster(testNode, nodesClient)
		assert.Nil(t, err)
		assert.Equal(t, testCase.expected, result)
	}
}

func TestIsNodeInCondition(t *testing.T) {
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
		testNode := generateNodeWithNameAndCondition("nodeName", testCase.conditionType, testCase.status)
		result := IsNodeInCondition(testNode, testCase.conditionType)
		assert.Equal(t, testCase.expected, result)
	}
}

func TestGetNumOfReadyNodesInCluster(t *testing.T) {
	testCases := []struct {
		readyNodes    int
		notReadyNodes int
	}{
		{
			readyNodes: 1,
		},
		{
			readyNodes: 2,
		},
		{
			readyNodes:    2,
			notReadyNodes: 1,
		},
	}

	for _, testCase := range testCases {
		var runtimeObjects []runtime.Object
		for i := 0; i < testCase.readyNodes; i++ {
			runtimeObjects = append(runtimeObjects, generateNodeWithNameAndCondition(
				fmt.Sprintf("node-%d", i), corev1.NodeReady, corev1.ConditionTrue))
		}

		for i := 0; i < testCase.notReadyNodes; i++ {
			runtimeObjects = append(runtimeObjects, generateNodeWithNameAndCondition(
				fmt.Sprintf("not-ready-%d", i), corev1.NodeReady, corev1.ConditionFalse))
		}

		client := k8sfake.NewClientset(runtimeObjects...)
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

	client := k8sfake.NewClientset(runtimeObjects...)
	assert.Nil(t, UnCordon(client.CoreV1().Nodes(), testNode.Name))

	// Get all of the nodes from the fake client and test their labels
	nodes, err := client.CoreV1().Nodes().List(t.Context(), metav1.ListOptions{})
	assert.Nil(t, err)
	assert.Equal(t, 1, len(nodes.Items))

	for _, node := range nodes.Items {
		assert.False(t, node.Spec.Unschedulable)
	}
}

func TestEnableMasterScheduling(t *testing.T) {
	testNode := generateMasterNode("test-value-1")

	var runtimeObjects []runtime.Object
	runtimeObjects = append(runtimeObjects, testNode)

	client := k8sfake.NewClientset(runtimeObjects...)
	assert.Nil(t, EnableMasterScheduling(client.CoreV1().Nodes(), true))

	// Get all of the nodes from the fake client and test their labels
	nodes, err := client.CoreV1().Nodes().List(t.Context(), metav1.ListOptions{})
	assert.Nil(t, err)
	assert.Equal(t, 1, len(nodes.Items))

	// Assert the node has correct taint
	assert.Zero(t, nodes.Items[0].Spec.Taints)

	// Disable master scheduling
	assert.Nil(t, EnableMasterScheduling(client.CoreV1().Nodes(), false))
	nodes, err = client.CoreV1().Nodes().List(t.Context(), metav1.ListOptions{})
	assert.Nil(t, err)
	assert.Equal(t, 1, len(nodes.Items))
	assert.NotZero(t, nodes.Items[0].Spec.Taints)
}

func TestGetNumOfReadyNodesInClusterListFailure(t *testing.T) {
	client := newFailingNodesClient("list")

	numOfReadyNodes, err := GetNumOfReadyNodesInCluster(client.CoreV1().Nodes())
	assert.ErrorIs(t, err, errInjected)
	assert.Equal(t, int32(0), numOfReadyNodes)
}

func TestUnCordonFailures(t *testing.T) {
	t.Run("get failure", func(t *testing.T) {
		client := newFailingNodesClient("get")

		assert.ErrorIs(t, UnCordon(client.CoreV1().Nodes(), "nodeName"), errInjected)
	})

	t.Run("update failure", func(t *testing.T) {
		client := newFailingNodesClient("update", generateNode("test-value-1"))

		err := UnCordon(client.CoreV1().Nodes(), "nodeName")
		assert.ErrorIs(t, err, errInjected)
		assert.Contains(t, err.Error(), "failed to update node nodeName")
	})
}

func TestEnsureAllNodesAreLabeledFailures(t *testing.T) {
	t.Run("list failure", func(t *testing.T) {
		client := newFailingNodesClient("list")

		err := ensureAllNodesAreLabeled(client.CoreV1(), "worker-cnf")
		assert.ErrorIs(t, err, errInjected)
		assert.Contains(t, err.Error(), "failed to get nodes")
	})

	t.Run("label update failure", func(t *testing.T) {
		client := newFailingNodesClient("update", generateNode("test-value-1"))

		// The node lacks the label, so labeling it updates the node and fails.
		assert.ErrorIs(t, ensureAllNodesAreLabeled(client.CoreV1(), "worker-cnf"), errInjected)
	})
}

func TestEnableMasterSchedulingFailures(t *testing.T) {
	t.Run("list failure", func(t *testing.T) {
		client := newFailingNodesClient("list")

		err := EnableMasterScheduling(client.CoreV1().Nodes(), true)
		assert.ErrorIs(t, err, errInjected)
		assert.Contains(t, err.Error(), "failed to get nodes")
	})

	t.Run("add taint update failure", func(t *testing.T) {
		client := newFailingNodesClient("update", generateMasterNode("test-value-1"))

		err := EnableMasterScheduling(client.CoreV1().Nodes(), false)
		assert.ErrorIs(t, err, errInjected)
		assert.Contains(t, err.Error(), "failed to set node nodeName schedulable value")
		assert.Contains(t, err.Error(), "failed to update node nodeName")
	})

	t.Run("remove taint update failure", func(t *testing.T) {
		masterNode := generateMasterNode("test-value-1")
		masterNode.Spec.Taints = append(masterNode.Spec.Taints, corev1.Taint{
			Key:    controlPlaneTaintKey,
			Effect: corev1.TaintEffectNoSchedule,
		})

		client := newFailingNodesClient("update", masterNode)

		err := EnableMasterScheduling(client.CoreV1().Nodes(), true)
		assert.ErrorIs(t, err, errInjected)
		assert.Contains(t, err.Error(), "failed to set node nodeName schedulable value")
		assert.Contains(t, err.Error(), "failed to update node nodeName")
	})
}

func TestWaitForNodesReady(t *testing.T) {
	t.Run("returns nil when all nodes are ready", func(t *testing.T) {
		client := k8sfake.NewClientset(
			generateNodeWithNameAndCondition("node-1", corev1.NodeReady, corev1.ConditionTrue),
			generateNodeWithNameAndCondition("node-2", corev1.NodeReady, corev1.ConditionTrue))

		assert.Nil(t, WaitForNodesReady(client.CoreV1().Nodes(), time.Second, 10*time.Millisecond))
	})

	t.Run("retries after a list failure", func(t *testing.T) {
		client := k8sfake.NewClientset(
			generateNodeWithNameAndCondition("node-1", corev1.NodeReady, corev1.ConditionTrue))

		var listCalls int

		client.PrependReactor("list", "nodes", func(k8stesting.Action) (bool, runtime.Object, error) {
			listCalls++
			if listCalls == 1 {
				return true, nil, errInjected
			}

			return false, nil, nil
		})

		assert.Nil(t, WaitForNodesReady(client.CoreV1().Nodes(), time.Second, time.Millisecond))
		assert.Equal(t, 2, listCalls)
	})

	t.Run("times out when a node is never ready", func(t *testing.T) {
		client := k8sfake.NewClientset(
			generateNodeWithNameAndCondition("ready-node", corev1.NodeReady, corev1.ConditionTrue),
			generateNodeWithNameAndCondition("not-ready-node", corev1.NodeReady, corev1.ConditionFalse))

		err := WaitForNodesReady(client.CoreV1().Nodes(), 50*time.Millisecond, 5*time.Millisecond)
		assert.ErrorIs(t, err, context.DeadlineExceeded)
	})

	t.Run("times out when listing keeps failing", func(t *testing.T) {
		client := newFailingNodesClient("list")

		err := WaitForNodesReady(client.CoreV1().Nodes(), 50*time.Millisecond, 5*time.Millisecond)
		assert.ErrorIs(t, err, context.DeadlineExceeded)
	})
}
