package globalhelper

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	corev1Typed "k8s.io/client-go/kubernetes/typed/core/v1"
)

const (
	controlPlaneTaintKey = "node-role.kubernetes.io/control-plane"
	masterTaintKey       = "node-role.kubernetes.io/master"
)

func addControlPlaneTaint(client corev1Typed.CoreV1Interface, node *corev1.Node) error {
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

	_, err := client.Nodes().Update(context.TODO(), node, metav1.UpdateOptions{})

	if err != nil {
		return fmt.Errorf("failed to update node %s - error: %w", node.Name, err)
	}

	return nil
}

func removeControlPlaneTaint(client corev1Typed.CoreV1Interface, node *corev1.Node) error {
	// remove the control-plane:NoSchedule taint from the master
	// remove the control-plane:NoSchedule taint from the master
	for i, taint := range node.Spec.Taints {
		if taint.Key == masterTaintKey || taint.Key == controlPlaneTaintKey {
			node.Spec.Taints = append(node.Spec.Taints[:i], node.Spec.Taints[i+1:]...)
		}
	}

	_, err := client.Nodes().Update(context.TODO(), node, metav1.UpdateOptions{})

	if err != nil {
		return fmt.Errorf("failed to update node %s - error: %w", node.Name, err)
	}

	return nil
}
