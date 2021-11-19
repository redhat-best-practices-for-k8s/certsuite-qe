package nethelper

import (
	"context"
	"fmt"
	"github.com/test-network-function/cnfcert-tests-verification/tests/networking/netparameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/deployment"
	"time"

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
func ValidateIfReportsAreValid(tcName string) error {
	glog.V(5).Info("Verify test case status in Junit report")
	junitTestReport, err := globalhelper.OpenJunitTestReport()
	if err != nil {
		return err
	}
	if !globalhelper.IsTestCasePassedInJunitReport(junitTestReport, tcName) {
		return fmt.Errorf("test case %s is not in expected passed state in junit report", tcName)
	}
	glog.V(5).Info("Verify test case status in claim report file")
	claimReport, err := globalhelper.OpenClaimReport()
	if err != nil {
		return err
	}
	testPassed, err := globalhelper.IsTestCasePassedInClaimReport(tcName, *claimReport)
	if err != nil {
		return err
	}
	if !testPassed {
		return fmt.Errorf("test case %s is not in expected passed state in claim report", tcName)
	}
	return nil
}

// DefineAndCreateDeploymentOnCluster defines deployment resource and creates it on cluster
func DefineAndCreateDeploymentOnCluster(replicaNumber int32) error {
	deploymentUnderTest := deployment.RedefineWithReplicaNumber(
		deployment.DefineDeployment(
			netparameters.TestNetworkingNameSpace,
			globalhelper.Configuration.General.TestImage,
			netparameters.TestDeploymentLabels),
		replicaNumber)

	err := CreateAndWaitUntilDeploymentIsReady(deploymentUnderTest, netparameters.WaitingTime)
	if err != nil {
		return err
	}
	return nil
}
