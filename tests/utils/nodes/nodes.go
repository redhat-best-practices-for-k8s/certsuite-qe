package nodes

import (
	"context"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/client"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"log"
	"time"
)

// WaitForNodesReady waits for all nodes become ready
func WaitForNodesReady(cs *client.ClientSet, timeout, interval time.Duration) error {
	return wait.PollImmediate(interval, timeout, func() (bool, error) {
		nodes_, err := cs.Nodes().List(context.Background(), metav1.ListOptions{})
		if err != nil {
			return false, nil
		}
		for _, node := range nodes_.Items {
			if !IsNodeInCondition(&node, corev1.NodeReady) {
				return false, nil
			}
		}
		log.Println("All nodes are Ready")
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
