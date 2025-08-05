package globalhelper

import (
	"context"
	"fmt"
	"time"

	"github.com/golang/glog"
	egiClients "github.com/openshift-kni/eco-goinfra/pkg/clients"
	egiDaemonset "github.com/openshift-kni/eco-goinfra/pkg/daemonset"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	. "github.com/onsi/gomega"
)

func CreateAndWaitUntilDaemonSetIsReady(daemonSet *appsv1.DaemonSet, timeout time.Duration) error {
	return createAndWaitUntilDaemonSetIsReady(egiClients.New(""), daemonSet, timeout)
}

// CreateAndWaitUntilDaemonSetIsReady creates daemonSet and waits until all pods are up and running.
func createAndWaitUntilDaemonSetIsReady(client *egiClients.Settings,
	daemonSet *appsv1.DaemonSet, timeout time.Duration) error {
	runningDaemonSet, err := client.AppsV1Interface.DaemonSets(daemonSet.Namespace).Create(
		context.TODO(), daemonSet, metav1.CreateOptions{})
	if k8serrors.IsAlreadyExists(err) {
		glog.V(5).Info(fmt.Sprintf("daemonset %s already exists", daemonSet.Name))

		return nil
	} else if err != nil {
		return fmt.Errorf("failed to create daemonset: %w", err)
	}

	Eventually(func() bool {
		status, err := isDaemonSetReady(client, runningDaemonSet.Namespace, runningDaemonSet.Name)
		if err != nil {
			glog.Errorf(
				"daemonset %s is not ready, retry in 1 second", runningDaemonSet.Name)

			return false
		}

		return status
	}, timeout, 1*time.Second).Should(Equal(true), "DaemonSet is not ready")

	return nil
}

func isDaemonSetReady(client *egiClients.Settings, namespace string, name string) (bool, error) {
	daemonset, err := egiDaemonset.Pull(client, name, namespace)
	if err != nil {
		return false, fmt.Errorf("failed to get daemonset %q (ns %s): %w", name, namespace, err)
	}

	// Get number of nodes and compare with the number of scheduled pods
	numNodes := GetNumberOfNodes(client.CoreV1Interface)

	if daemonset.Object.Status.DesiredNumberScheduled == int32(numNodes) &&
		daemonset.Object.Status.NumberReady == daemonset.Object.Status.DesiredNumberScheduled &&
		daemonset.Object.Status.NumberAvailable == daemonset.Object.Status.DesiredNumberScheduled &&
		daemonset.Object.Status.CurrentNumberScheduled == daemonset.Object.Status.DesiredNumberScheduled &&
		daemonset.Object.Status.NumberUnavailable == 0 &&
		daemonset.Object.Status.NumberReady > 0 {
		return true, nil
	}

	return false, nil
}

func GetDaemonSetPullPolicy(ds *appsv1.DaemonSet) (corev1.PullPolicy, error) {
	return getDaemonSetPullPolicy(ds, egiClients.New(""))
}

func getDaemonSetPullPolicy(daemonset *appsv1.DaemonSet, client *egiClients.Settings) (corev1.PullPolicy, error) {
	ds, err := egiDaemonset.Pull(client, daemonset.Name, daemonset.Namespace)
	if err != nil {
		return "", fmt.Errorf("failed to get daemonset %q (ns %s): %w", daemonset.Name, daemonset.Namespace, err)
	}

	return ds.Object.Spec.Template.Spec.Containers[0].ImagePullPolicy, nil
}

func GetRunningDaemonsetByName(name, namespace string) (*appsv1.DaemonSet, error) {
	return getRunningDaemonsetByName(name, namespace, egiClients.New(""))
}

func getRunningDaemonsetByName(name, namespace string, client *egiClients.Settings) (*appsv1.DaemonSet, error) {
	ds, err := egiDaemonset.Pull(client, name, namespace)
	if err != nil {
		return nil, fmt.Errorf("failed to get daemonset %q (ns %s): %w", name, namespace, err)
	}

	return ds.Object, nil
}

func GetRunningDaemonset(ds *appsv1.DaemonSet) (*appsv1.DaemonSet, error) {
	return getRunningDaemonset(ds, egiClients.New(""))
}

func getRunningDaemonset(daemonset *appsv1.DaemonSet, client *egiClients.Settings) (*appsv1.DaemonSet, error) {
	ds, err := egiDaemonset.Pull(client, daemonset.Name, daemonset.Namespace)
	if err != nil {
		return nil, fmt.Errorf("failed to get daemonset %q (ns %s): %w", daemonset.Name, daemonset.Namespace, err)
	}

	return ds.Object, nil
}
