package cluster

import (
	"context"
	"fmt"

	"github.com/golang/glog"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	testclient "github.com/test-network-function/cnfcert-tests-verification/tests/utils/client"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/nodes"
)

func IsClusterStable(clients *testclient.ClientSet) (bool, error) {
	nodeInCluster, err := clients.Nodes().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return false, err
	}

	for _, node := range nodeInCluster.Items {
		if node.Spec.Unschedulable {
			fmt.Println("node", node.Name, "is cordoned, uncordoning it")
			err := nodes.UnCordon(clients, node.Name)

			if err != nil {
				return false, err
			}

			glog.V(5).Info(fmt.Sprintf("node %s is in unschedulable state", node.Name))
		}
	}

	return true, nil
}
