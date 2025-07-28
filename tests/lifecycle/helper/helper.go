package helper

import (
	"errors"
	"fmt"
	"time"

	. "github.com/onsi/gomega"

	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/utils/cluster"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/utils/daemonset"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/utils/nodes"

	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/globalhelper"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/utils/deployment"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/utils/pod"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/utils/replicaset"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/utils/statefulset"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"

	tsparams "github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/lifecycle/parameters"
)

// DefineDeployment defines a deployment.
func DefineDeployment(replica int32, containers int, name, namespace string) (*appsv1.Deployment, error) {
	if containers < 1 {
		return nil, errors.New("invalid containers number")
	}

	deploymentStruct := deployment.DefineDeployment(name, namespace,
		tsparams.SampleWorkloadImage, tsparams.TestTargetLabels)

	globalhelper.AppendContainersToDeployment(deploymentStruct, containers-1, tsparams.SampleWorkloadImage)
	deployment.RedefineWithReplicaNumber(deploymentStruct, replica)

	return deploymentStruct, nil
}

func DefineReplicaSet(name, namespace string) *appsv1.ReplicaSet {
	return replicaset.DefineReplicaSet(name,
		namespace,
		tsparams.SampleWorkloadImage,
		tsparams.TestTargetLabels)
}

func DefineStatefulSet(name, namespace string) *appsv1.StatefulSet {
	return statefulset.DefineStatefulSet(name,
		namespace,
		tsparams.SampleWorkloadImage,
		tsparams.TestTargetLabels)
}

func DefinePod(name, namespace string) *corev1.Pod {
	return pod.DefinePod(name, namespace,
		tsparams.SampleWorkloadImage, tsparams.TestTargetLabels)
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
