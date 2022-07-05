package helper

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang/glog"
	. "github.com/onsi/gomega"

	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/cluster"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/daemonset"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/nodes"

	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/deployment"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/pod"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/replicaset"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/statefulset"

	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	tsparams "github.com/test-network-function/cnfcert-tests-verification/tests/lifecycle/parameters"
)

const retryInterval = 5

// DefineDeployment defines a deployment.
func DefineDeployment(replica int32, containers int, name string) (*v1.Deployment, error) {
	if containers < 1 {
		return nil, errors.New("invalid containers number")
	}

	deploymentStruct := globalhelper.AppendContainersToDeployment(
		deployment.RedefineWithReplicaNumber(
			deployment.DefineDeployment(
				name,
				tsparams.LifecycleNamespace,
				globalhelper.Configuration.General.TestImage,
				tsparams.TestDeploymentLabels), replica),
		containers-1,
		globalhelper.Configuration.General.TestImage)

	return deploymentStruct, nil
}

func DefineReplicaSet(name string) *v1.ReplicaSet {
	return replicaset.DefineReplicaSet(name,
		tsparams.LifecycleNamespace,
		globalhelper.Configuration.General.TestImage,
		tsparams.TestDeploymentLabels)
}

func DefineStatefulSet(name string) *v1.StatefulSet {
	return statefulset.DefineStatefulSet(name,
		tsparams.LifecycleNamespace,
		globalhelper.Configuration.General.TestImage,
		tsparams.TestDeploymentLabels)
}

func DefinePod(name string) *corev1.Pod {
	return pod.DefinePod(name, tsparams.LifecycleNamespace,
		globalhelper.Configuration.General.TestImage)
}

func DefineDaemonSetWithImagePullPolicy(name string, image string, pullPolicy corev1.PullPolicy) *v1.DaemonSet {
	return daemonset.RedefineWithImagePullPolicy(
		daemonset.DefineDaemonSet(tsparams.LifecycleNamespace, image,
			tsparams.TestDeploymentLabels, name), pullPolicy)
}

// WaitUntilClusterIsStable validates that all nodes are schedulable, and in ready state.
func WaitUntilClusterIsStable() error {
	Eventually(func() bool {
		isClusterReady, err := cluster.IsClusterStable(globalhelper.APIClient)
		Expect(err).ToNot(HaveOccurred())

		return isClusterReady
	}, tsparams.WaitingTime, tsparams.RetryInterval*time.Second).Should(BeTrue())

	err := nodes.WaitForNodesReady(globalhelper.APIClient,
		tsparams.WaitingTime, tsparams.RetryInterval)

	return err
}

// CreateAndWaitUntilReplicaSetIsReady creates replicaSet and waits until all it's replicas are ready.
func CreateAndWaitUntilReplicaSetIsReady(replicaSet *v1.ReplicaSet, timeout time.Duration) error {
	runningReplica, err := globalhelper.APIClient.ReplicaSets(replicaSet.Namespace).Create(
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

// isReplicaSetReady checks if a replicaset is ready.
func isReplicaSetReady(namespace string, replicaSetName string) (bool, error) {
	testReplicaSet, err := globalhelper.APIClient.ReplicaSets(namespace).Get(
		context.Background(),
		replicaSetName,
		metav1.GetOptions{},
	)
	if err != nil {
		return false, err
	}

	if *testReplicaSet.Spec.Replicas == testReplicaSet.Status.ReadyReplicas {
		return true, nil
	}

	return false, nil
}
