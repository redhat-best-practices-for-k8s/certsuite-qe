package globalhelper

import (
	"context"
	"fmt"
	"time"

	egiClients "github.com/openshift-kni/eco-goinfra/pkg/clients"
	egiDeployment "github.com/openshift-kni/eco-goinfra/pkg/deployment"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	klog "k8s.io/klog/v2"
)

var (
	deploymentPollInterval      = 10 * time.Second
	deploymentReadyProbeTimeout = time.Second
	deploymentReadyChecker      = IsDeploymentReady
)

// IsDeploymentReady checks if a deployment is ready.
func IsDeploymentReady(client *egiClients.Settings, namespace, deploymentName string) (bool, error) {
	dep, err := egiDeployment.Pull(client, deploymentName, namespace)
	if err != nil {
		return false, fmt.Errorf("failed to get deployment %q (ns %s): %w", deploymentName, namespace, err)
	}

	return dep.IsReady(deploymentReadyProbeTimeout), nil
}

// CreateAndWaitUntilDeploymentIsReady creates deployment and wait until all deployment replicas are up and running.
func CreateAndWaitUntilDeploymentIsReady(deployment *appsv1.Deployment,
	timeout time.Duration) error {
	return createAndWaitUntilDeploymentIsReady(GetEcoGoinfraClient(), deployment, timeout)
}

// createAndWaitUntilDeploymentIsReady creates deployment and wait until all deployment replicas are up and running.
func createAndWaitUntilDeploymentIsReady(client *egiClients.Settings, deployment *appsv1.Deployment,
	timeout time.Duration) error {
	runningDeployment, err := client.Deployments(deployment.Namespace).Create(
		context.TODO(),
		deployment,
		metav1.CreateOptions{})

	if k8serrors.IsAlreadyExists(err) {
		klog.V(5).Infof("deployment %s already exists", deployment.Name)

		return nil
	} else if err != nil {
		return fmt.Errorf("failed to create deployment %q (ns %s): %w", deployment.Name, deployment.Namespace, err)
	}

	// Wait until every replica is running, failing fast when the deployment cannot be scheduled.
	err = WaitForState(context.TODO(),
		fmt.Sprintf("deployment %s to run", deployment.Name),
		deploymentPollInterval, timeout,
		func(ctx context.Context) (bool, error) {
			testDeployment, err := client.Deployments(deployment.Namespace).Get(
				ctx, deployment.Name, metav1.GetOptions{})
			if err != nil {
				klog.V(5).Infof("deployment %s is not running, retry in %v: %v",
					deployment.Name, deploymentPollInterval, err)

				return false, nil
			}

			if testDeployment.Status.ReadyReplicas == *testDeployment.Spec.Replicas {
				klog.V(5).Infof("deployment %s is running", testDeployment.Name)
				klog.V(5).Infof("deployment %s status: %v", testDeployment.Name, testDeployment.Status)

				return true, nil
			}

			klog.V(5).Infof("Deployment %s conditions: %v", testDeployment.Name, testDeployment.Status.Conditions)

			if hasReplicaFailure(testDeployment) {
				return false, fmt.Errorf("deployment %s is not schedulable", testDeployment.Name)
			}

			return false, nil
		})
	if err != nil {
		return err
	}

	// Wait until the deployment reports itself ready.
	return WaitForState(context.TODO(),
		fmt.Sprintf("deployment %s to become ready", runningDeployment.Name),
		deploymentPollInterval, timeout,
		func(context.Context) (bool, error) {
			status, err := deploymentReadyChecker(client, runningDeployment.Namespace, runningDeployment.Name)
			if err != nil {
				klog.V(5).Infof("deployment %s is not ready, retrying in %v: %v",
					runningDeployment.Name, deploymentPollInterval, err)

				return false, nil
			}

			return status, nil
		})
}

// hasReplicaFailure reports whether the deployment hit a condition that no amount of polling will fix.
func hasReplicaFailure(deployment *appsv1.Deployment) bool {
	for _, condition := range deployment.Status.Conditions {
		if condition.Type == appsv1.DeploymentReplicaFailure && condition.Status == corev1.ConditionTrue {
			return true
		}
	}

	return false
}

// GetRunningDeployment returns a running deployment.
func GetRunningDeployment(namespace, deploymentName string) (*appsv1.Deployment, error) {
	return getRunningDeployment(GetEcoGoinfraClient(), namespace, deploymentName)
}

func getRunningDeployment(client *egiClients.Settings, namespace, deploymentName string) (*appsv1.Deployment, error) {
	dep, err := egiDeployment.Pull(client, deploymentName, namespace)
	if err != nil {
		return nil, fmt.Errorf("failed to get deployment %q (ns %s): %w", deploymentName, namespace, err)
	}

	return dep.Object, nil
}

func DeleteDeployment(name, namespace string) error {
	return deleteDeployment(GetEcoGoinfraClient(), name, namespace)
}

func deleteDeployment(client *egiClients.Settings, name, namespace string) error {
	return egiDeployment.NewBuilder(client, name, namespace, map[string]string{"test-app": "test"}, corev1.Container{}).Delete()
}
