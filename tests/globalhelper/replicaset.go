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

// CreateAndWaitUntilReplicaSetIsReady creates replicaSet and wait until all replicas are ready.
func CreateAndWaitUntilReplicaSetIsReady(replicaSet *v1.ReplicaSet, timeout time.Duration) error {
	runningReplica, err := APIClient.ReplicaSets(replicaSet.Namespace).Create(
		context.Background(),
		replicaSet,
		metav1.CreateOptions{})
	if err != nil {
		return err
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

func isReplicaSetReady(namespace string, replicaSetName string) (bool, error) {
	testReplicaSet, err := APIClient.ReplicaSets(namespace).Get(
		context.Background(),
		replicaSetName,
		metav1.GetOptions{},
	)
	if err != nil {
		return false, err
	}

	if testReplicaSet.Status.ReadyReplicas > 0 {
		if testReplicaSet.Status.Replicas == testReplicaSet.Status.ReadyReplicas {
			return true, nil
		}
	}

	return false, nil
}
