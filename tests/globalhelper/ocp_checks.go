package globalhelper

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// IsMCOHealthy checks if MachineConfigOperator is running and accessible.
// This verifies that we can list MachineConfigs and MachineConfigPools,
// which is required for platform-alteration tests that depend on MCO.
func IsMCOHealthy() (bool, error) {
	// Check we can list MachineConfigs
	_, err := GetAPIClient().MachineConfigs().List(context.TODO(), metav1.ListOptions{Limit: 1})
	if err != nil {
		return false, err
	}

	// Check we can list MachineConfigPools
	mcpList, err := GetAPIClient().MachineConfigPools().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return false, err
	}

	// Verify at least one MachineConfigPool exists
	if len(mcpList.Items) == 0 {
		return false, nil
	}

	return true, nil
}

// HasWorkerNodes checks if cluster has worker nodes (not just control-plane).
// This is required for tests that need pods to run on worker nodes.
func HasWorkerNodes() bool {
	nodes, err := GetAPIClient().Nodes().List(context.TODO(), metav1.ListOptions{
		LabelSelector: "node-role.kubernetes.io/worker",
	})
	if err != nil {
		return false
	}

	return len(nodes.Items) > 0
}
