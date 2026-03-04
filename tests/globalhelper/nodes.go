package globalhelper

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
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
	filteredTaints := make([]corev1.Taint, 0, len(node.Spec.Taints))

	for _, taint := range node.Spec.Taints {
		if taint.Key != masterTaintKey && taint.Key != controlPlaneTaintKey {
			filteredTaints = append(filteredTaints, taint)
		}
	}

	node.Spec.Taints = filteredTaints

	_, err := client.Nodes().Update(context.TODO(), node, metav1.UpdateOptions{})

	if err != nil {
		return fmt.Errorf("failed to update node %s - error: %w", node.Name, err)
	}

	return nil
}

func NodesHaveHugePagesEnabled(resourceName string) bool {
	// check if the node has hugepages enabled
	nodes, err := GetAPIClient().Nodes().List(context.TODO(), metav1.ListOptions{
		LabelSelector: "node-role.kubernetes.io/worker-cnf",
	})

	if err != nil {
		return false
	}

	for _, node := range nodes.Items {
		if NodeHasHugePagesEnabled(&node, resourceName) {
			return true
		}
	}

	return false
}

func NodeHasHugePagesEnabled(node *corev1.Node, resourceName string) bool {
	// check if the node has hugepages enabled
	hugepagesEnabled := false
	resourceNameStr := "hugepages-" + resourceName

	if node.Status.Capacity != nil {
		if _, ok := node.Status.Capacity[corev1.ResourceName(resourceNameStr)]; ok {
			if node.Status.Capacity[corev1.ResourceName(resourceNameStr)] != resource.MustParse("0") {
				hugepagesEnabled = true
			}
		}
	}

	return hugepagesEnabled
}

func GetNumberOfNodes(client corev1Typed.CoreV1Interface) int {
	return getNumberOfNodes(client)
}

func getNumberOfNodes(client corev1Typed.CoreV1Interface) int {
	nodes, err := client.Nodes().List(context.TODO(), metav1.ListOptions{})

	if err != nil {
		return 0
	}

	return len(nodes.Items)
}

func IsClusterOvercommitted() (bool, error) {
	return isClusterOvercommitted(GetAPIClient().Nodes())
}

func isClusterOvercommitted(client corev1Typed.NodeInterface) (bool, error) {
	nodes, err := client.List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return false, err
	}

	for _, node := range nodes.Items {
		if isNodeOvercommitted(&node) {
			return true, nil
		}
	}

	return false, nil
}

func isNodeOvercommitted(node *corev1.Node) bool {
	// check if the node is overcommitted
	if node.Status.Allocatable.Memory().Cmp(resource.MustParse("0")) > 0 {
		return true
	}

	if node.Status.Allocatable.Cpu().Cmp(resource.MustParse("0")) > 0 {
		return true
	}

	return false
}
