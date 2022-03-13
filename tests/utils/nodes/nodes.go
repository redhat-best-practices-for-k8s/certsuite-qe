package nodes

import (
	"context"
	"encoding/json"
	"time"

	"github.com/golang/glog"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/client"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/wait"
)

type patchStringValue struct {
	Operation string `json:"op"`
	Path      string `json:"path"`
	Value     bool   `json:"value"`
}

// WaitForNodesReady waits for all nodes become ready.
func WaitForNodesReady(cs *client.ClientSet, timeout, interval time.Duration) error {
	return wait.PollImmediate(interval, timeout, func() (bool, error) {
		nodesList, err := cs.Nodes().List(context.Background(), metav1.ListOptions{})
		if err != nil {
			return false, nil
		}
		for _, node := range nodesList.Items {
			if !IsNodeInCondition(&node, corev1.NodeReady) {
				return false, nil
			}
		}
		glog.V(5).Info("All nodes are Ready")

		return true, nil
	})
}

// IsNodeInCondition parses node conditions. Returns true if node is in given condition, otherwise false.
func IsNodeInCondition(node *corev1.Node, condition corev1.NodeConditionType) bool {
	for _, c := range node.Status.Conditions {
		if c.Type == condition && c.Status == corev1.ConditionTrue {
			return true
		}
	}

	return false
}

// GetNumOfReadyNodesInCluster gets the number of ready nodes in the cluster.
func GetNumOfReadyNodesInCluster(cs *client.ClientSet) (int32, error) {
	nodesList, err := cs.Nodes().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return 0, err
	}

	numOfNodesExistsInCluster := len(nodesList.Items)

	for _, node := range nodesList.Items {
		if !IsNodeInCondition(&node, corev1.NodeReady) {
			numOfNodesExistsInCluster--
		}
	}

	return int32(numOfNodesExistsInCluster), nil
}

// CordonOrUnCordonNodeByName cordones/uncordones a node by a given node name.
func CordonOrUnCordonNodeByName(nodeName string, clients *client.ClientSet, cordon bool) (bool, error) {
	payload := []patchStringValue{{
		Operation: "replace",
		Path:      "/spec/unschedulable",
		Value:     cordon,
	}}
	payloadBytes, _ := json.Marshal(payload)

	_, err := clients.Nodes().Patch(context.Background(), nodeName, types.JSONPatchType,
		payloadBytes, metav1.PatchOptions{})

	if err != nil {
		return false, err
	}

	return true, nil
}
