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
	corev1Typed "k8s.io/client-go/kubernetes/typed/core/v1"
)

type resourceSpecs struct {
	Operation string `json:"op"`
	Path      string `json:"path"`
	Value     bool   `json:"value"`
}

// WaitForNodesReady waits for all the nodes to become ready.
func WaitForNodesReady(clients *client.ClientSet, timeout, interval time.Duration) error {
	return wait.PollUntilContextTimeout(context.TODO(), interval, timeout, true,
		func(ctx context.Context) (bool, error) {
			nodesList, err := clients.Nodes().List(ctx, metav1.ListOptions{})
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
	nodesList, err := clients.Nodes().List(context.TODO(), metav1.ListOptions{})
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

	_, err = clients.Nodes().Patch(context.TODO(), nodeName, types.JSONPatchType,
		cordonPatchBytes, metav1.PatchOptions{})
	if err != nil {
		return fmt.Errorf("failed to patch node unschedulable value: %w", err)
	}

	return nil
}

func IsNodeMaster(name string, clients *client.ClientSet) (bool, error) {
	node, err := clients.Nodes().Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return false, err
	}

	masterLabels := []string{"node-role.kubernetes.io/master", "node-role.kubernetes.io/control-plane"}

	for _, label := range masterLabels {
		if _, exists := node.Labels[label]; exists {
			return true, nil
		}
	}

	return false, nil
}

// EnsureAllNodesAreLabeled ensures that all nodes are labeled with the given label.
func EnsureAllNodesAreLabeled(client corev1Typed.CoreV1Interface, label string) error {
	nodesList, err := client.Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return err
	}

	for _, node := range nodesList.Items {
		if _, exists := node.Labels[label]; !exists {
			err = LabelNode(client, &node, label, "")

			if err != nil {
				return err
			}
		}
	}

	return nil
}

// LabelNode labels a node by a given node name.
func LabelNode(client corev1Typed.CoreV1Interface, node *corev1.Node, label, value string) error {
	// Set the label
	node.Labels[label] = value

	var err error

	_, err = client.Nodes().Update(context.TODO(), node, metav1.UpdateOptions{})

	return err
}
