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
	"k8s.io/utils/pointer"

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
				tsparams.TestTargetLabels), replica),
		containers-1,
		globalhelper.Configuration.General.TestImage)

	return deploymentStruct, nil
}

func DefineReplicaSet(name string) *v1.ReplicaSet {
	return replicaset.DefineReplicaSet(name,
		tsparams.LifecycleNamespace,
		globalhelper.Configuration.General.TestImage,
		tsparams.TestTargetLabels)
}

func DefineStatefulSet(name string) *v1.StatefulSet {
	return statefulset.DefineStatefulSet(name,
		tsparams.LifecycleNamespace,
		globalhelper.Configuration.General.TestImage,
		tsparams.TestTargetLabels)
}

func DefinePod(name string) *corev1.Pod {
	return pod.DefinePod(name, tsparams.LifecycleNamespace,
		globalhelper.Configuration.General.TestImage)
}

func DefineDaemonSetWithImagePullPolicy(name string, image string, pullPolicy corev1.PullPolicy) *v1.DaemonSet {
	return daemonset.RedefineWithImagePullPolicy(
		daemonset.DefineDaemonSet(tsparams.LifecycleNamespace, image,
			tsparams.TestTargetLabels, name), pullPolicy)
}

// WaitUntilClusterIsStable validates that all nodes are schedulable, and in ready state.
func WaitUntilClusterIsStable() error {
	Eventually(func() bool {
		isClusterReady, err := cluster.IsClusterStable(globalhelper.APIClient)
		Expect(err).ToNot(HaveOccurred())

		return isClusterReady
	}, tsparams.WaitingTime, tsparams.RetryInterval*time.Second).Should(BeTrue())

	err := nodes.WaitForNodesReady(globalhelper.APIClient,
		tsparams.WaitingTime, tsparams.RetryInterval*time.Second)

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

func CreatePersistentVolume(pv *corev1.PersistentVolume, timeout time.Duration) error {
	_, err := globalhelper.APIClient.PersistentVolumes().Create(context.Background(), pv, metav1.CreateOptions{})

	return err
}

func CreateAndWaitUntilPVCIsBound(pvc *corev1.PersistentVolumeClaim, namespace string, timeout time.Duration, pvName string) error {
	pvc, err := globalhelper.APIClient.PersistentVolumeClaims(namespace).Create(context.Background(), pvc, metav1.CreateOptions{})
	if err != nil {
		return err
	}

	Eventually(func() bool {

		status, err := isPvcBound(pvc.Name, pvc.Namespace, pvName)
		if err != nil {

			glog.V(5).Info(fmt.Sprintf(
				"pvc %s is not bound, retry in %d seconds", pvc.Name, retryInterval))

			return false
		}

		return status
	}, timeout, retryInterval*time.Second).Should(Equal(true), "pvc is not bound")

	return nil
}

func isPvcBound(pvcName string, namespace string, pvName string) (bool, error) {
	pvc, err := globalhelper.APIClient.PersistentVolumeClaims(namespace).Get(context.Background(), pvcName, metav1.GetOptions{})
	if err != nil {
		return false, err
	}

	return pvc.Status.Phase == corev1.ClaimBound && pvc.Spec.VolumeName == pvName, nil
}

func DeletePV(persistentVolume string, timeout time.Duration) error {
	err := globalhelper.APIClient.PersistentVolumes().Delete(context.Background(), persistentVolume, metav1.DeleteOptions{
		GracePeriodSeconds: pointer.Int64Ptr(0),
	})
	if err != nil {
		return fmt.Errorf("failed to delete persistent volume %w", err)
	}

	Eventually(func() bool {
		// if the pv was deleted, we will get an error.
		_, err := globalhelper.APIClient.PersistentVolumes().Get(context.Background(), persistentVolume, metav1.GetOptions{})

		return err != nil
	}, timeout, tsparams.RetryInterval*time.Second).Should(Equal(true), "PV is not removed yet.")

	return nil
}

func DeleteRunTimeClass(rtcName string) error {
	err := globalhelper.APIClient.RuntimeClasses().Delete(context.Background(), rtcName,
		metav1.DeleteOptions{GracePeriodSeconds: pointer.Int64Ptr(0)})
	if err != nil {
		return fmt.Errorf("failed to delete RunTimeClasses %w", err)
	}

	return err
}
