package cluster

import (
	"context"
	"fmt"

	klog "k8s.io/klog/v2"

	nodesutils "github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/utils/nodes"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	corev1Typed "k8s.io/client-go/kubernetes/typed/core/v1"
)

// IsClusterStable tests if cluster is stable.
func IsClusterStable(client corev1Typed.NodeInterface) (bool, error) {
	nodes, err := client.List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return false, err
	}

	for _, node := range nodes.Items {
		if node.Spec.Unschedulable {
			klog.V(5).Info(fmt.Sprintf("node %s is in unschedulable state, trying to uncordon it", node.Name))
			err := nodesutils.UnCordon(client, node.Name)

			if err != nil {
				return false, err
			}

			updatedNode, err := client.Get(context.TODO(), node.Name, metav1.GetOptions{})
			if err != nil {
				return false, err
			}

			if updatedNode.Spec.Unschedulable {
				return false, nil
			}

			klog.V(5).Info(fmt.Sprintf("node %s is in schedulable state", node.Name))
		}
	}

	return true, nil
}
