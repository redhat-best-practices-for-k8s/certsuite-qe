package nodes

import (
	"context"
	"fmt"
	"time"

	"github.com/golang/glog"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	corev1Typed "k8s.io/client-go/kubernetes/typed/core/v1"
)

const (
	controlPlaneTaintKey = "node-role.kubernetes.io/control-plane"
	masterTaintKey       = "node-role.kubernetes.io/master"
)

// WaitForNodesReady waits for all the nodes to become ready.
func WaitForNodesReady(client corev1Typed.NodeInterface, timeout, interval time.Duration) error {
	return wait.PollUntilContextTimeout(context.TODO(), interval, timeout, true,
		func(ctx context.Context) (bool, error) {
			nodesList, err := client.List(ctx, metav1.ListOptions{})
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
func GetNumOfReadyNodesInCluster(client corev1Typed.NodeInterface) (int32, error) {
	nodesList, err := client.List(context.TODO(), metav1.ListOptions{})
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
func UnCordon(client corev1Typed.NodeInterface, nodeName string) error {
	return setUnSchedulableValue(client, nodeName, false)
}

// setUnSchedulableValue cordones/uncordons a node by a given node name.
func setUnSchedulableValue(client corev1Typed.NodeInterface, nodeName string, unSchedulable bool) error {
	node, err := client.Get(context.TODO(), nodeName, metav1.GetOptions{})
	if err != nil {
		return err
	}

	node.Spec.Unschedulable = unSchedulable
	_, err = client.Update(context.TODO(), node, metav1.UpdateOptions{})

	if err != nil {
		return fmt.Errorf("failed to update node %s - error: %w", node.Name, err)
	}

	return nil
}

func IsNodeMaster(node *corev1.Node, client corev1Typed.NodeInterface) (bool, error) {
	masterLabels := []string{"node-role.kubernetes.io/master", "node-role.kubernetes.io/control-plane"}

	for _, label := range masterLabels {
		if _, exists := node.Labels[label]; exists {
			return true, nil
		}
	}

	return false, nil
}

func EnsureAllNodesAreLabeled(label string) error {
	return ensureAllNodesAreLabeled(globalhelper.GetAPIClient().K8sClient.CoreV1(), label)
}

// EnsureAllNodesAreLabeled ensures that all nodes are labeled with the given label.
func ensureAllNodesAreLabeled(client corev1Typed.CoreV1Interface, label string) error {
	// Get all nodes in the cluster
	nodes, err := client.Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return fmt.Errorf("failed to get nodes: %w", err)
	}

	for _, node := range nodes.Items {
		if _, exists := node.Labels[label]; !exists {
			err := LabelNode(client, &node, label, "")

			if err != nil {
				return err
			}
		}
	}

	return nil
}

// LabelNode labels a node by a given node name.
func LabelNode(client corev1Typed.CoreV1Interface, node *corev1.Node, label, value string) error {
	node.Labels[label] = value

	// Set the label
	_, err := client.Nodes().Update(context.TODO(), node, metav1.UpdateOptions{})

	return err
}

// EnableMasterScheduling enables/disables master nodes scheduling.
func EnableMasterScheduling(client corev1Typed.NodeInterface, scheduleable bool) error {
	// Get all nodes in the cluster
	nodes, err := client.List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return fmt.Errorf("failed to get nodes: %w", err)
	}

	// Loop through the nodes and modify the taints
	for _, node := range nodes.Items {
		isMaster, err := IsNodeMaster(&node, client)
		if err != nil {
			return fmt.Errorf("failed to get node %s: %w", node.Name, err)
		}

		if isMaster {
			if scheduleable {
				err = removeControlPlaneTaint(client, &node)
				if err != nil {
					return fmt.Errorf("failed to set node %s schedulable value: %w", node.Name, err)
				}
			} else {
				err = addControlPlaneTaint(client, &node)
				if err != nil {
					return fmt.Errorf("failed to set node %s schedulable value: %w", node.Name, err)
				}
			}
		}
	}

	return nil
}

func addControlPlaneTaint(client corev1Typed.NodeInterface, node *corev1.Node) error {
	// add the control-plane:NoSchedule taint to the master
	// check if the tainted already exists to avoid duplicate key error
	for _, taint := range node.Spec.Taints {
		if taint.Key == masterTaintKey || taint.Key == controlPlaneTaintKey {
			return nil
		}
	}
	node.Spec.Taints = append(node.Spec.Taints, corev1.Taint{
		Key:    controlPlaneTaintKey,
		Effect: corev1.TaintEffectNoSchedule,
	})
	_, err := client.Update(context.TODO(), node, metav1.UpdateOptions{})

	if err != nil {
		return fmt.Errorf("failed to update node %s - error: %w", node.Name, err)
	}

	return nil
}

func removeControlPlaneTaint(client corev1Typed.NodeInterface, node *corev1.Node) error {
	// remove the control-plane:NoSchedule taint from the master
	for i, taint := range node.Spec.Taints {
		if taint.Key == masterTaintKey || taint.Key == controlPlaneTaintKey {
			node.Spec.Taints = append(node.Spec.Taints[:i], node.Spec.Taints[i+1:]...)
		}
	}

	_, err := client.Update(context.TODO(), node, metav1.UpdateOptions{})

	if err != nil {
		return fmt.Errorf("failed to update node %s - error: %w", node.Name, err)
	}

	return nil
}
