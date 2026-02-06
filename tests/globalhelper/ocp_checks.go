package globalhelper

import (
	"context"
	"encoding/json"
	"fmt"

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

// HasPerformanceProfiles checks if there are any PerformanceProfile resources configured.
// The PerformanceProfile CRD may exist but without any actual resources, exclusive CPU pools
// won't work because the CPU Manager policy will be 'none' instead of 'static'.
func HasPerformanceProfiles() (bool, error) {
	// Use REST client to query PerformanceProfiles (cluster-scoped)
	data, err := GetAPIClient().CoreV1Interface.RESTClient().
		Get().
		AbsPath("/apis/performance.openshift.io/v2/performanceprofiles").
		DoRaw(context.TODO())

	if err != nil {
		return false, fmt.Errorf("failed to list PerformanceProfiles: %w", err)
	}

	// Parse the response to check if there are items
	var result struct {
		Items []interface{} `json:"items"`
	}
	if err := json.Unmarshal(data, &result); err != nil {
		return false, fmt.Errorf("failed to parse PerformanceProfiles response: %w", err)
	}

	return len(result.Items) > 0, nil
}

// CPUManagerState represents the CPU manager state from kubelet.
type CPUManagerState struct {
	PolicyName string `json:"policyName"`
}

// HasStaticCPUManagerPolicy checks if nodes have CPU Manager configured with 'static' policy.
// This is required for exclusive CPU pools to work properly.
func HasStaticCPUManagerPolicy() (bool, error) {
	// Get a worker node
	nodes, err := GetAPIClient().Nodes().List(context.TODO(), metav1.ListOptions{
		LabelSelector: "node-role.kubernetes.io/worker",
	})
	if err != nil {
		return false, fmt.Errorf("failed to list worker nodes: %w", err)
	}

	if len(nodes.Items) == 0 {
		return false, fmt.Errorf("no worker nodes found")
	}

	// We cannot directly check CPU manager state from the test harness without debug pod access.
	// However, if PerformanceProfiles exist, the CPU Manager policy should be 'static'.
	// For now, we rely on HasPerformanceProfiles as a proxy for static CPU manager policy.
	return HasPerformanceProfiles()
}

// IsClusterConfiguredForExclusiveCPUs checks if the cluster has the necessary configuration
// to support exclusive CPU pool tests. This includes:
// 1. Worker nodes exist
// 2. PerformanceProfile resources exist (not just the CRD)
// 3. Implicitly, CPU Manager policy is 'static' (configured via PerformanceProfile).
func IsClusterConfiguredForExclusiveCPUs() (bool, string, error) {
	// Check for worker nodes
	if !HasWorkerNodes() {
		return false, "cluster has no worker nodes", nil
	}

	// Check for PerformanceProfile resources (not just CRD)
	hasProfiles, err := HasPerformanceProfiles()
	if err != nil {
		// If we can't query PerformanceProfiles, the CRD likely doesn't exist
		return false, "cannot query PerformanceProfiles: " + err.Error(), nil
	}

	if !hasProfiles {
		return false, "no PerformanceProfile resources configured (CPU Manager policy is likely 'none')", nil
	}

	return true, "", nil
}
