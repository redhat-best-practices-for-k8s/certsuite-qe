package globalhelper

import (
	"context"
	"fmt"
	"time"

	"github.com/golang/glog"
	egiClients "github.com/openshift-kni/eco-goinfra/pkg/clients"
	egiDeployment "github.com/openshift-kni/eco-goinfra/pkg/deployment"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	. "github.com/onsi/gomega"
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
	return createAndWaitUntilDeploymentIsReady(egiClients.New(""), deployment, timeout)
}

// createAndWaitUntilDeploymentIsReady creates deployment and wait until all deployment replicas are up and running.
func createAndWaitUntilDeploymentIsReady(client *egiClients.Settings, deployment *appsv1.Deployment,
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

	// Check if all pods in the deployment are schedulable
	deploymentUnschedulable := false

	Eventually(func() bool {
		testDeployment, err := client.Deployments(deployment.Namespace).Get(
			context.TODO(), deployment.Name, metav1.GetOptions{})

		if err != nil {
			glog.V(5).Info(fmt.Sprintf(
				"deployment %s is not running, retry in 1 second", testDeployment.Name))

			return false
		}

		// If it is running, we can break the loop
		if testDeployment.Status.ReadyReplicas == *testDeployment.Spec.Replicas {
			glog.V(5).Info(fmt.Sprintf("deployment %s is running", testDeployment.Name))
			glog.V(5).Info(fmt.Sprintf("deployment %s status: %v", testDeployment.Name, testDeployment.Status))

			return true
		}

		// print the conditions
		fmt.Printf("Deployment %s conditions: %v\n", testDeployment.Name, testDeployment.Status.Conditions)

		for _, condition := range testDeployment.Status.Conditions {
			if condition.Type == appsv1.DeploymentReplicaFailure && condition.Status == corev1.ConditionTrue {
				deploymentUnschedulable = true

				break
			}
		}

		return deploymentUnschedulable
	}, timeout, 1*time.Second).Should(Equal(true), "Deployment is not running")

	if deploymentUnschedulable {
		return fmt.Errorf("deployment %s is not schedulable", runningDeployment.Name)
	}

	Eventually(func() bool {
		status, err := IsDeploymentReady(client, runningDeployment.Namespace, runningDeployment.Name)
		if err != nil {
			glog.V(5).Info(fmt.Sprintf(
				"deployment %s is not ready, retry in 1 second", runningDeployment.Name))

			return false
		}

		return status
	}, timeout, 1*time.Second).Should(Equal(true), "Deployment is not ready")

	return nil
}

// GetRunningDeployment returns a running deployment.
func GetRunningDeployment(namespace, deploymentName string) (*appsv1.Deployment, error) {
	return getRunningDeployment(egiClients.New(""), namespace, deploymentName)
}

func getRunningDeployment(client *egiClients.Settings, namespace, deploymentName string) (*appsv1.Deployment, error) {
	dep, err := egiDeployment.Pull(client, deploymentName, namespace)
	if err != nil {
		return nil, fmt.Errorf("failed to get deployment %q (ns %s): %w", deploymentName, namespace, err)
	}

	return dep.Object, nil
}

func DeleteDeployment(name, namespace string) error {
	return deleteDeployment(egiClients.New(""), name, namespace)
}

func deleteDeployment(client *egiClients.Settings, name, namespace string) error {
	return egiDeployment.NewBuilder(client, name, namespace, map[string]string{"test-app": "test"}, corev1.Container{}).Delete()
}
