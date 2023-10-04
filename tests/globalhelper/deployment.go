package globalhelper

import (
	"context"
	"fmt"
	"time"

	"github.com/golang/glog"
	appsv1 "k8s.io/api/apps/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	typedappsv1 "k8s.io/client-go/kubernetes/typed/apps/v1"

	. "github.com/onsi/gomega"
)

// IsDeploymentReady checks if a deployment is ready.
func IsDeploymentReady(client typedappsv1.AppsV1Interface, namespace, deploymentName string) (bool, error) {
	testDeployment, err := client.Deployments(namespace).Get(
		context.TODO(),
		deploymentName,
		metav1.GetOptions{},
	)
	if err != nil {
		return false, err
	}

	// Ensure the number of ready replicas matches the desired number of replicas.
	if testDeployment.Status.AvailableReplicas == *testDeployment.Spec.Replicas {
		return true, nil
	}

	return false, nil
}

// CreateAndWaitUntilDeploymentIsReady creates deployment and wait until all deployment replicas are up and running.
func CreateAndWaitUntilDeploymentIsReady(deployment *appsv1.Deployment,
	timeout time.Duration) error {
	return createAndWaitUntilDeploymentIsReady(GetAPIClient().K8sClient.AppsV1(), deployment, timeout)
}

// createAndWaitUntilDeploymentIsReady creates deployment and wait until all deployment replicas are up and running.
func createAndWaitUntilDeploymentIsReady(client typedappsv1.AppsV1Interface, deployment *appsv1.Deployment,
	timeout time.Duration) error {
	runningDeployment, err := client.Deployments(deployment.Namespace).Create(
		context.TODO(),
		deployment,
		metav1.CreateOptions{})

	if k8serrors.IsAlreadyExists(err) {
		glog.V(5).Info(fmt.Sprintf("deployment %s already exists", deployment.Name))

		return nil
	} else if err != nil {
		return fmt.Errorf("failed to create deployment %q (ns %s): %w", deployment.Name, deployment.Namespace, err)
	}

	Eventually(func() bool {
		status, err := IsDeploymentReady(client, runningDeployment.Namespace, runningDeployment.Name)
		if err != nil {
			glog.V(5).Info(fmt.Sprintf(
				"deployment %s is not ready, retry in 5 seconds", runningDeployment.Name))

			return false
		}

		return status
	}, timeout, retryInterval*time.Second).Should(Equal(true), "Deployment is not ready")

	return nil
}
