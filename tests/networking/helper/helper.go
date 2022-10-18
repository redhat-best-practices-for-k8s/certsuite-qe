package helper

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/nad"

	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/daemonset"

	tsparams "github.com/test-network-function/cnfcert-tests-verification/tests/networking/parameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/deployment"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/service"
	corev1 "k8s.io/api/core/v1"

	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// DefineAndCreateDeploymentOnCluster defines deployment resource and creates it on cluster.
func DefineAndCreateDeploymentOnCluster(replicaNumber int32) error {
	deploymentUnderTest := defineDeploymentBasedOnArgs(tsparams.TestDeploymentAName, replicaNumber, false, nil, nil)

	return globalhelper.CreateAndWaitUntilDeploymentIsReady(deploymentUnderTest, tsparams.WaitingTime)
}

// DefineAndCreateDeploymentWithMultusOnCluster defines deployment resource and creates it on cluster.
func DefineAndCreateDeploymentWithMultusOnCluster(name string, nadNames []string, replicaNumber int32) error {
	deploymentUnderTest := defineDeploymentBasedOnArgs(name, replicaNumber, true, nadNames, nil)

	return globalhelper.CreateAndWaitUntilDeploymentIsReady(deploymentUnderTest, tsparams.WaitingTime)
}

// DefineAndCreateDeploymentWithMultusAndSkipLabelOnCluster defines deployment resource and creates it on cluster.
func DefineAndCreateDeploymentWithMultusAndSkipLabelOnCluster(
	name string, nadNames []string, replicaNumber int32) error {
	deploymentUnderTest := defineDeploymentBasedOnArgs(
		name, replicaNumber,
		false,
		nadNames,
		tsparams.NetworkingTestMultusSkipLabel)

	return globalhelper.CreateAndWaitUntilDeploymentIsReady(deploymentUnderTest, tsparams.WaitingTime)
}

// DefineAndCreatePrivilegedDeploymentOnCluster defines deployment resource and creates it on cluster.
func DefineAndCreatePrivilegedDeploymentOnCluster(replicaNumber int32) error {
	deploymentUnderTest := defineDeploymentBasedOnArgs(tsparams.TestDeploymentAName, replicaNumber, true,
		nil, nil)

	return globalhelper.CreateAndWaitUntilDeploymentIsReady(deploymentUnderTest, tsparams.WaitingTime)
}

// DefineAndCreateDeploymentWithSkippedLabelOnCluster defines deployment resource and creates it on cluster.
func DefineAndCreateDeploymentWithSkippedLabelOnCluster(replicaNumber int32) error {
	deploymentUnderTest := defineDeploymentBasedOnArgs(tsparams.TestDeploymentAName, replicaNumber,
		true, nil, tsparams.NetworkingTestSkipLabel)

	err := globalhelper.CreateAndWaitUntilDeploymentIsReady(deploymentUnderTest, tsparams.WaitingTime)
	if err != nil {
		return fmt.Errorf("failed to create deployment: %w", err)
	}

	return nil
}

func DefineAndCreateDeamonsetWithMultusOnCluster(nadName string) error {
	return defineDaemonSetBasedOnArgs(nadName, nil)
}

func DefineAndCreateDeamonsetWithMultusAndSkipLabelOnCluster(nadName string) error {
	return defineDaemonSetBasedOnArgs(nadName, tsparams.NetworkingTestMultusSkipLabel)
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

func DefineAndCreateNadOnCluster(name string, intName string, network string) error {
	nadOneInterface := nad.DefineNad(name, tsparams.TestNetworkingNameSpace, intName)

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
		nodeInterfaces, err := getInterfacesList(runningPod)
		if err != nil {
			return nil, err
		}
		nodesInterfacesList = append(nodesInterfacesList, nodeInterfaces)
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
	name string, replicaNumber int32, privileged bool, multus []string, label map[string]string) *v1.Deployment {
	deploymentStruct := deployment.DefineDeployment(name, tsparams.TestNetworkingNameSpace,
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
