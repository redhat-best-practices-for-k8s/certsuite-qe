package cluster

import (
	"context"
	"fmt"

	"github.com/golang/glog"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	testclient "github.com/test-network-function/cnfcert-tests-verification/tests/utils/client"
)

func IsClusterStable(clients *testclient.ClientSet) (bool, error) {
	nodes, err := clients.Nodes().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return false, err
	}

	for _, node := range nodes.Items {
		if node.Spec.Unschedulable {
			glog.V(5).Info(fmt.Sprintf("node %s is in unschedulable state", node.Name))
			return false, nil
		}

	}
	return true, nil
}
