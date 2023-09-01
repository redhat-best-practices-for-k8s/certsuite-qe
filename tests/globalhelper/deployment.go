package globalhelper

import (
	"context"
	"fmt"
	"time"

	"github.com/golang/glog"
	appsv1 "k8s.io/api/apps/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	. "github.com/onsi/gomega"
)

// IsDeploymentReady checks if a deployment is ready.
func IsDeploymentReady(operatorNamespace string, deploymentName string) (bool, error) {
	testDeployment, err := GetAPIClient().Deployments(operatorNamespace).Get(
		context.TODO(),
		deploymentName,
		metav1.GetOptions{},
	)
	if err != nil {
		return false, err
	}

	// Ensure the number of ready replicas matches the desired number of replicas.
	if testDeployment.Status.ReadyReplicas == *testDeployment.Spec.Replicas {
		return true, nil
	}

	return false, nil
}

// CreateAndWaitUntilDeploymentIsReady creates deployment and wait until all deployment replicas are up and running.
func CreateAndWaitUntilDeploymentIsReady(deployment *appsv1.Deployment, timeout time.Duration) error {
	runningDeployment, err := GetAPIClient().Deployments(deployment.Namespace).Create(
		context.TODO(),
		deployment,
		metav1.CreateOptions{})

	if k8serrors.IsAlreadyExists(err) {
		glog.V(5).Info(fmt.Sprintf("deployment %s already exists", deployment.Name))
	} else if err != nil {
		return fmt.Errorf("failed to create deployment %q (ns %s): %w", deployment.Name, deployment.Namespace, err)
	}

	Eventually(func() bool {
		status, err := IsDeploymentReady(runningDeployment.Namespace, runningDeployment.Name)
		if err != nil {
			glog.V(5).Info(fmt.Sprintf(
				"deployment %s is not ready, retry in 5 seconds", runningDeployment.Name))

			return false
		}

		return status
	}, timeout, retryInterval*time.Second).Should(Equal(true), "Deployment is not ready ",
		getDeploymentStatus(deployment.Name, deployment.Namespace))

	return nil
}

func getDeploymentStatus(name, namespaces string) string {
	dep, err := GetAPIClient().Deployments(namespaces).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return time.Now().String() + " " + err.Error()
	}

	return time.Now().String() + " " + dep.Status.String()
}
