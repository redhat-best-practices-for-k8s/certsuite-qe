package nodes

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/golang/glog"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/client"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/wait"
)

type resourceSpecs struct {
	Operation string `json:"op"`
	Path      string `json:"path"`
	Value     bool   `json:"value"`
}

// WaitForNodesReady waits for all the nodes to become ready.
func WaitForNodesReady(clients *client.ClientSet, timeout, interval time.Duration) error {
	return wait.PollImmediate(interval, timeout, func() (bool, error) {
		nodesList, err := clients.Nodes().List(context.Background(), metav1.ListOptions{})
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
func GetNumOfReadyNodesInCluster(clients *client.ClientSet) (int32, error) {
	nodesList, err := clients.Nodes().List(context.Background(), metav1.ListOptions{})
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

// UnCordon removes cordon label from the given node.
func UnCordon(clients *client.ClientSet, nodeName string) error {
	return setUnSchedulableValue(clients, nodeName, false)
}

// setUnSchedulableValue cordones/uncordons a node by a given node name.
func setUnSchedulableValue(clients *client.ClientSet, nodeName string, unSchedulable bool) error {
	cordonPatchBytes, err := json.Marshal(
		[]resourceSpecs{{
			Operation: "replace",
			Path:      "/spec/unschedulable",
			Value:     unSchedulable,
		}})

	if err != nil {
		return err
	}

	_, err = clients.Nodes().Patch(context.Background(), nodeName, types.JSONPatchType,
		cordonPatchBytes, metav1.PatchOptions{})
	if err != nil {
		return fmt.Errorf("failed to patch node unschedulable value: %w", err)
	}

	return nil
}

func IsNodeMaster(name string, clients *client.ClientSet) (bool, error) {
	node, err := clients.Nodes().Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		return false, err
	}

	if _, exists := node.Labels["node-role.kubernetes.io/master"]; exists {
		return true, nil
	}

	return false, nil
}
