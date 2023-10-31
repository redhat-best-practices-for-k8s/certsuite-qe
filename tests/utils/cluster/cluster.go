package cluster

import (
	"context"
	"fmt"

	"github.com/golang/glog"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	nodesutils "github.com/test-network-function/cnfcert-tests-verification/tests/utils/nodes"
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
			glog.V(5).Info(fmt.Sprintf("node %s is in unschedulable state, trying to uncordon it", node.Name))
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

			glog.V(5).Info(fmt.Sprintf("node %s is in schedulable state", node.Name))
		}
	}

	return true, nil
}
