package helper

import (
	"errors"
	"fmt"
	"time"

	. "github.com/onsi/gomega"

	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/cluster"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/daemonset"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/nodes"

	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/deployment"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/pod"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/replicaset"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/statefulset"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"

	tsparams "github.com/test-network-function/cnfcert-tests-verification/tests/lifecycle/parameters"
)

// DefineDeployment defines a deployment.
func DefineDeployment(replica int32, containers int, name, namespace string) (*appsv1.Deployment, error) {
	if containers < 1 {
		return nil, errors.New("invalid containers number")
	}

	deploymentStruct := deployment.DefineDeployment(name, namespace,
		globalhelper.GetConfiguration().General.TestImage, tsparams.TestTargetLabels)

	globalhelper.AppendContainersToDeployment(deploymentStruct, containers-1, globalhelper.GetConfiguration().General.TestImage)
	deployment.RedefineWithReplicaNumber(deploymentStruct, replica)

	return deploymentStruct, nil
}

func DefineReplicaSet(name, namespace string) *appsv1.ReplicaSet {
	return replicaset.DefineReplicaSet(name,
		namespace,
		globalhelper.GetConfiguration().General.TestImage,
		tsparams.TestTargetLabels)
}

func DefineStatefulSet(name, namespace string) *appsv1.StatefulSet {
	return statefulset.DefineStatefulSet(name,
		namespace,
		globalhelper.GetConfiguration().General.TestImage,
		tsparams.TestTargetLabels)
}

func DefinePod(name, namespace string) *corev1.Pod {
	return pod.DefinePod(name, namespace,
		globalhelper.GetConfiguration().General.TestImage, tsparams.TestTargetLabels)
}

func DefineDaemonSetWithImagePullPolicy(name, namespace string, image string, pullPolicy corev1.PullPolicy) *appsv1.DaemonSet {
	daemonSet := daemonset.DefineDaemonSet(namespace, image, tsparams.TestTargetLabels, name)
	daemonset.RedefineWithImagePullPolicy(daemonSet, pullPolicy)

	return daemonSet
}

// WaitUntilClusterIsStable validates that all nodes are schedulable, and in ready state.
func WaitUntilClusterIsStable() error {
	Eventually(func() bool {
		isClusterReady, err := cluster.IsClusterStable(globalhelper.GetAPIClient().Nodes())
		Expect(err).ToNot(HaveOccurred())

		return isClusterReady
	}, tsparams.WaitingTime, tsparams.RetryInterval*time.Second).Should(BeTrue())

	err := nodes.WaitForNodesReady(globalhelper.GetAPIClient().Nodes(),
		tsparams.WaitingTime, tsparams.RetryInterval*time.Second)
	if err != nil {
		return fmt.Errorf("failed to wait for node to become ready: %w", err)
	}

	return nil
}
