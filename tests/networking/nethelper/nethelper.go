package nethelper

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/nad"

	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/daemonset"

	"github.com/test-network-function/cnfcert-tests-verification/tests/networking/netparameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/deployment"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/namespaces"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/rbac"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/service"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"

	"github.com/golang/glog"

	. "github.com/onsi/gomega"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func isDaemonSetReady(operatorNamespace string, daemonSetName string) (bool, error) {
	daemonSet, err := globalhelper.APIClient.DaemonSets(operatorNamespace).Get(
		context.Background(),
		daemonSetName,
		metav1.GetOptions{},
	)
	if err != nil {
		return false, err
	}

	if daemonSet.Status.NumberReady > 0 {
		if daemonSet.Status.NumberUnavailable == 0 {
			return true, nil
		}
	}

	return false, nil
}

func defineDeploymentBasedOnArgs(
	name string, replicaNumber int32, privileged bool, multus string, label map[string]string) *v1.Deployment {
	deploymentStruct := deployment.RedefineWithReplicaNumber(
		deployment.DefineDeployment(
			name,
			netparameters.TestNetworkingNameSpace,
			globalhelper.Configuration.General.TestImage,
			netparameters.TestDeploymentLabels),
		replicaNumber)
	if privileged {
		deploymentStruct = deployment.RedefineWithContainersSecurityContextAll(deploymentStruct)
	}

	if label != nil {
		deploymentStruct = deployment.RedefineWithLabels(deploymentStruct, label)
	}

	if len(multus) > 1 {
		deploymentStruct = deployment.RedefineWithMultus(deploymentStruct, multus)
	}

	return deploymentStruct
}

// CreateAndWaitUntilDaemonSetIsReady creates daemonSet and wait until all deployment replicas are up and running.
func CreateAndWaitUntilDaemonSetIsReady(daemonSet *v1.DaemonSet, timeout time.Duration) error {
	runningDaemonSet, err := globalhelper.APIClient.DaemonSets(daemonSet.Namespace).Create(
		context.Background(),
		daemonSet,
		metav1.CreateOptions{})
	if err != nil {
		return err
	}

	Eventually(func() bool {
		status, err := isDaemonSetReady(runningDaemonSet.Namespace, runningDaemonSet.Name)
		if err != nil {
			glog.Fatal(fmt.Sprintf(
				"daemonset %s is not ready, retry in 5 seconds", runningDaemonSet.Name))

			return false
		}

		return status
	}, timeout, 5*time.Second).Should(Equal(true), "DaemonSet is not ready")

	return nil
}

// DefineAndCreateDeploymentOnCluster defines deployment resource and creates it on cluster.
func DefineAndCreateDeploymentOnCluster(replicaNumber int32) error {
	deploymentUnderTest := defineDeploymentBasedOnArgs(netparameters.TestDeploymentAName, replicaNumber, false, "", nil)

	return globalhelper.CreateAndWaitUntilDeploymentIsReady(deploymentUnderTest, netparameters.WaitingTime)
}

// DefineAndCreateDeploymentWithMultusOnCluster defines deployment resource and creates it on cluster.
func DefineAndCreateDeploymentWithMultusOnCluster(name string, nadName string, replicaNumber int32) error {
	deploymentUnderTest := defineDeploymentBasedOnArgs(name, replicaNumber, false, nadName, nil)

	return globalhelper.CreateAndWaitUntilDeploymentIsReady(deploymentUnderTest, netparameters.WaitingTime)
}

// DefineAndCreatePrivilegedDeploymentOnCluster defines deployment resource and creates it on cluster.
func DefineAndCreatePrivilegedDeploymentOnCluster(replicaNumber int32) error {
	deploymentUnderTest := defineDeploymentBasedOnArgs(
		netparameters.TestDeploymentAName,
		replicaNumber,
		true,
		"",
		nil)

	return globalhelper.CreateAndWaitUntilDeploymentIsReady(deploymentUnderTest, netparameters.WaitingTime)
}

// DefineAndCreateDeploymentWithSkippedLabelOnCluster defines deployment resource and creates it on cluster.
func DefineAndCreateDeploymentWithSkippedLabelOnCluster(replicaNumber int32) error {
	deploymentUnderTest := defineDeploymentBasedOnArgs(
		netparameters.TestDeploymentAName,
		replicaNumber,
		true,
		"",
		netparameters.NetworkingTestSkipLabel)
	err := globalhelper.CreateAndWaitUntilDeploymentIsReady(deploymentUnderTest, netparameters.WaitingTime)

	if err != nil {
		return err
	}

	return nil
}

func DefineAndCreateDeamonsetWithMultusOnCluster(nadName string) error {
	testDaemonset := daemonset.RedefineDaemonSetWithNodeSelector(daemonset.RedefineWithMultus(
		daemonset.DefineDaemonSet(
			netparameters.TestNetworkingNameSpace,
			globalhelper.Configuration.General.TestImage,
			netparameters.TestDeploymentLabels),
		nadName,
	), map[string]string{globalhelper.Configuration.General.CnfNodeLabel: ""})

	err := CreateAndWaitUntilDaemonSetIsReady(testDaemonset, netparameters.WaitingTime)
	if err != nil {
		return err
	}

	return nil
}

func defineAndCreatePrivilegedDaemonset() error {
	daemonSet := daemonset.RedefineWithPrivilegeAndHostNetwork(daemonset.RedefineDaemonSetWithNodeSelector(
		daemonset.DefineDaemonSet(
			netparameters.TestNetworkingNameSpace,
			globalhelper.Configuration.General.TestImage,
			netparameters.TestDeploymentLabels,
		), map[string]string{globalhelper.Configuration.General.WorkerNodeLabel: ""}))
	err := CreateAndWaitUntilDaemonSetIsReady(daemonSet, netparameters.WaitingTime)

	if err != nil {
		return err
	}

	return nil
}

// AllowAuthenticatedUsersRunPrivilegedContainers adds all authenticated users to privileged group.
func AllowAuthenticatedUsersRunPrivilegedContainers() error {
	_, err := globalhelper.APIClient.ClusterRoleBindings().Get(
		context.Background(),
		"system:openshift:scc:privileged",
		metav1.GetOptions{},
	)
	if k8serrors.IsNotFound(err) {
		glog.V(5).Info("RBAC policy is not found")

		roleBind := rbac.DefineClusterRoleBinding(
			*rbac.DefineRbacAuthorizationClusterRoleRef("system:openshift:scc:privileged"),
			*rbac.DefineRbacAuthorizationClusterGroupSubjects([]string{"system:authenticated"}),
		)
		_, err = globalhelper.APIClient.ClusterRoleBindings().Create(
			context.Background(),
			roleBind,
			metav1.CreateOptions{},
		)

		if err != nil {
			return err
		}

		glog.V(5).Info("RBAC policy created")

		return nil
	} else if err == nil {
		glog.V(5).Info("RBAC policy detected")
	}

	glog.V(5).Info("error to query RBAC policy")

	return err
}

func execCmdOnPodsListInNamespace(command []string, execOn string) error {
	runningTestPods, err := globalhelper.APIClient.Pods(netparameters.TestNetworkingNameSpace).List(
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
		netparameters.TestNetworkingNameSpace,
		port,
		targetPort,
		corev1.ProtocolTCP,
		netparameters.TestDeploymentLabels)

	if withNodePort {
		var err error

		testService, err = service.RedefineWithNodePort(testService)
		if err != nil {
			return err
		}
	}

	_, err := globalhelper.APIClient.Services(netparameters.TestNetworkingNameSpace).Create(
		context.Background(),
		testService, metav1.CreateOptions{})

	return err
}

func DefineAndCreateNadOnCluster(name string, intName string, network string) error {
	nadOneInterface := nad.DefineNad(name, netparameters.TestNetworkingNameSpace, intName)

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

	podsList, err := globalhelper.GetListOfPodsInNamespace(netparameters.TestNetworkingNameSpace)
	if err != nil {
		return nil, err
	}

	var nodesInterfacesList [][]string

	for _, runningPod := range podsList.Items {
		nodeIterfaces, err := getInterfacesList(runningPod)
		if err != nil {
			return nil, err
		}
		nodesInterfacesList = append(nodesInterfacesList, nodeIterfaces)
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

	var linuxInterfaces []netparameters.IPOutputInterface
	err = json.Unmarshal(links.Bytes(), &linuxInterfaces)

	if err != nil {
		return nil, err
	}

	for _, nodeInterface := range linuxInterfaces {
		if nodeInterface.Master == "" && nodeInterface.Operstate == "UP" {
			interfaceList = append(interfaceList, nodeInterface.IfName)
		}
	}

	if len(interfaceList) < 1 {
		return nil, fmt.Errorf("there is no multus interfaces available on node")
	}

	return interfaceList, nil
}

func PutDownInterfaceOnNode(name string, state string) error {
	err := defineAndCreatePrivilegedDaemonset()

	if err != nil {
		return err
	}

	podsList, err := globalhelper.GetListOfPodsInNamespace(netparameters.TestNetworkingNameSpace)
	if err != nil {
		return err
	}

	command := []string{"ip", "link", "set", name}

	if state == "up" || state == "down" {
		command = append(command, state)
	} else {
		return fmt.Errorf("invalid argument %s", state)
	}

	output, err := globalhelper.ExecCommand(podsList.Items[0], command)
	if err != nil {
		glog.Fatal(output.String())

		return err
	}

	err = namespaces.Clean(netparameters.TestNetworkingNameSpace, globalhelper.APIClient)

	if err != nil {
		return err
	}

	return nil
}

func PutAllInterfaceUPonAllNodes(name string) error {
	err := namespaces.Clean(netparameters.TestNetworkingNameSpace, globalhelper.APIClient)
	if err != nil {
		return err
	}

	err = defineAndCreatePrivilegedDaemonset()

	if err != nil {
		return err
	}

	podsList, err := globalhelper.GetListOfPodsInNamespace(netparameters.TestNetworkingNameSpace)
	if err != nil {
		return err
	}

	for _, runningPod := range podsList.Items {
		output, err := globalhelper.ExecCommand(runningPod, []string{"ip", "link", "set", name, "up"})

		if err != nil {
			glog.Fatal(output.String())

			return err
		}
	}

	err = namespaces.Clean(netparameters.TestNetworkingNameSpace, globalhelper.APIClient)
	if err != nil {
		return err
	}

	return nil
}
