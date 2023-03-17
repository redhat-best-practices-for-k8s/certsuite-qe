package helper

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/nad"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/nodes"

	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/daemonset"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/networkpolicy"

	tsparams "github.com/test-network-function/cnfcert-tests-verification/tests/networking/parameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/deployment"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/service"
	corev1 "k8s.io/api/core/v1"

	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// DefineAndCreateDeploymentOnCluster defines deployment resource and creates it on cluster.
func DefineAndCreateDeploymentOnCluster(replicaNumber int32) error {
	deploymentUnderTest := defineDeploymentBasedOnArgs(tsparams.TestDeploymentAName, tsparams.TestNetworkingNameSpace,
		replicaNumber, false, nil, nil)

	return globalhelper.CreateAndWaitUntilDeploymentIsReady(deploymentUnderTest, tsparams.WaitingTime)
}

// DefineAndCreateDeploymentWithMultusOnCluster defines deployment resource and creates it on cluster.
func DefineAndCreateDeploymentWithMultusOnCluster(name string, nadNames []string, replicaNumber int32) error {
	deploymentUnderTest := defineDeploymentBasedOnArgs(name, tsparams.TestNetworkingNameSpace,
		replicaNumber, true, nadNames, nil)

	return globalhelper.CreateAndWaitUntilDeploymentIsReady(deploymentUnderTest, tsparams.WaitingTime)
}

// DefineAndCreateDeploymentWithMultusAndSkipLabelOnCluster defines deployment resource and creates it on cluster.
func DefineAndCreateDeploymentWithMultusAndSkipLabelOnCluster(
	name string, nadNames []string, replicaNumber int32) error {
	deploymentUnderTest := defineDeploymentBasedOnArgs(
		name, tsparams.TestNetworkingNameSpace,
		replicaNumber,
		false,
		nadNames,
		tsparams.NetworkingTestMultusSkipLabel)

	return globalhelper.CreateAndWaitUntilDeploymentIsReady(deploymentUnderTest, tsparams.WaitingTime)
}

// DefineAndCreatePrivilegedDeploymentOnCluster defines deployment resource and creates it on cluster.
func DefineAndCreatePrivilegedDeploymentOnCluster(replicaNumber int32) error {
	deploymentUnderTest := defineDeploymentBasedOnArgs(tsparams.TestDeploymentAName, tsparams.TestNetworkingNameSpace,
		replicaNumber, true,
		nil, nil)

	return globalhelper.CreateAndWaitUntilDeploymentIsReady(deploymentUnderTest, tsparams.WaitingTime)
}

// DefineAndCreateDeploymentWithSkippedLabelOnCluster defines deployment resource and creates it on cluster.
func DefineAndCreateDeploymentWithSkippedLabelOnCluster(replicaNumber int32) error {
	deploymentUnderTest := defineDeploymentBasedOnArgs(tsparams.TestDeploymentAName, tsparams.TestNetworkingNameSpace,
		replicaNumber,
		true, nil, tsparams.NetworkingTestSkipLabel)

	err := globalhelper.CreateAndWaitUntilDeploymentIsReady(deploymentUnderTest, tsparams.WaitingTime)
	if err != nil {
		return fmt.Errorf("failed to create deployment: %w", err)
	}

	return nil
}

func DefineAndCreateDeploymentWithNamespace(namespace string, replicaNumber int32) error {
	deploymentUnderTest := defineDeploymentBasedOnArgs(tsparams.TestDeploymentBName, namespace,
		replicaNumber, false, nil, nil)

	return globalhelper.CreateAndWaitUntilDeploymentIsReady(deploymentUnderTest, tsparams.WaitingTime)
}

func DefineAndCreateDeployment(deploymentName string, replicaNumber int32) error {
	deploymentUnderTest := defineDeploymentBasedOnArgs(deploymentName, tsparams.TestNetworkingNameSpace,
		replicaNumber, false, nil, nil)

	return globalhelper.CreateAndWaitUntilDeploymentIsReady(deploymentUnderTest, tsparams.WaitingTime)
}

func DefineAndCreateDeamonsetWithMultusOnCluster(nadName string) error {
	return defineDaemonSetBasedOnArgs(nadName, nil)
}

func DefineAndCreateDeamonsetWithMultusAndSkipLabelOnCluster(nadName string) error {
	return defineDaemonSetBasedOnArgs(nadName, tsparams.NetworkingTestMultusSkipLabel)
}

// DefineAndCreateDeploymentOnCluster defines deployment resource and creates it on cluster.
func DefineAndCreateDeploymentWithContainerPorts(replicaNumber int32, ports []corev1.ContainerPort) error {
	deploymentUnderTest, err := defineDeploymentWithContainers(replicaNumber, len(ports), tsparams.TestDeploymentAName)
	if err != nil {
		return err
	}

	portSpecs := createContainerSpecsFromContainerPorts(ports)

	deployment.RedefineWithContainerSpecs(deploymentUnderTest, portSpecs)

	return globalhelper.CreateAndWaitUntilDeploymentIsReady(deploymentUnderTest, tsparams.WaitingTime)
}

// ExecCmdOnOnePodInNamespace runs command on the first available pod in namespace.
func ExecCmdOnOnePodInNamespace(command []string) error {
	return execCmdOnPodsListInNamespace(command, "first")
}

func ExecCmdOnAllPodInNamespace(command []string) error {
	return execCmdOnPodsListInNamespace(command, "all")
}

// DefineAndCreateServiceOnCluster defines service resource and creates it on cluster.
func DefineAndCreateServiceOnCluster(name string, port int32, targetPort int32, withNodePort bool) error {
	testService := service.DefineService(
		name,
		tsparams.TestNetworkingNameSpace,
		port,
		targetPort,
		corev1.ProtocolTCP,
		tsparams.TestDeploymentLabels)

	if withNodePort {
		var err error

		testService, err = service.RedefineWithNodePort(testService)
		if err != nil {
			return err
		}
	}

	_, err := globalhelper.APIClient.Services(tsparams.TestNetworkingNameSpace).Create(
		context.Background(),
		testService, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("failed to create service on cluster: %w", err)
	}

	return nil
}

func DefineAndCreateNadOnCluster(name string, network string) error {
	nadOneInterface := nad.DefineNad(name, tsparams.TestNetworkingNameSpace)

	if network != "" {
		nadOneInterface = nad.RedefineNadWithWhereaboutsIpam(nadOneInterface, network)
	}

	return globalhelper.APIClient.Create(context.Background(), nadOneInterface)
}

func GetClusterMultusInterfaces() ([]string, error) {
	err := defineAndCreatePrivilegedDaemonset()
	if err != nil {
		return nil, err
	}

	podsList, err := globalhelper.GetListOfPodsInNamespace(tsparams.TestNetworkingNameSpace)
	if err != nil {
		return nil, err
	}

	var nodesInterfacesList [][]string

	for _, runningPod := range podsList.Items {
		isMasterNode, err := nodes.IsNodeMaster(runningPod.Spec.NodeName, globalhelper.APIClient)
		if err != nil {
			return nil, err
		}

		if !isMasterNode {
			nodeInterfaces, err := getInterfacesList(runningPod)
			if err != nil {
				return nil, err
			}
			nodesInterfacesList = append(nodesInterfacesList, nodeInterfaces)
		}
	}

	var lastMatch []string

	for _, nodeInterfaces := range nodesInterfacesList {
		if len(lastMatch) == 0 {
			lastMatch = findListIntersections(nodesInterfacesList[0], nodeInterfaces)
		}

		lastMatch = findListIntersections(lastMatch, nodeInterfaces)
	}

	return lastMatch, nil
}

func DefineAndCreateNetworkPolicy(name, ns string, policyTypes []string, labels map[string]string) error {
	types := networkpolicy.DefinePolicyTypes(policyTypes)
	policy := networkpolicy.DefineDenyAllNetworkPolicy(name, ns, types, labels)

	return globalhelper.CreateAndWaitUntilNetworkPolicyIsReady(policy, tsparams.WaitingTime)
}

func findListIntersections(listA []string, listB []string) []string {
	var overlap []string

	for _, elementA := range listA {
		for _, elementB := range listB {
			if elementA == elementB {
				overlap = append(overlap, elementA)
			}
		}
	}

	return overlap
}

func getInterfacesList(runningPod corev1.Pod) ([]string, error) {
	links, err := globalhelper.ExecCommand(
		runningPod,
		[]string{"ip", "-j", "link", "show"},
	)
	if err != nil {
		return nil, err
	}

	var interfaceList []string

	var linuxInterfaces []tsparams.IPOutputInterface
	err = json.Unmarshal(links.Bytes(), &linuxInterfaces)

	if err != nil {
		return nil, err
	}

	for _, nodeInterface := range linuxInterfaces {
		if nodeInterface.Master == "" &&
			!strings.Contains(nodeInterface.IfName, "ovn") &&
			!strings.Contains(nodeInterface.IfName, "br") &&
			!strings.Contains(nodeInterface.IfName, "ovs") &&
			!strings.Contains(nodeInterface.IfName, "lo") {
			interfaceList = append(interfaceList, nodeInterface.IfName)
		}
	}

	if len(interfaceList) < 1 {
		return nil, fmt.Errorf("there is no multus interfaces available on node")
	}

	return interfaceList, nil
}

func defineDeploymentBasedOnArgs(
	name string, namespace string, replicaNumber int32, privileged bool, multus []string, label map[string]string) *appsv1.Deployment {
	deploymentStruct := deployment.DefineDeployment(name, namespace,
		globalhelper.Configuration.General.TestImage, tsparams.TestDeploymentLabels)
	deployment.RedefineWithReplicaNumber(deploymentStruct, replicaNumber)

	if privileged {
		deployment.RedefineWithContainersSecurityContextAll(deploymentStruct)
	}

	if label != nil {
		deployment.RedefineWithLabels(deploymentStruct, label)
	}

	if len(multus) > 0 {
		deploymentStruct = deployment.RedefineWithMultus(deploymentStruct, multus)
	}

	return deploymentStruct
}

func defineDeploymentWithContainers(replica int32, containers int,
	name string) (*appsv1.Deployment, error) {
	if containers < 1 {
		return nil, errors.New("invalid containers number")
	}

	deploymentStruct := defineDeploymentBasedOnArgs(name, tsparams.TestNetworkingNameSpace, replica, false, nil, nil)

	globalhelper.AppendContainersToDeployment(deploymentStruct, containers-1, globalhelper.Configuration.General.TestImage)
	deployment.RedefineWithReplicaNumber(deploymentStruct, replica)

	return deploymentStruct, nil
}

func defineDaemonSetBasedOnArgs(nadName string, labels map[string]string) error {
	testDaemonset := daemonset.DefineDaemonSet(tsparams.TestNetworkingNameSpace,
		globalhelper.Configuration.General.TestImage, tsparams.TestDeploymentLabels, "daemonsetnetworkingput")
	daemonset.RedefineWithMultus(testDaemonset, nadName)
	daemonset.RedefineDaemonSetWithNodeSelector(testDaemonset, map[string]string{globalhelper.Configuration.General.CnfNodeLabel: ""})

	if labels != nil {
		daemonset.RedefineDaemonSetWithLabel(testDaemonset, labels)
	}

	return globalhelper.CreateAndWaitUntilDaemonSetIsReady(testDaemonset, tsparams.WaitingTime)
}

func defineAndCreatePrivilegedDaemonset() error {
	daemonSet := daemonset.DefineDaemonSet(tsparams.TestNetworkingNameSpace, globalhelper.Configuration.General.TestImage,
		tsparams.TestDeploymentLabels, "daemonsetnetworkingput")
	daemonset.RedefineDaemonSetWithNodeSelector(daemonSet, map[string]string{globalhelper.Configuration.General.WorkerNodeLabel: ""})
	daemonset.RedefineWithPrivilegeAndHostNetwork(daemonSet)

	err := globalhelper.CreateAndWaitUntilDaemonSetIsReady(daemonSet, tsparams.WaitingTime)
	if err != nil {
		return err
	}

	return nil
}

func execCmdOnPodsListInNamespace(command []string, execOn string) error {
	runningTestPods, err := globalhelper.APIClient.Pods(tsparams.TestNetworkingNameSpace).List(
		context.Background(),
		metav1.ListOptions{})
	if err != nil {
		return err
	}

	var execOcPods *corev1.PodList

	switch execOn {
	case "all":
		execOcPods = runningTestPods

	case "first":
		execOcPods = &corev1.PodList{
			TypeMeta: runningTestPods.TypeMeta,
			ListMeta: runningTestPods.ListMeta,
			Items:    []corev1.Pod{runningTestPods.Items[0]}}
	default:
		return fmt.Errorf("invalid parameter %s", execOn)
	}

	for _, runningPod := range execOcPods.Items {
		_, err := globalhelper.ExecCommand(runningPod, command)
		if err != nil {
			return err
		}
	}

	return nil
}

func createContainerSpecsFromContainerPorts(ports []corev1.ContainerPort) []corev1.Container {
	numContainers := len(ports)
	containerSpecs := []corev1.Container{}

	for index := 0; index < numContainers; index++ {
		containerSpecs = append(containerSpecs,
			corev1.Container{
				Name:    fmt.Sprintf("%s-%d", tsparams.TestDeploymentAName, index),
				Image:   globalhelper.Configuration.General.TestImage,
				Command: []string{"/bin/bash", "-c", "sleep INF"},
				Ports:   []corev1.ContainerPort{ports[index]},
			},
		)
	}

	return containerSpecs
}
