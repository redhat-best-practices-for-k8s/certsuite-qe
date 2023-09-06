package globalhelper

import (
	"context"
	"fmt"
	"time"

	"github.com/golang/glog"
	. "github.com/onsi/gomega"
	appsv1 "k8s.io/api/apps/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// CreateAndWaitUntilReplicaSetIsReady creates replicaSet and waits until all it's replicas are ready.
func CreateAndWaitUntilReplicaSetIsReady(replicaSet *appsv1.ReplicaSet, timeout time.Duration) error {
	runningReplica, err := GetAPIClient().ReplicaSets(replicaSet.Namespace).Create(context.TODO(),
		replicaSet, metav1.CreateOptions{})
	if k8serrors.IsAlreadyExists(err) {
		glog.V(5).Info(fmt.Sprintf("replicaSet %s already exists", replicaSet.Name))
	} else if err != nil {
		return fmt.Errorf("failed to create replicaSet %q (ns %s): %w", replicaSet.Name, replicaSet.Namespace, err)
	}

	Eventually(func() bool {
		status, err := isReplicaSetReady(runningReplica.Namespace, runningReplica.Name)
		if err != nil {
			glog.V(5).Info(fmt.Sprintf(
				"replicaSet %s is not ready, retry in %d seconds", runningReplica.Name, retryInterval))

			return false
		}

		return status
	}, timeout, retryInterval*time.Second).Should(Equal(true), "replicaSet is not ready")

	return nil
}

// isReplicaSetReady checks if a replicaset is ready.
func isReplicaSetReady(namespace, replicaSetName string) (bool, error) {
	testReplicaSet, err := GetAPIClient().ReplicaSets(namespace).Get(
		context.TODO(),
		replicaSetName,
		metav1.GetOptions{},
	)
	if err != nil {
		return false, err
	}

	if *testReplicaSet.Spec.Replicas == testReplicaSet.Status.AvailableReplicas {
		return true, nil
	}

	return false, nil
}
