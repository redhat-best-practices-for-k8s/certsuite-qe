package globalhelper

import (
	"context"
	"fmt"
	"time"

	"github.com/golang/glog"
	appsv1 "k8s.io/api/apps/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	appsv1Typed "k8s.io/client-go/kubernetes/typed/apps/v1"

	. "github.com/onsi/gomega"
)

// CreateAndWaitUntilDaemonSetIsReady creates daemonSet and waits until all pods are up and running.
func CreateAndWaitUntilDaemonSetIsReady(client appsv1Typed.AppsV1Interface, daemonSet *appsv1.DaemonSet, timeout time.Duration) error {
	runningDaemonSet, err := client.DaemonSets(daemonSet.Namespace).Create(
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
			glog.Fatal(fmt.Sprintf(
				"daemonset %s is not ready, retry in 5 seconds", runningDaemonSet.Name))

			return false
		}

		return status
	}, timeout, retryInterval*time.Second).Should(Equal(true), "DaemonSet is not ready")

	return nil
}

func isDaemonSetReady(client appsv1Typed.AppsV1Interface, namespace string, name string) (bool, error) {
	daemonSet, err := client.DaemonSets(namespace).Get(
		context.TODO(),
		name,
		metav1.GetOptions{},
	)
	if err != nil {
		return false, err
	}

	if daemonSet.Status.NumberAvailable == daemonSet.Status.DesiredNumberScheduled && daemonSet.Status.NumberUnavailable == 0 {
		return true, nil
	}

	return false, nil
}
