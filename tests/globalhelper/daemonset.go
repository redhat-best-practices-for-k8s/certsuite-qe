package globalhelper

import (
	"context"
	"fmt"
	"time"

	"github.com/golang/glog"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	appsv1Typed "k8s.io/client-go/kubernetes/typed/apps/v1"
	corev1Typed "k8s.io/client-go/kubernetes/typed/core/v1"

	. "github.com/onsi/gomega"
)

func CreateAndWaitUntilDaemonSetIsReady(daemonSet *appsv1.DaemonSet, timeout time.Duration) error {
	return createAndWaitUntilDaemonSetIsReady(GetAPIClient().K8sClient.AppsV1(), GetAPIClient().K8sClient.CoreV1(), daemonSet, timeout)
}

// CreateAndWaitUntilDaemonSetIsReady creates daemonSet and waits until all pods are up and running.
func createAndWaitUntilDaemonSetIsReady(appsClient appsv1Typed.AppsV1Interface,
	coreClient corev1Typed.CoreV1Interface,
	daemonSet *appsv1.DaemonSet, timeout time.Duration) error {
	runningDaemonSet, err := appsClient.DaemonSets(daemonSet.Namespace).Create(
		context.TODO(), daemonSet, metav1.CreateOptions{})
	if k8serrors.IsAlreadyExists(err) {
		glog.V(5).Info(fmt.Sprintf("daemonset %s already exists", daemonSet.Name))

		return nil
	} else if err != nil {
		return fmt.Errorf("failed to create daemonset: %w", err)
	}

	Eventually(func() bool {
		status, err := isDaemonSetReady(appsClient, coreClient, runningDaemonSet.Namespace, runningDaemonSet.Name)
		if err != nil {
			glog.Errorf(
				"daemonset %s is not ready, retry in 5 seconds", runningDaemonSet.Name)

			return false
		}

		return status
	}, timeout, retryInterval*time.Second).Should(Equal(true), "DaemonSet is not ready")

	return nil
}

func isDaemonSetReady(client appsv1Typed.AppsV1Interface, coreClient corev1Typed.CoreV1Interface, namespace string, name string) (bool, error) {
	daemonSet, err := client.DaemonSets(namespace).Get(
		context.TODO(),
		name,
		metav1.GetOptions{},
	)
	if err != nil {
		return false, err
	}

	// Get number of nodes and compare with the number of scheduled pods
	numNodes := GetNumberOfNodes(coreClient)
	if daemonSet.Status.DesiredNumberScheduled == int32(numNodes) &&
		daemonSet.Status.NumberReady == daemonSet.Status.DesiredNumberScheduled {
		return true, nil
	}

	return false, nil
}

func GetDaemonSetPullPolicy(ds *appsv1.DaemonSet) (corev1.PullPolicy, error) {
	return getDaemonSetPullPolicy(ds, GetAPIClient().K8sClient.AppsV1())
}

func getDaemonSetPullPolicy(daemonset *appsv1.DaemonSet, client appsv1Typed.AppsV1Interface) (corev1.PullPolicy, error) {
	runningDaemonSet, err := client.DaemonSets(daemonset.Namespace).Get(
		context.TODO(),
		daemonset.Name,
		metav1.GetOptions{},
	)
	if err != nil {
		glog.Errorf("failed to get daemonset %s: %v", daemonset.Name, err)

		return "", err
	}

	return runningDaemonSet.Spec.Template.Spec.Containers[0].ImagePullPolicy, nil
}

func GetRunningDaemonset(ds *appsv1.DaemonSet) (*appsv1.DaemonSet, error) {
	return getRunningDaemonset(ds, GetAPIClient().K8sClient.AppsV1())
}

func getRunningDaemonset(daemonset *appsv1.DaemonSet, client appsv1Typed.AppsV1Interface) (*appsv1.DaemonSet, error) {
	runningDaemonSet, err := client.DaemonSets(daemonset.Namespace).Get(
		context.TODO(),
		daemonset.Name,
		metav1.GetOptions{},
	)
	if err != nil {
		glog.Errorf("failed to get daemonset %s: %v", daemonset.Name, err)

		return nil, err
	}

	return runningDaemonSet, nil
}
