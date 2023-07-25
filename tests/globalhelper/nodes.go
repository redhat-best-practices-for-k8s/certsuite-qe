package globalhelper

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	controlPlaneTaintKey = "node-role.kubernetes.io/control-plane"
	masterTaintKey       = "node-role.kubernetes.io/master"
)

// EnableMasterScheduling enables/disables master nodes scheduling.
func EnableMasterScheduling(scheduleable bool) error {
	// Get all nodes in the cluster
	nodes, err := GetAPIClient().CoreV1Interface.Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return fmt.Errorf("failed to get nodes: %w", err)
	}

	// Loop through the nodes and modify the taints
	for _, node := range nodes.Items {
		if isMasterNode(&node) {
			if scheduleable {
				err = removeControlPlaneTaint(&node)
				if err != nil {
					return fmt.Errorf("failed to set node %s schedulable value: %w", node.Name, err)
				}
			} else {
				err = addControlPlaneTaint(&node)
				if err != nil {
					return fmt.Errorf("failed to set node %s schedulable value: %w", node.Name, err)
				}
			}
		}
	}

	return nil
}

func addControlPlaneTaint(node *corev1.Node) error {
	// add the control-plane:NoSchedule taint to the master
	node.Spec.Taints = append(node.Spec.Taints, corev1.Taint{
		Key:    controlPlaneTaintKey,
		Effect: corev1.TaintEffectNoSchedule,
	})

	_, err := GetAPIClient().CoreV1Interface.Nodes().Update(context.TODO(), node, metav1.UpdateOptions{})

	if err != nil {
		return fmt.Errorf("failed to update node %s - error: %w", node.Name, err)
	}

	return nil
}

func removeControlPlaneTaint(node *corev1.Node) error {
	// remove the control-plane:NoSchedule taint from the master
	var updatedTaints []corev1.Taint

	for _, taint := range node.Spec.Taints {
		switch taint.Key {
		case masterTaintKey, controlPlaneTaintKey:
			// Skip the taints that need to be removed
			continue
		default:
			updatedTaints = append(updatedTaints, taint)
		}
	}

	node.Spec.Taints = updatedTaints

	_, err := GetAPIClient().CoreV1Interface.Nodes().Update(context.TODO(), node, metav1.UpdateOptions{})

	if err != nil {
		return fmt.Errorf("failed to update node %s - error: %w", node.Name, err)
	}

	return nil
}

func isMasterNode(node *corev1.Node) bool {
	masterLabels := []string{masterTaintKey, controlPlaneTaintKey}
	for _, label := range masterLabels {
		if _, exists := node.Labels[label]; exists {
			return true
		}
	}

	return false
}
