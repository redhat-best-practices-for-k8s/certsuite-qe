package globalhelper

import (
	"context"
	"fmt"
	"time"

	"github.com/golang/glog"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/client"
	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	. "github.com/onsi/gomega"
)

// IsDeploymentReady checks if a deployment is ready.
func IsDeploymentReady(operatorNamespace string, deploymentName string) (bool, error) {
	testDeployment, err := APIClient.Deployments(operatorNamespace).Get(
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

// IsDeploymentInstalled checks if deployment is installed.
func IsDeploymentInstalled(
	cs *client.ClientSet, operatorNamespace string, operatorDeploymentName string) (bool, error) {
	_, err := APIClient.Deployments(operatorNamespace).Get(context.Background(), operatorDeploymentName, metav1.GetOptions{})
	if err != nil {
		return false, err
	}

	return true, nil
}

// CreateAndWaitUntilDeploymentIsReady creates deployment and wait until all deployment replicas are up and running.
func CreateAndWaitUntilDeploymentIsReady(deployment *v1.Deployment, timeout time.Duration) error {
	runningDeployment, err := APIClient.Deployments(deployment.Namespace).Create(
		context.Background(),
		deployment,
		metav1.CreateOptions{})
	if err != nil {
		return err
	}

	Eventually(func() bool {
		status, err := IsDeploymentReady(runningDeployment.Namespace, runningDeployment.Name)
		if err != nil {
			glog.V(5).Info(fmt.Sprintf(
				"deployment %s is not ready, retry in 5 seconds", runningDeployment.Name))

			return false
		}

		return status
	}, timeout, 5*time.Second).Should(Equal(true), "Deployment is not ready")

	return nil
}
