package nethelper

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/test-network-function/cnfcert-tests-verification/tests/globalparameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/networking/netparameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/deployment"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/rbac"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"

	"github.com/golang/glog"

	. "github.com/onsi/gomega"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func isDeploymentReady(operatorNamespace string, deploymentName string) (bool, error) {
	testDeployment, err := globalhelper.ApiClient.Deployments(operatorNamespace).Get(context.Background(), deploymentName, metav1.GetOptions{})
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
	daemonSet, err := globalhelper.ApiClient.DaemonSets(operatorNamespace).Get(context.Background(), daemonSetName, metav1.GetOptions{})
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

func defineDeploymentBasedOnArgs(replicaNumber int32, privileged bool) *v1.Deployment {
	deploymentStruct := deployment.RedefineWithReplicaNumber(
		deployment.DefineDeployment(
			netparameters.TestNetworkingNameSpace,
			globalhelper.Configuration.General.TestImage,
			netparameters.TestDeploymentLabels),
		replicaNumber)
	if privileged {
		deploymentStruct = deployment.RedefineWithContainersSecurityContextAll(deploymentStruct)
	}
	return deploymentStruct
}

// CreateAndWaitUntilDeploymentIsReady creates deployment and wait until all deployment replicas are up and running
func CreateAndWaitUntilDeploymentIsReady(deployment *v1.Deployment, timeout time.Duration) error {
	runningDeployment, err := globalhelper.ApiClient.Deployments(deployment.Namespace).Create(
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

// CreateAndWaitUntilDaemonSetIsReady creates daemonSet and wait until all deployment replicas are up and running
func CreateAndWaitUntilDaemonSetIsReady(daemonSet *v1.DaemonSet, timeout time.Duration) error {
	runningDaemonSet, err := globalhelper.ApiClient.DaemonSets(daemonSet.Namespace).Create(
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

// ValidateIfReportsAreValid test if report is valid for given test case
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

// DefineAndCreateDeploymentOnCluster defines deployment resource and creates it on cluster
func DefineAndCreateDeploymentOnCluster(replicaNumber int32) error {
	deploymentUnderTest := defineDeploymentBasedOnArgs(replicaNumber, false)
	err := CreateAndWaitUntilDeploymentIsReady(deploymentUnderTest, netparameters.WaitingTime)
	if err != nil {
		return err
	}
	return nil
}

// DefineAndCreatePrivilegedDeploymentOnCluster defines deployment resource and creates it on cluster
func DefineAndCreatePrivilegedDeploymentOnCluster(replicaNumber int32) error {
	deploymentUnderTest := defineDeploymentBasedOnArgs(replicaNumber, true)
	err := CreateAndWaitUntilDeploymentIsReady(deploymentUnderTest, netparameters.WaitingTime)
	if err != nil {
		return err
	}
	return nil
}

// AllowAuthenticatedUsersRunPrivilegedContainers adds all authenticated users to privileged group
func AllowAuthenticatedUsersRunPrivilegedContainers() error {
	_, err := globalhelper.ApiClient.ClusterRoleBindings().Get(
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
		_, err = globalhelper.ApiClient.ClusterRoleBindings().Create(
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

// GetPartnerPodDefinition returns partner's pod struct
func GetPartnerPodDefinition() (*corev1.Pod, error) {
	podsList, err := globalhelper.GetListOfPodsInNamespace(netparameters.DefaultPartnerPodNamespace)
	if err != nil {
		return nil, err
	}
	var partnerPodIP *corev1.Pod
	for _, runningPod := range podsList.Items {
		if strings.Contains(runningPod.Name, netparameters.DefaultPartnerPodPrefixName) {
			partnerPodIP = &runningPod
		}
	}
	if partnerPodIP == nil {
		return nil, fmt.Errorf("can not detect partner pods in %s namespace",
			netparameters.DefaultPartnerPodNamespace)
	}
	return partnerPodIP, nil
}

// ExecCmdCommandOnOnePodInNamespace runs command on the first available pod in namespace
func ExecCmdCommandOnOnePodInNamespace(command []string) error {
	runningTestPods, err := globalhelper.ApiClient.Pods(netparameters.TestNetworkingNameSpace).List(
		context.Background(),
		metav1.ListOptions{})
	if err != nil {
		return err
	}
	if len(runningTestPods.Items) < 1 {
		return fmt.Errorf("there is no running pods under %s namespace", netparameters.TestNetworkingNameSpace)
	}
	_, err = globalhelper.ExecCommand(runningTestPods.Items[1], command)
	return err
}
