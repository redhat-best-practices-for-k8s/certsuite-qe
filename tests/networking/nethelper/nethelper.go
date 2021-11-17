package nethelper

import (
	"context"
	"fmt"
	"time"

	"github.com/golang/glog"

	. "github.com/onsi/gomega"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func isDeploymentReady(operatorNamespace string, deploymentName string) (bool, error) {
	deployment, err := globalhelper.ApiClient.Deployments(operatorNamespace).Get(context.Background(), deploymentName, metav1.GetOptions{})
	if err != nil {
		return false, err
	}
	if deployment.Status.ReadyReplicas > 0 {
		if deployment.Status.Replicas == deployment.Status.ReadyReplicas {
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
