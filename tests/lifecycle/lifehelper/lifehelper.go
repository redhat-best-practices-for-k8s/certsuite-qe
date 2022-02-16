package lifehelper

import (
	"context"
	"fmt"
	"time"

	"github.com/golang/glog"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/lifecycle/lifeparameters"
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
func DefineDeployment(replica int32, containers int, name string) *v1.Deployment {
	deploymentStruct := globalhelper.AppendContainersToDeployment(
		deployment.RedefineWithReplicaNumber(
			deployment.DefineDeployment(
				name,
				lifeparameters.LifecycleNamespace,
				globalhelper.Configuration.General.TnfImage,
				lifeparameters.TestDeploymentLabels), replica),
		containers,
		globalhelper.Configuration.General.TnfImage)

	return deploymentStruct
}

// RemoveterminationGracePeriod removes terminationGracePeriodSeconds field in a deployment.
func RemoveterminationGracePeriod(deploymentStruct *v1.Deployment) *v1.Deployment {
	return deployment.RedefineWithTerminationGracePeriod(deploymentStruct, nil)
}

func DefineReplicaSet(name string) *v1.ReplicaSet {
	return replicaset.DefineReplicaSet(name,
		lifeparameters.LifecycleNamespace,
		globalhelper.Configuration.General.TnfImage,
		lifeparameters.TestDeploymentLabels)
}

func DefineStatefulSet(name string) *v1.StatefulSet {
	return statefulset.DefineStatefulSet(name,
		lifeparameters.LifecycleNamespace,
		globalhelper.Configuration.General.TnfImage,
		lifeparameters.TestDeploymentLabels)
}

func DefindPod(name string) *corev1.Pod {
	return pod.DefinePod(name, lifeparameters.LifecycleNamespace,
		globalhelper.Configuration.General.TnfImage)
}

// CreateAndWaitUntilStatefulSetIsReady creates statefulset and wait until all replicas are ready.
func CreateAndWaitUntilStatefulSetIsReady(statefulSet *v1.StatefulSet, timeout time.Duration) error {
	runningReplica, err := globalhelper.APIClient.StatefulSets(statefulSet.Namespace).Create(
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
				"statefulSet %s is not ready, retry in %d seconds", runningReplica.Name, lifeparameters.RetryInterval))

			return false
		}

		return status
	}, timeout, lifeparameters.RetryInterval*time.Second).Should(Equal(true), "statefulSet is not ready")

	return nil
}

func isStatefulSetReady(namespace string, statefulSetName string) (bool, error) {
	testStatefulSet, err := globalhelper.APIClient.StatefulSets(namespace).Get(
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

// CreateAndWaitUntilReplicaSetIsReady creates replicaSet and wait until all replicas are ready.
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
				"replicaSet %s is not ready, retry in %d seconds", runningReplica.Name, lifeparameters.RetryInterval))

			return false
		}

		return status
	}, timeout, lifeparameters.RetryInterval*time.Second).Should(Equal(true), "replicaSet is not ready")

	return nil
}

func isReplicaSetReady(namespace string, replicaSetName string) (bool, error) {
	testReplicaSet, err := globalhelper.APIClient.ReplicaSets(namespace).Get(
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

// CreateAndWaitUntilPodIsReady create and wait until pod is in a "Running" phase.
func CreateAndWaitUntilPodIsReady(pod *corev1.Pod, timeout time.Duration) error {
	pod, err := globalhelper.APIClient.Pods(pod.Namespace).Create(
		context.Background(),
		pod,
		metav1.CreateOptions{})
	if err != nil {
		return err
	}

	Eventually(func() bool {
		status, err := isPodReady(pod.Namespace, pod.Name)
		if err != nil {

			glog.V(5).Info(fmt.Sprintf(
				"deployment %s is not ready, retry in %d seconds", pod.Name, lifeparameters.RetryInterval))

			return false
		}

		return status
	}, timeout, lifeparameters.RetryInterval*time.Second).Should(Equal(true), "Deployment is not ready")

	return nil
}

func isPodReady(namespace string, podName string) (bool, error) {
	pod, err := globalhelper.APIClient.Pods(namespace).Get(
		context.Background(),
		podName,
		metav1.GetOptions{},
	)

	if err != nil {
		return false, err
	}

	if pod.Status.Phase == "Running" {
		return true, nil
	}

	return false, nil
}
