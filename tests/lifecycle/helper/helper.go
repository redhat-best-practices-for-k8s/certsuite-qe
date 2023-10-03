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

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"

	tsparams "github.com/test-network-function/cnfcert-tests-verification/tests/lifecycle/parameters"
	storagev1 "k8s.io/api/storage/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
)

const retryInterval = 5

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
		isClusterReady, err := cluster.IsClusterStable(globalhelper.GetAPIClient())
		Expect(err).ToNot(HaveOccurred())

		return isClusterReady
	}, tsparams.WaitingTime, tsparams.RetryInterval*time.Second).Should(BeTrue())

	err := nodes.WaitForNodesReady(globalhelper.GetAPIClient(),
		tsparams.WaitingTime, tsparams.RetryInterval*time.Second)
	if err != nil {
		return fmt.Errorf("failed to wait for node to become ready: %w", err)
	}

	return nil
}

func CreatePersistentVolume(persistentVolume *corev1.PersistentVolume, timeout time.Duration) error {
	_, err := globalhelper.GetAPIClient().PersistentVolumes().Create(context.TODO(), persistentVolume, metav1.CreateOptions{})
	if k8serrors.IsAlreadyExists(err) {
		glog.V(5).Info(fmt.Sprintf("persistent volume %s already created", persistentVolume.Name))
	} else if err != nil {
		return fmt.Errorf("failed to create persistent volume: %w", err)
	}

	return nil
}

func CreateAndWaitUntilPVCIsBound(pvc *corev1.PersistentVolumeClaim, namespace string, timeout time.Duration, pvName string) error {
	pvc, err := globalhelper.GetAPIClient().PersistentVolumeClaims(namespace).Create(context.TODO(), pvc, metav1.CreateOptions{})
	if k8serrors.IsAlreadyExists(err) {
		glog.V(5).Info(fmt.Sprintf("persistent volume claim %s already created", pvc.Name))
	} else if err != nil {
		return fmt.Errorf("failed to create persistent volume claim: %w", err)
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
	pvc, err := globalhelper.GetAPIClient().PersistentVolumeClaims(namespace).Get(context.TODO(), pvcName, metav1.GetOptions{})
	if err != nil {
		return false, err
	}

	return pvc.Status.Phase == corev1.ClaimBound && pvc.Spec.VolumeName == pvName, nil
}

func DeletePV(persistentVolume string, timeout time.Duration) error {
	err := globalhelper.GetAPIClient().PersistentVolumes().Delete(context.TODO(), persistentVolume, metav1.DeleteOptions{
		GracePeriodSeconds: ptr.To[int64](0),
	})
	if err != nil {
		return fmt.Errorf("failed to delete persistent volume %w", err)
	}

	Eventually(func() bool {
		// if the pv was deleted, we will get an error.
		_, err := globalhelper.GetAPIClient().PersistentVolumes().Get(context.TODO(), persistentVolume, metav1.GetOptions{})

		return err != nil
	}, timeout, tsparams.RetryInterval*time.Second).Should(Equal(true), "PV is not removed yet.")

	return nil
}

func DeleteRunTimeClass(rtcName string) error {
	err := globalhelper.GetAPIClient().RuntimeClasses().Delete(context.TODO(), rtcName,
		metav1.DeleteOptions{GracePeriodSeconds: ptr.To[int64](0)})
	if err != nil {
		return fmt.Errorf("failed to delete RunTimeClasses %w", err)
	}

	return nil
}

func CreateStorageClass(storageClassName string, defaultSC bool) error {
	storageClassTemplate := storagev1.StorageClass{
		ObjectMeta: metav1.ObjectMeta{
			Name: storageClassName,
		},
		Provisioner: "kubernetes.io/no-provisioner",
	}

	// Set the storageclass as default if needed.
	if defaultSC {
		storageClassTemplate.Annotations = map[string]string{
			"storageclass.kubernetes.io/is-default-class": "true",
		}
	}

	_, err := globalhelper.GetAPIClient().K8sClient.StorageV1().StorageClasses().Create(context.Background(),
		&storageClassTemplate, metav1.CreateOptions{})

	if k8serrors.IsAlreadyExists(err) {
		glog.V(5).Info(fmt.Sprintf("storageclass %s already installed", storageClassName))

		return nil
	}

	return err
}

func DeleteStorageClass(storageClassName string) error {
	err := globalhelper.GetAPIClient().K8sClient.StorageV1().StorageClasses().Delete(context.Background(),
		storageClassName, metav1.DeleteOptions{GracePeriodSeconds: ptr.To[int64](0)})

	if k8serrors.IsNotFound(err) {
		glog.V(5).Info(fmt.Sprintf("storageclass %s already deleted", storageClassName))

		return nil
	} else if err != nil {
		return fmt.Errorf("failed to delete storageclass %w", err)
	}

	return nil
}
