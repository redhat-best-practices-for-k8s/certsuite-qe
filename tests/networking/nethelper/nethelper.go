package nethelper

import (
	"context"
	"fmt"
	"time"

	"github.com/test-network-function/cnfcert-tests-verification/tests/globalparameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/networking/netparameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/deployment"
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

func isDeploymentReady(operatorNamespace string, deploymentName string) (bool, error) {
	testDeployment, err := globalhelper.APIClient.Deployments(operatorNamespace).Get(
		context.Background(),
		deploymentName,
		metav1.GetOptions{},
	)
	if err != nil {
		return false, err
	}

	if testDeployment.Status.ReadyReplicas > 0 {
		if testDeployment.Status.Replicas == testDeployment.Status.ReadyReplicas {
			return true, nil
		}
	}

	return false, nil
}

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

func defineDeploymentBasedOnArgs(replicaNumber int32, privileged bool, label map[string]string) *v1.Deployment {
	deploymentStruct := deployment.RedefineWithReplicaNumber(
		deployment.DefineDeployment(
			"networkingput",
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

	return deploymentStruct
}

// CreateAndWaitUntilDeploymentIsReady creates deployment and wait until all deployment replicas are up and running.
func CreateAndWaitUntilDeploymentIsReady(deployment *v1.Deployment, timeout time.Duration) error {
	runningDeployment, err := globalhelper.APIClient.Deployments(deployment.Namespace).Create(
		context.Background(),
		deployment,
		metav1.CreateOptions{})
	if err != nil {
		return err
	}

	Eventually(func() bool {
		status, err := isDeploymentReady(runningDeployment.Namespace, runningDeployment.Name)
		if err != nil {
			glog.V(5).Info(fmt.Sprintf(
				"deployment %s is not ready, retry in 5 seconds", runningDeployment.Name))

			return false
		}

		return status
	}, timeout, 5*time.Second).Should(Equal(true), "Deployment is not ready")

	return nil
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
			glog.V(5).Info(fmt.Sprintf(
				"daemonset %s is not ready, retry in 5 seconds", runningDaemonSet.Name))

			return false
		}

		return status
	}, timeout, 5*time.Second).Should(Equal(true), "DaemonSet is not ready")

	return nil
}

// ValidateIfReportsAreValid test if report is valid for given test case.
func ValidateIfReportsAreValid(tcName string, tcExpectedStatus string) error {
	glog.V(5).Info("Verify test case status in Junit report")

	junitTestReport, err := globalhelper.OpenJunitTestReport()

	if err != nil {
		return err
	}

	claimReport, err := globalhelper.OpenClaimReport()

	if err != nil {
		return err
	}

	err = globalhelper.IsExpectedStatusParamValid(tcExpectedStatus)

	if err != nil {
		return err
	}

	isTestCaseInValidStatusInJunitReport := globalhelper.IsTestCasePassedInJunitReport
	isTestCaseInValidStatusInClaimReport := globalhelper.IsTestCasePassedInClaimReport

	if tcExpectedStatus == globalparameters.TestCaseFailed {
		isTestCaseInValidStatusInJunitReport = globalhelper.IsTestCaseFailedInJunitReport
		isTestCaseInValidStatusInClaimReport = globalhelper.IsTestCaseFailedInClaimReport
	}

	if tcExpectedStatus == globalparameters.TestCaseSkipped {
		isTestCaseInValidStatusInJunitReport = globalhelper.IsTestCaseSkippedInJunitReport
		isTestCaseInValidStatusInClaimReport = globalhelper.IsTestCaseSkippedInClaimReport
	}

	if !isTestCaseInValidStatusInJunitReport(junitTestReport, tcName) {
		return fmt.Errorf("test case %s is not in expected %s state in junit report", tcName, tcExpectedStatus)
	}

	glog.V(5).Info("Verify test case status in claim report file")

	testPassed, err := isTestCaseInValidStatusInClaimReport(tcName, *claimReport)

	if err != nil {
		return err
	}

	if !testPassed {
		return fmt.Errorf("test case %s is not in expected %s state in claim report", tcName, tcExpectedStatus)
	}

	return nil
}

// DefineAndCreateDeploymentOnCluster defines deployment resource and creates it on cluster.
func DefineAndCreateDeploymentOnCluster(replicaNumber int32) error {
	deploymentUnderTest := defineDeploymentBasedOnArgs(replicaNumber, false, nil)

	return CreateAndWaitUntilDeploymentIsReady(deploymentUnderTest, netparameters.WaitingTime)
}

// DefineAndCreatePrivilegedDeploymentOnCluster defines deployment resource and creates it on cluster.
func DefineAndCreatePrivilegedDeploymentOnCluster(replicaNumber int32) error {
	deploymentUnderTest := defineDeploymentBasedOnArgs(replicaNumber, true, nil)

	return CreateAndWaitUntilDeploymentIsReady(deploymentUnderTest, netparameters.WaitingTime)
}

// DefineAndCreateDeploymentWithSkippedLabelOnCluster defines deployment resource and creates it on cluster.
func DefineAndCreateDeploymentWithSkippedLabelOnCluster(replicaNumber int32) error {
	deploymentUnderTest := defineDeploymentBasedOnArgs(
		replicaNumber,
		true,
		netparameters.NetworkingTestSkipLabel)
	err := CreateAndWaitUntilDeploymentIsReady(deploymentUnderTest, netparameters.WaitingTime)

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
