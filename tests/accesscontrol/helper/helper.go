package helper

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	tsparams "github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/accesscontrol/parameters"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/globalhelper"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/utils/container"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/utils/deployment"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/utils/service"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/utils/statefulset"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"sigs.k8s.io/controller-runtime/pkg/client"

	machineconfigv1 "github.com/openshift/api/machineconfiguration/v1"
	mcv1 "github.com/openshift/machine-config-operator/pkg/apis/machineconfiguration.openshift.io/v1"

	egiMco "github.com/openshift-kni/eco-goinfra/pkg/mco"
	egiNodes "github.com/openshift-kni/eco-goinfra/pkg/nodes"
	klog "k8s.io/klog/v2"
)

func DefineDeploymentWithImage(replica int32, containers int, name, namespace, image string) (*appsv1.Deployment, error) {
	if containers < 1 {
		return nil, errors.New("invalid number of containers")
	}

	deploymentStruct := deployment.DefineDeployment(name, namespace,
		image, tsparams.TestDeploymentLabels)

	globalhelper.AppendContainersToDeployment(deploymentStruct, containers-1, image)
	deployment.RedefineWithReplicaNumber(deploymentStruct, replica)

	return deploymentStruct, nil
}

func DefineDeployment(replica int32, containers int, name, namespace string) (*appsv1.Deployment, error) {
	if containers < 1 {
		return nil, errors.New("invalid number of containers")
	}

	deploymentStruct := deployment.DefineDeployment(name, namespace,
		tsparams.SampleWorkloadImage, tsparams.TestDeploymentLabels)

	globalhelper.AppendContainersToDeployment(deploymentStruct, containers-1, tsparams.SampleWorkloadImage)
	deployment.RedefineWithReplicaNumber(deploymentStruct, replica)

	return deploymentStruct, nil
}

func DefineStatefulSet(replica int32, containers int, name, namespace string) (*appsv1.StatefulSet, error) {
	if containers < 1 {
		return nil, errors.New("invalid number of containers")
	}

	sts := statefulset.DefineStatefulSet(name, namespace,
		tsparams.SampleWorkloadImage, tsparams.TestStatefulSetLabels)

	globalhelper.AppendContainersToStatefulSet(sts, containers-1, tsparams.SampleWorkloadImage)
	statefulset.RedefineWithReplicaNumber(sts, replica)

	return sts, nil
}

func DefineDeploymentWithClusterRoleBindingWithServiceAccount(replica int32,
	containers int, name, namespace, serviceAccountName string) (*appsv1.Deployment, error) {
	deploymentStruct := deployment.DefineDeployment(name, namespace,
		tsparams.SampleWorkloadImage, tsparams.TestDeploymentLabels)

	globalhelper.AppendContainersToDeployment(deploymentStruct, containers-1, tsparams.SampleWorkloadImage)
	deployment.RedefineWithReplicaNumber(deploymentStruct, replica)
	deployment.AppendServiceAccount(deploymentStruct, serviceAccountName)

	return deploymentStruct, nil
}

func DefineDeploymentWithNamespace(replica int32, containers int, name string, namespace string) (*appsv1.Deployment, error) {
	if containers < 1 {
		return nil, errors.New("invalid number of containers")
	}

	deploymentStruct := deployment.DefineDeployment(name, namespace,
		tsparams.SampleWorkloadImage, tsparams.TestDeploymentLabels)

	globalhelper.AppendContainersToDeployment(deploymentStruct, containers-1, tsparams.SampleWorkloadImage)
	deployment.RedefineWithReplicaNumber(deploymentStruct, replica)

	return deploymentStruct, nil
}

func DefineDeploymentWithContainerPorts(name, namespace string, replicaNumber int32,
	ports []corev1.ContainerPort) (*appsv1.Deployment, error) {
	if len(ports) < 1 {
		return nil, errors.New("invalid number of containers")
	}

	deploymentStruct := deployment.DefineDeployment(name, namespace,
		tsparams.SampleWorkloadImage, tsparams.TestDeploymentLabels)

	globalhelper.AppendContainersToDeployment(deploymentStruct, len(ports)-1, tsparams.SampleWorkloadImage)
	deployment.RedefineWithReplicaNumber(deploymentStruct, replicaNumber)

	portSpecs := container.CreateContainerSpecsFromContainerPorts(ports,
		tsparams.SampleWorkloadImage, "test")

	deployment.RedefineWithContainerSpecs(deploymentStruct, portSpecs)

	return deploymentStruct, nil
}

func SetServiceAccountAutomountServiceAccountToken(namespace, saname, value string) error {
	var boolVal bool

	serviceacct, err := globalhelper.GetAPIClient().ServiceAccounts(namespace).
		Get(context.TODO(), saname, metav1.GetOptions{})

	if err != nil {
		return fmt.Errorf("error getting service account: %w", err)
	}

	switch value {
	case "true":
		boolVal = true
		serviceacct.AutomountServiceAccountToken = &boolVal

	case "false":
		boolVal = false
		serviceacct.AutomountServiceAccountToken = &boolVal

	case "nil":
		serviceacct.AutomountServiceAccountToken = nil

	default:
		return fmt.Errorf("invalid value for token value")
	}

	_, err = globalhelper.GetAPIClient().ServiceAccounts(namespace).
		Update(context.TODO(), serviceacct, metav1.UpdateOptions{})

	return err
}

// DefineAndCreateServiceOnCluster defines service resource and creates it on cluster.
func DefineAndCreateServiceOnCluster(name, namespace string, port int32, targetPort int32, withNodePort bool,
	ipFams []corev1.IPFamily, ipFamPolicy string) error {
	var testService *corev1.Service

	if ipFamPolicy == "" {
		testService = service.DefineService(
			name,
			namespace,
			port,
			targetPort,
			corev1.ProtocolTCP,
			tsparams.TestDeploymentLabels,
			ipFams,
			nil)
	} else {
		ipPolicy := corev1.IPFamilyPolicy(ipFamPolicy)

		testService = service.DefineService(
			name,
			namespace,
			port,
			targetPort,
			corev1.ProtocolTCP,
			tsparams.TestDeploymentLabels,
			ipFams,
			&ipPolicy)
	}

	if withNodePort {
		err := service.RedefineWithNodePort(testService)
		if err != nil {
			return err
		}
	}

	_, err := globalhelper.GetAPIClient().Services(namespace).Create(
		context.TODO(),
		testService, metav1.CreateOptions{})

	if k8serrors.IsAlreadyExists(err) {
		return nil
	} else if err != nil {
		return fmt.Errorf("failed to create service on cluster: %w", err)
	}

	return nil
}

// Finds the MachineConfigPools that manage worker nodes of the cluster. Returns a map of MCP names
// to slices of worker names.
func GetWorkersMCPs() (mcpNodes map[string][]string, err error) {
	// Get the eco-goinfra client.
	egiClient := globalhelper.GetEcoGoinfraClient()
	if egiClient == nil {
		return nil, fmt.Errorf("eco-goinfra client is not initialized")
	}

	nodeBuilders, err := egiNodes.List(egiClient, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get nodes: %w", err)
	}

	mcpBuilders, err := egiMco.ListMCP(egiClient, client.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list MachineConfigPools: %w", err)
	}

	mcps := map[string][]string{}

	// For each mcp, add entry to the herlpe mcps map if it has worker nodes.
	for i := range mcpBuilders {
		mcpNodes := []string{}
		mcp := mcpBuilders[i].Object

		if mcp.Status.MachineCount == 0 {
			continue
		}

		labelSelector, err := metav1.LabelSelectorAsSelector(mcp.Spec.NodeSelector)
		if err != nil {
			return nil, fmt.Errorf("failed to get labelselector from mcp %s: %w", mcp.Name, err)
		}

		for i := range nodeBuilders {
			node := nodeBuilders[i].Object

			nodeLabelsSet := labels.Set(node.Labels)
			if !labelSelector.Matches(nodeLabelsSet) {
				continue
			}

			if _, exist := node.Labels["node-role.kubernetes.io/worker"]; !exist {
				break
			}
			mcpNodes = append(mcpNodes, node.Name)
		}

		// Continue if this MCP doesn't manage any worker node
		if len(mcpNodes) == 0 {
			klog.V(5).Infof("MCP=%s doesn't manage any worker node", mcp.Name)

			continue
		}

		mcps[mcp.Name] = mcpNodes

		klog.V(5).Infof("MCP=%s is candidate, num workers=%d", mcp.Name, len(mcpNodes))
	}

	return mcps, nil
}

// Gets the label set that a MachineConfig needs so it's used by a MachineConfigPool.
//
// The MachineConfigPool watches all the MachineConfig objects in the cluster but it only uses the ones that matches
// the selector in MachineconfigPool.Spec.MachineConfigSelector to create the final rendered-xxxx MachineConfig.
//
// The code uses the MachineConfigSelector on each MachineConfig found and returns the labels of the MachineConfig that
// matches that selector.
func GetMachineConfigTargetLabels(mcpName string) (map[string]string, error) {
	// Get the eco-goinfra client.
	egiClient := globalhelper.GetEcoGoinfraClient()
	if egiClient == nil {
		return nil, fmt.Errorf("eco-goinfra client is not initialized")
	}

	// Create a labelSelector from MCP's Spec.MachineConfigSelector
	mcpBuilder := egiMco.NewMCPBuilder(egiClient, mcpName)
	if mcpBuilder == nil {
		return nil, fmt.Errorf("failed to get MachineConfigPool builder")
	}

	if !mcpBuilder.Exists() {
		return nil, fmt.Errorf("mcp %s not found", mcpName)
	}

	labelSelector, err := metav1.LabelSelectorAsSelector(mcpBuilder.Object.Spec.MachineConfigSelector)
	if err != nil {
		return nil, fmt.Errorf("failed to create label selector: %w", err)
	}

	// Now, list the MachineConfigs and get the labels of the first one that would be matched by the MCP's selector.
	mcBuilders, err := egiMco.ListMC(egiClient, client.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed getting MachineConfigs: %w", err)
	}

	for i := range mcBuilders {
		mc := mcBuilders[i].Object

		mcLabelsSet := labels.Set(mc.Labels)
		if labelSelector.Matches(mcLabelsSet) {
			return mc.Labels, nil
		}
	}

	return nil, fmt.Errorf("no machineconfig matches the label selector %v", mcpBuilder.Object.Spec.MachineConfigSelector)
}

// Checks whether a node is running using a realtime kernel type.
func HasNodeRtKernel(nodeName string) (bool, error) {
	node, err := globalhelper.GetAPIClient().Nodes().Get(context.TODO(), nodeName, metav1.GetOptions{})
	if err != nil {
		return false, fmt.Errorf("failed to get node %s: %w", nodeName, err)
	}

	return strings.Contains(node.Status.NodeInfo.KernelVersion, "+rt"), nil
}

// Deploys a Machineconfig to switch the kernel to realtime type on all th enodes of a MachineConfigPool.
//
// Warning: This is a disruptive operation that will reboot all the worker nodes one by one.
func DeployRTKernelMachineConfig(mcpName string, mcName string, mcpLabels map[string]string, timeout time.Duration) error {
	// Get the eco-goinfra client.
	egiClient := globalhelper.GetEcoGoinfraClient()
	if egiClient == nil {
		return fmt.Errorf("eco-goinfra client is not initialized")
	}

	// Create custom MachineConfig with kerneltype realtime.
	mcBuilder := egiMco.NewMCBuilder(egiClient, mcName).WithKernelType("realtime")
	if mcBuilder == nil {
		return fmt.Errorf("failed to get MachineConfig builder")
	}

	// Add labels to MC so the MCP can match it.
	for k, v := range mcpLabels {
		mcBuilder.WithLabel(k, v)
	}

	klog.V(5).Infof("Creating MC=%s in the cluster.", mcName)
	// Apply MachineConfig
	_, err := mcBuilder.Create()
	if err != nil {
		return fmt.Errorf("failed to create MachineConfig: %w", err)
	}

	// MachineConfigPool
	mcpBuilder := egiMco.NewMCPBuilder(egiClient, mcpName)
	if mcpBuilder == nil {
		return fmt.Errorf("failed to get MachineConfigPool builder")
	}

	// Wait for MCP to start updating: condition Updating=True
	klog.V(5).Infof("MCP=%s: waiting %s for condition Updating to be True", mcpName, tsparams.McpStartTimeout)

	err = mcpBuilder.WaitToBeInCondition(machineconfigv1.MachineConfigPoolConditionType(
		mcv1.MachineConfigPoolUpdating), corev1.ConditionTrue, tsparams.McpStartTimeout)
	if err != nil {
		return fmt.Errorf("failed waiting for MachineConfigPool to set condition Updating=True: %w", err)
	}

	// Wait for MCP to finish the update using the dynamic timeout as deadline.
	klog.V(5).Infof("MCP=%s: waiting %s for condition Updating to be False", mcpName, timeout)

	err = mcpBuilder.WaitToBeInCondition(machineconfigv1.MachineConfigPoolConditionType(
		mcv1.MachineConfigPoolUpdating), corev1.ConditionFalse, timeout)
	if err != nil {
		return fmt.Errorf("failed waiting for MachineConfigPool to set condition Updating=False: %w", err)
	}

	// Make sure everything went ok
	if mcpBuilder.IsInCondition(machineconfigv1.MachineConfigPoolConditionType(mcv1.MachineConfigPoolDegraded)) {
		return fmt.Errorf("machineconfigpool appears as Degraded after machineconfig was applied")
	}

	return nil
}

// Removes the Machineconfig that switches the kernel to realtime type on all the nodes of a MachineConfigPool.
//
// Warning: This is a disruptive operation that will reboot all the worker nodes one by one.
func RemoveRTKernelMachineConfig(mcpName string, mcName string, timeout time.Duration) error {
	// Get the eco-goinfra client.
	egiClient := globalhelper.GetEcoGoinfraClient()
	if egiClient == nil {
		return fmt.Errorf("eco-goinfra client is not initialized")
	}

	// Delete the MachineConfig
	err := egiMco.NewMCBuilder(egiClient, mcName).Delete()
	if err != nil {
		return fmt.Errorf("failed to delete machineconfig %s: %w", mcName, err)
	}

	// Wait for MCP to start updating.
	mcpBuilder := egiMco.NewMCPBuilder(egiClient, mcpName)
	if mcpBuilder == nil {
		return fmt.Errorf("failed to get MachineConfigPool builder")
	}

	// Wait for MCP to start updating: condition Updating=True
	klog.V(5).Infof("MCP=%s: waiting %s for condition Updating to be True", mcpName, tsparams.McpStartTimeout)

	err = mcpBuilder.WaitToBeInCondition(machineconfigv1.MachineConfigPoolConditionType(
		mcv1.MachineConfigPoolUpdating), corev1.ConditionTrue, tsparams.McpStartTimeout)
	if err != nil {
		return fmt.Errorf("failed waiting for MachineConfigPool to set condition Updating=True: %w", err)
	}

	// Wait for MCP to finish the update using the dynamic timeout as deadline.
	klog.V(5).Infof("MCP=%s: waiting %s for condition Updating to be False", mcpName, timeout)

	err = mcpBuilder.WaitToBeInCondition(machineconfigv1.MachineConfigPoolConditionType(
		mcv1.MachineConfigPoolUpdating), corev1.ConditionFalse, timeout)
	if err != nil {
		return fmt.Errorf("failed waiting for MachineConfigPool to set condition Updating=False: %w", err)
	}

	// Make sure everything went ok
	if mcpBuilder.IsInCondition(machineconfigv1.MachineConfigPoolConditionType(mcv1.MachineConfigPoolDegraded)) {
		return fmt.Errorf("machineconfigpool appears as Degraded after machineconfig was applied")
	}

	return nil
}

// Returns true if MachineConfigPool's Condition "Updating" (MachineConfigPoolUpdated) exists and is True.
// Returns error if conditions slice is empty or condition Updating is not found.
func IsMCPConditionUpdatingTrue(mcpConditions []mcv1.MachineConfigPoolCondition) (bool, error) {
	// No conditions available yet.
	if len(mcpConditions) == 0 {
		return false, fmt.Errorf("conditions are empty/nil")
	}

	// Search for condition type "Update".
	for _, cond := range mcpConditions {
		if cond.Type == mcv1.MachineConfigPoolUpdating {
			return cond.Status == corev1.ConditionTrue, nil
		}
	}

	// Conditions exist but Updated condition not available yet.
	return false, fmt.Errorf("condition %s not found yet", mcv1.MachineConfigPoolUpdating)
}
