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

// IsDeploymentReady checks if a deployment is ready.
func IsDeploymentReady(client *egiClients.Settings, namespace, deploymentName string) (bool, error) {
	dep, err := egiDeployment.Pull(client, deploymentName, namespace)
	if err != nil {
		return false, fmt.Errorf("failed to get deployment %q (ns %s): %w", deploymentName, namespace, err)
	}

	return dep.IsReady(1 * time.Second), nil
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
		klog.V(5).Info(fmt.Sprintf("deployment %s already exists", deployment.Name))

		return nil
	} else if err != nil {
		return fmt.Errorf("failed to create deployment %q (ns %s): %w", deployment.Name, deployment.Namespace, err)
	}

	// Check if all pods in the deployment are schedulable
	deploymentUnschedulable := false

	timeoutChan := time.After(timeout)

	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	running := false
	for !running && !deploymentUnschedulable {
		select {
		case <-timeoutChan:
			return fmt.Errorf("deployment %s timed out after %v waiting to become running", deployment.Name, timeout)
		case <-ticker.C:
			testDeployment, err := client.Deployments(deployment.Namespace).Get(
				context.TODO(), deployment.Name, metav1.GetOptions{})

			if err != nil {
				klog.V(5).Info(fmt.Sprintf(
					"deployment %s is not running, retry in 1 second", deployment.Name))

				continue
			}

			// If it is running, we can break the loop
			if testDeployment.Status.ReadyReplicas == *testDeployment.Spec.Replicas {
				klog.V(5).Info(fmt.Sprintf("deployment %s is running", testDeployment.Name))
				klog.V(5).Info(fmt.Sprintf("deployment %s status: %v", testDeployment.Name, testDeployment.Status))
				running = true

				break
			}

			// print the conditions
			klog.V(5).Infof("Deployment %s conditions: %v", testDeployment.Name, testDeployment.Status.Conditions)

			for _, condition := range testDeployment.Status.Conditions {
				if condition.Type == appsv1.DeploymentReplicaFailure && condition.Status == corev1.ConditionTrue {
					deploymentUnschedulable = true

					break
				}
			}
		}
	}

	if deploymentUnschedulable {
		return fmt.Errorf("deployment %s is not schedulable", runningDeployment.Name)
	}

	timeoutChan = time.After(timeout)

	ticker = time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-timeoutChan:
			return fmt.Errorf("deployment %s timed out after %v waiting to become ready", runningDeployment.Name, timeout)
		case <-ticker.C:
			status, err := IsDeploymentReady(client, runningDeployment.Namespace, runningDeployment.Name)
			if err != nil {
				klog.V(5).Infof("deployment %s is not ready, retrying in 10 seconds", runningDeployment.Name)

				continue
			}

			if status {
				return nil
			}
		}
	}
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
