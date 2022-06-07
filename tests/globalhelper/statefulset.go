package globalhelper

import (
	"context"
	"fmt"
	"time"

	"github.com/golang/glog"
	. "github.com/onsi/gomega"
	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	replicaSetRetryIntervalSecs = 5
)

// CreateAndWaitUntilStatefulSetIsReady creates statefulset and wait until all replicas are ready.
func CreateAndWaitUntilStatefulSetIsReady(statefulSet *v1.StatefulSet, timeout time.Duration) error {
	runningReplica, err := APIClient.StatefulSets(statefulSet.Namespace).Create(
		context.Background(),
		statefulSet,
		metav1.CreateOptions{})
	if err != nil {
		return err
	}

	Eventually(func() bool {
		status, err := isStatefulSetReady(runningReplica.Namespace, runningReplica.Name)
		if err != nil {
			glog.V(5).Info(fmt.Sprintf(
				"statefulSet %s is not ready, retry in %d seconds", runningReplica.Name, replicaSetRetryIntervalSecs))

			return false
		}

		return status
	}, timeout, replicaSetRetryIntervalSecs*time.Second).Should(Equal(true), "statefulSet is not ready")

	return nil
}

func isStatefulSetReady(namespace string, statefulSetName string) (bool, error) {
	testStatefulSet, err := APIClient.StatefulSets(namespace).Get(
		context.Background(),
		statefulSetName,
		metav1.GetOptions{},
	)
	if err != nil {
		return false, err
	}

	if testStatefulSet.Status.ReadyReplicas > 0 {
		if testStatefulSet.Status.Replicas == testStatefulSet.Status.ReadyReplicas {
			return true, nil
		}
	}

	return false, nil
}
