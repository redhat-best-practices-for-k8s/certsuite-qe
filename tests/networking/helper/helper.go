package helper

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/container"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/nad"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/nodes"

	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/daemonset"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/networkpolicy"

	tsparams "github.com/test-network-function/cnfcert-tests-verification/tests/networking/parameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/deployment"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/service"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/utils/pointer"

	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
	appsv1 "k8s.io/api/apps/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// DefineAndCreateDeploymentOnCluster defines deployment resource and creates it on cluster.
func DefineAndCreateDeploymentOnCluster(replicaNumber int32, namespace string) error {
	deploymentUnderTest := defineDeploymentBasedOnArgs(tsparams.TestDeploymentAName, namespace,
		replicaNumber, false, nil, nil)

	return globalhelper.CreateAndWaitUntilDeploymentIsReady(deploymentUnderTest, tsparams.WaitingTime)
}

// DefineAndCreateDeploymentWithMultusOnCluster defines deployment resource and creates it on cluster.
func DefineAndCreateDeploymentWithMultusOnCluster(name, namespace string, nadNames []string, replicaNumber int32) error {
	deploymentUnderTest := defineDeploymentBasedOnArgs(name, namespace,
		replicaNumber, true, nadNames, nil)

	return globalhelper.CreateAndWaitUntilDeploymentIsReady(deploymentUnderTest, tsparams.WaitingTime)
}

// DefineAndCreateDeploymentWithMultusAndSkipLabelOnCluster defines deployment resource and creates it on cluster.
func DefineAndCreateDeploymentWithMultusAndSkipLabelOnCluster(
	name, namespace string, nadNames []string, replicaNumber int32) error {
	deploymentUnderTest := defineDeploymentBasedOnArgs(
		name, namespace,
		replicaNumber,
		false,
		nadNames,
		tsparams.NetworkingTestMultusSkipLabel)

	return globalhelper.CreateAndWaitUntilDeploymentIsReady(deploymentUnderTest, tsparams.WaitingTime)
}

// DefineAndCreatePrivilegedDeploymentOnCluster defines deployment resource and creates it on cluster.
func DefineAndCreatePrivilegedDeploymentOnCluster(replicaNumber int32, namespace string) error {
	deploymentUnderTest := defineDeploymentBasedOnArgs(tsparams.TestDeploymentAName, namespace,
		replicaNumber, true,
		nil, nil)

	return globalhelper.CreateAndWaitUntilDeploymentIsReady(deploymentUnderTest, tsparams.WaitingTime)
}

// DefineAndCreateDeploymentWithSkippedLabelOnCluster defines deployment resource and creates it on cluster.
func DefineAndCreateDeploymentWithSkippedLabelOnCluster(replicaNumber int32, namespace string) error {
	deploymentUnderTest := defineDeploymentBasedOnArgs(tsparams.TestDeploymentAName, namespace,
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

func DefineAndCreateDeployment(deploymentName, namespace string, replicaNumber int32) error {
	deploymentUnderTest := defineDeploymentBasedOnArgs(deploymentName, namespace,
		replicaNumber, false, nil, nil)

	return globalhelper.CreateAndWaitUntilDeploymentIsReady(deploymentUnderTest, tsparams.WaitingTime)
}

func DefineAndCreateDaemonsetWithMultusOnCluster(nadName, namespace, daemonsetName string) error {
	return defineDaemonSetBasedOnArgs(nadName, namespace, daemonsetName, nil)
}

func DefineAndCreateDaemonsetWithMultusAndSkipLabelOnCluster(nadName, namespace, daemonsetName string) error {
	return defineDaemonSetBasedOnArgs(nadName, namespace, daemonsetName, tsparams.NetworkingTestMultusSkipLabel)
}

// DefineAndCreateDeploymentOnCluster defines deployment resource and creates it on cluster.
func DefineAndCreateDeploymentWithContainerPorts(replicaNumber int32, ports []corev1.ContainerPort, namespace string) error {
	deploymentUnderTest, err := DefineDeploymentWithContainers(replicaNumber, len(ports), tsparams.TestDeploymentAName, namespace)
	if err != nil {
		return err
	}

	portSpecs := createContainerSpecsFromContainerPorts(ports)

	deployment.RedefineWithContainerSpecs(deploymentUnderTest, portSpecs)

	return globalhelper.CreateAndWaitUntilDeploymentIsReady(deploymentUnderTest, tsparams.WaitingTime)
}

// ExecCmdOnOnePodInNamespace runs command on the first available pod in namespace.
func ExecCmdOnOnePodInNamespace(command []string, namespace string) error {
	return execCmdOnPodsListInNamespace(command, "first", namespace)
}

func ExecCmdOnAllPodInNamespace(command []string, namespace string) error {
	return execCmdOnPodsListInNamespace(command, "all", namespace)
}

func RedefineServiceToHeadless(service *corev1.Service) {
	service.Spec.ClusterIP = corev1.ClusterIPNone
}

// DefineAndCreateServiceOnCluster defines service resource and creates it on cluster.
func DefineAndCreateServiceOnCluster(name, namespace string, port int32, targetPort int32, withNodePort, headless bool,
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
		var err error

		testService, err = service.RedefineWithNodePort(testService)
		if err != nil {
			return err
		}
	}

	if headless {
		RedefineServiceToHeadless(testService)
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

func DefineAndCreateNadOnCluster(name, namespace string, network string) error {
	nadOneInterface := nad.DefineNad(name, namespace)

	if network != "" {
		nadOneInterface = nad.RedefineNadWithWhereaboutsIpam(nadOneInterface, network)
	}

	err := globalhelper.GetAPIClient().Create(context.TODO(), nadOneInterface)

	if k8serrors.IsAlreadyExists(err) {
		return nil
	} else if err != nil {
		return fmt.Errorf("failed to create nad on cluster: %w", err)
	}

	return nil
}

func GetClusterMultusInterfaces(namespace string) ([]string, error) {
	err := defineAndCreatePrivilegedDaemonset(namespace)
	if err != nil {
		return nil, err
	}

	podsList, err := globalhelper.GetListOfPodsInNamespace(namespace)
	if err != nil {
		return nil, err
	}

	var nodesInterfacesList [][]string

	for _, runningPod := range podsList.Items {
		isMasterNode, err := nodes.IsNodeMaster(runningPod.Spec.NodeName, globalhelper.GetAPIClient())
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
		globalhelper.GetConfiguration().General.TestImage, tsparams.TestDeploymentLabels)
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

func DefineDeploymentWithContainers(replica int32, containers int,
	name, namespace string) (*appsv1.Deployment, error) {
	if containers < 1 {
		return nil, errors.New("invalid containers number")
	}

	deploymentStruct := defineDeploymentBasedOnArgs(name, namespace, replica, false, nil, nil)

	globalhelper.AppendContainersToDeployment(deploymentStruct, containers-1, globalhelper.GetConfiguration().General.TestImage)
	deployment.RedefineWithReplicaNumber(deploymentStruct, replica)

	return deploymentStruct, nil
}

func defineDaemonSetBasedOnArgs(nadName, namespace, daemonsetName string, labels map[string]string) error {
	testDaemonset := daemonset.DefineDaemonSet(namespace,
		globalhelper.GetConfiguration().General.TestImage, tsparams.TestDeploymentLabels, daemonsetName)
	daemonset.RedefineWithMultus(testDaemonset, nadName)
	//nolint:lll
	daemonset.RedefineDaemonSetWithNodeSelector(testDaemonset, map[string]string{globalhelper.GetConfiguration().General.CnfNodeLabel: ""})

	if labels != nil {
		daemonset.RedefineDaemonSetWithLabel(testDaemonset, labels)
	}

	return globalhelper.CreateAndWaitUntilDaemonSetIsReady(testDaemonset, tsparams.WaitingTime)
}

func defineAndCreatePrivilegedDaemonset(namespace string) error {
	daemonSet := daemonset.DefineDaemonSet(namespace, globalhelper.GetConfiguration().General.TestImage,
		tsparams.TestDeploymentLabels, "daemonsetnetworkingput")
	daemonset.RedefineDaemonSetWithNodeSelector(daemonSet, map[string]string{globalhelper.GetConfiguration().General.WorkerNodeLabel: ""})
	daemonset.RedefineWithPrivilegeAndHostNetwork(daemonSet)

	err := globalhelper.CreateAndWaitUntilDaemonSetIsReady(daemonSet, tsparams.WaitingTime)
	if err != nil {
		return err
	}

	return nil
}

func execCmdOnPodsListInNamespace(command []string, execOn, namespace string) error {
	runningTestPods, err := globalhelper.GetAPIClient().Pods(namespace).List(
		context.TODO(),
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
				Image:   globalhelper.GetConfiguration().General.TestImage,
				Command: []string{"/bin/bash", "-c", "sleep INF"},
				Ports:   []corev1.ContainerPort{ports[index]},
			},
		)
	}

	return containerSpecs
}

func DefineDeploymentWithContainerPorts(name, namespace string,
	replicaNumber int32, ports []corev1.ContainerPort) (*appsv1.Deployment, error) {
	if len(ports) < 1 {
		return nil, errors.New("invalid number of containers")
	}

	deploymentStruct := deployment.DefineDeployment(name, namespace,
		globalhelper.GetConfiguration().General.TestImage, tsparams.TestDeploymentLabels)

	globalhelper.AppendContainersToDeployment(deploymentStruct, len(ports)-1, globalhelper.GetConfiguration().General.TestImage)
	deployment.RedefineWithReplicaNumber(deploymentStruct, replicaNumber)

	portSpecs := container.CreateContainerSpecsFromContainerPorts(ports, globalhelper.GetConfiguration().General.TestImage, "test")

	deployment.RedefineWithContainerSpecs(deploymentStruct, portSpecs)

	return deploymentStruct, nil
}

func DefineDpdkPod(podName, namespace string) *corev1.Pod {
	cpuLimit := "1"
	memoryLimit := "512Mi"
	containerCommand := []string{"/bin/sh", "-c", "touch /tmp/healthy", "sleep infinity"}

	containerResource := corev1.ResourceRequirements{
		Limits: corev1.ResourceList{
			corev1.ResourceCPU:    resource.MustParse(cpuLimit),
			corev1.ResourceMemory: resource.MustParse(memoryLimit),
		},
		Requests: corev1.ResourceList{
			corev1.ResourceCPU:    resource.MustParse(cpuLimit),
			corev1.ResourceMemory: resource.MustParse(memoryLimit),
		},
	}

	containerSecurityContext := &corev1.SecurityContext{
		Privileged: pointer.Bool(true),
		RunAsUser:  pointer.Int64(0),
		Capabilities: &corev1.Capabilities{
			Add: []corev1.Capability{"IPC_LOCK", "SYS_RESOURCE", "NET_RAW"}},
	}

	annotations := make(map[string]string)
	annotations["k8s.v1.cni.cncf.io/networks"] = "sriovnet1"

	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:        podName,
			Namespace:   namespace,
			Labels:      tsparams.TestDeploymentLabels,
			Annotations: annotations,
		},
		Spec: corev1.PodSpec{
			TerminationGracePeriodSeconds: pointer.Int64(0),
			Containers: []corev1.Container{
				{
					Name:            "app-container",
					Image:           "registry.redhat.io/openshift4/dpdk-base-rhel8:v4.9",
					Command:         containerCommand,
					Resources:       containerResource,
					SecurityContext: containerSecurityContext,
					VolumeMounts: []corev1.VolumeMount{
						{
							Name:      "hugepage",
							MountPath: "/dev/hugepages",
						},
					}},
			},
			Volumes: []corev1.Volume{
				{
					Name: "hugepage",
					VolumeSource: corev1.VolumeSource{
						EmptyDir: &corev1.EmptyDirVolumeSource{
							Medium: corev1.StorageMediumHugePages,
						},
					},
				},
			},
		},
	}
}
