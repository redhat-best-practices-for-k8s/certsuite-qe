package helper

import (
	"context"
	"errors"
	"time"

	"github.com/test-network-function/cnfcert-tests-verification/tests/lifecycle/parameters"
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

	. "github.com/onsi/gomega"
)

// DefineDeployment defines a deployment.
func DefineDeployment(replica int32, containers int, name string) (*v1.Deployment, error) {
	if containers < 1 {
		return nil, errors.New("invalid containers number")
	}

	deploymentStruct := globalhelper.AppendContainersToDeployment(
		deployment.RedefineWithReplicaNumber(
			deployment.DefineDeployment(
				name,
				parameters.LifecycleNamespace,
				globalhelper.Configuration.General.TestImage,
				parameters.TestDeploymentLabels), replica),
		containers-1,
		globalhelper.Configuration.General.TestImage)

	return deploymentStruct, nil
}

func DefineReplicaSet(name string) *v1.ReplicaSet {
	return replicaset.DefineReplicaSet(name,
		parameters.LifecycleNamespace,
		globalhelper.Configuration.General.TestImage,
		parameters.TestDeploymentLabels)
}

func DefineStatefulSet(name string) *v1.StatefulSet {
	return statefulset.DefineStatefulSet(name,
		parameters.LifecycleNamespace,
		globalhelper.Configuration.General.TestImage,
		parameters.TestDeploymentLabels)
}

func DefinePod(name string) *corev1.Pod {
	return pod.DefinePod(name, parameters.LifecycleNamespace,
		globalhelper.Configuration.General.TestImage)
}

// EnableMasterScheduling enables/disables master nodes scheduling.
func EnableMasterScheduling(scheduleable bool) error {
	scheduler, err := globalhelper.APIClient.ConfigV1Interface.Schedulers().Get(
		context.TODO(), "cluster", metav1.GetOptions{})
	if err != nil {
		return err
	}

	scheduler.Spec.MastersSchedulable = scheduleable
	_, err = globalhelper.APIClient.ConfigV1Interface.Schedulers().Update(context.TODO(),
		scheduler, metav1.UpdateOptions{})

	return err
}

func DefineDaemonSetWithImagePullPolicy(name string, image string, pullPolicy corev1.PullPolicy) *v1.DaemonSet {
	return daemonset.RedefineWithImagePullPolicy(
		daemonset.DefineDaemonSet(parameters.LifecycleNamespace, image,
			parameters.TestDeploymentLabels, name), pullPolicy)
}

// WaitUntilClusterIsStable validates that all nodes are schedulable, and in ready state.
func WaitUntilClusterIsStable() error {
	Eventually(func() bool {
		isClusterReady, err := cluster.IsClusterStable(globalhelper.APIClient)
		Expect(err).ToNot(HaveOccurred())

		return isClusterReady
	}, parameters.WaitingTime, parameters.RetryInterval*time.Second).Should(BeTrue())

	err := nodes.WaitForNodesReady(globalhelper.APIClient,
		parameters.WaitingTime, parameters.RetryInterval)

	return err
}
