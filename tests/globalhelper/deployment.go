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

		// Get detailed pod and replicaset information for debugging
		printDetailedDeploymentDebugInfo(client, testDeployment)

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

// printDetailedDeploymentDebugInfo provides comprehensive debugging information for deployment issues.
func printDetailedDeploymentDebugInfo(client *egiClients.Settings, deployment *appsv1.Deployment) {
	fmt.Printf("=== DEBUG: Deployment %s in namespace %s ===\n", deployment.Name, deployment.Namespace)

	// 1. Print deployment status details
	fmt.Printf("Deployment Status: Replicas=%d, Ready=%d, Available=%d, Updated=%d\n",
		deployment.Status.Replicas,
		deployment.Status.ReadyReplicas,
		deployment.Status.AvailableReplicas,
		deployment.Status.UpdatedReplicas)

	// 2. Get and print ReplicaSet information
	replicaSets, err := client.ReplicaSets(deployment.Namespace).List(context.TODO(), metav1.ListOptions{
		LabelSelector: metav1.FormatLabelSelector(deployment.Spec.Selector),
	})
	if err != nil {
		fmt.Printf("ERROR: Failed to get ReplicaSets: %v\n", err)
	} else {
		for _, rs := range replicaSets.Items {
			fmt.Printf("ReplicaSet %s: Replicas=%d, Ready=%d, Available=%d\n",
				rs.Name, rs.Status.Replicas, rs.Status.ReadyReplicas, rs.Status.AvailableReplicas)

			// Print ReplicaSet conditions
			for _, condition := range rs.Status.Conditions {
				if condition.Status == corev1.ConditionFalse {
					fmt.Printf("  RS Condition: %s=%s, Reason=%s, Message=%s\n",
						condition.Type, condition.Status, condition.Reason, condition.Message)
				}
			}
		}
	}

	// 3. Get and print Pod information with detailed status
	pods, err := client.Pods(deployment.Namespace).List(context.TODO(), metav1.ListOptions{
		LabelSelector: metav1.FormatLabelSelector(deployment.Spec.Selector),
	})
	if err != nil {
		fmt.Printf("ERROR: Failed to get Pods: %v\n", err)
	} else {
		fmt.Printf("Found %d pods for deployment\n", len(pods.Items))

		for _, pod := range pods.Items {
			fmt.Printf("Pod %s: Phase=%s, Ready=%s\n", pod.Name, pod.Status.Phase, getPodReadyStatus(&pod))

			// Print container statuses
			for _, containerStatus := range pod.Status.ContainerStatuses {
				fmt.Printf("  Container %s: Ready=%t, RestartCount=%d\n",
					containerStatus.Name, containerStatus.Ready, containerStatus.RestartCount)

				// Print detailed container state
				if containerStatus.State.Waiting != nil {
					fmt.Printf("    Waiting: Reason=%s, Message=%s\n",
						containerStatus.State.Waiting.Reason, containerStatus.State.Waiting.Message)
				}

				if containerStatus.State.Terminated != nil {
					fmt.Printf("    Terminated: Reason=%s, Message=%s, ExitCode=%d\n",
						containerStatus.State.Terminated.Reason,
						containerStatus.State.Terminated.Message,
						containerStatus.State.Terminated.ExitCode)
				}
			}

			// Print pod conditions
			for _, condition := range pod.Status.Conditions {
				if condition.Status == corev1.ConditionFalse {
					fmt.Printf("  Pod Condition: %s=%s, Reason=%s, Message=%s\n",
						condition.Type, condition.Status, condition.Reason, condition.Message)
				}
			}
		}
	}

	// 4. Get and print recent events related to this deployment
	events, err := client.Events(deployment.Namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		fmt.Printf("ERROR: Failed to get Events: %v\n", err)
	} else {
		fmt.Printf("Recent Events (last 5 minutes):\n")
		recentTime := time.Now().Add(-5 * time.Minute)
		eventCount := 0

		for _, event := range events.Items {
			if event.FirstTimestamp.After(recentTime) {
				// Filter events related to our deployment, replicasets, or pods
				if containsDeploymentRelatedObject(event, deployment.Name) {
					fmt.Printf("  %s [%s]: %s - %s\n",
						event.FirstTimestamp.Format("15:04:05"),
						event.Type,
						event.Reason,
						event.Message)

					eventCount++
				}
			}
		}

		if eventCount == 0 {
			fmt.Printf("  No recent events found for this deployment\n")
		}
	}

	fmt.Printf("=== END DEBUG INFO ===\n")
}

// getPodReadyStatus returns a human-readable ready status for a pod.
func getPodReadyStatus(pod *corev1.Pod) string {
	for _, condition := range pod.Status.Conditions {
		if condition.Type == corev1.PodReady {
			return string(condition.Status)
		}
	}

	return "Unknown"
}

// containsDeploymentRelatedObject checks if an event is related to our deployment.
func containsDeploymentRelatedObject(event corev1.Event, deploymentName string) bool {
	objectName := event.InvolvedObject.Name

	// Check if event is about the deployment itself
	if objectName == deploymentName {
		return true
	}

	// Check if event is about a replicaset owned by this deployment
	// ReplicaSet names typically follow pattern: deploymentname-<hash>
	if len(objectName) > len(deploymentName) &&
		objectName[:len(deploymentName)] == deploymentName &&
		objectName[len(deploymentName)] == '-' {
		return true
	}

	return false
}
