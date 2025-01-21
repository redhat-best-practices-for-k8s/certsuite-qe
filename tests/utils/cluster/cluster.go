package cluster

import (
	"context"
	"fmt"

	"github.com/golang/glog"

	egiNodes "github.com/openshift-kni/eco-goinfra/pkg/nodes"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/globalhelper"
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

func IsSNO() (bool, error) {
	nodes, err := egiNodes.List(globalhelper.GetEcoGoinfraClient(), metav1.ListOptions{})
	if err != nil {
		return false, fmt.Errorf("failed to list nodes with eco-goinfra client: %w", err)
	}

	return len(nodes) == 1, nil
}

func IsCompact() (bool, error) {
	nodes, err := egiNodes.List(globalhelper.GetEcoGoinfraClient(), metav1.ListOptions{})
	if err != nil {
		return false, fmt.Errorf("failed to list nodes with eco-goinfra client: %w", err)
	}

	return len(nodes) == 3, nil
}
