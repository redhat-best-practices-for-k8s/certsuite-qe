package globalhelper

import (
	"context"
	"fmt"
	"time"

	. "github.com/onsi/gomega"

	"github.com/golang/glog"
	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// CreateAndWaitUntilStatefulSetIsReady creates statefulset and waits until all it's replicas are ready.
func CreateAndWaitUntilStatefulSetIsReady(statefulSet *v1.StatefulSet, timeout time.Duration) error {
	runningStatefulSet, err := APIClient.StatefulSets(statefulSet.Namespace).Create(
		context.Background(),
		statefulSet,
		metav1.CreateOptions{})
	if err != nil {
		return err
	}

	Eventually(func() bool {
		status, err := isStatefulSetReady(runningStatefulSet.Namespace, runningStatefulSet.Name)
		if err != nil {
			glog.V(5).Info(fmt.Sprintf(
				"statefulSet %s is not ready, retry in %d seconds", runningStatefulSet.Name, retryInterval))

			return false
		}

		return status
	}, timeout, retryInterval*time.Second).Should(Equal(true), "statefulSet is not ready")

	return nil
}

// isStatefulSetReady checks if a statefulset is ready.
func isStatefulSetReady(namespace string, statefulSetName string) (bool, error) {
	testStatefulSet, err := APIClient.StatefulSets(namespace).Get(
		context.Background(),
		statefulSetName,
		metav1.GetOptions{},
	)
	if err != nil {
		return false, err
	}

	if *testStatefulSet.Spec.Replicas == testStatefulSet.Status.ReadyReplicas {
		return true, nil
	}

	return false, nil
}
