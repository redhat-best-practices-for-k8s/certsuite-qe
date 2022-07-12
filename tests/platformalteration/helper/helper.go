package helper

import (
	"context"
	"time"

	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/client"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/nodes"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
)

const WaitingTime = 5 * time.Minute

func WaitForSpecificNodeCondition(clients *client.ClientSet, timeout, interval time.Duration, nodeName string,
	ready bool) error {
	return wait.PollImmediate(interval, timeout, func() (bool, error) {
		nodesList, err := clients.Nodes().List(context.Background(), metav1.ListOptions{})
		if err != nil {
			return false, err
		}

		// Verify the node condition
		for _, node := range nodesList.Items {
			if node.Name == nodeName && nodes.IsNodeInCondition(&node, corev1.NodeReady) == ready {
				return true, nil
			}
		}

		return false, nil
	})
}
