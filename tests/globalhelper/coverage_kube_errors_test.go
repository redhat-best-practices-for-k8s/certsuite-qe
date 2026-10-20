//nolint:lll // Inline fake API fixtures and reactor callbacks are intentionally compact.
package globalhelper

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	nodev1 "k8s.io/api/node/v1"
	policyv1 "k8s.io/api/policy/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	k8stesting "k8s.io/client-go/testing"
)

//nolint:funlen // This test covers related API error outcomes with shared fake clients.
func TestCoverageKubernetesErrorPaths(t *testing.T) {
	gomega.RegisterTestingT(t)
	oldPVCInterval := persistentVolumePollInterval
	oldPDBInterval, oldRuntimeInterval := podDisruptionBudgetReadyPollInterval, runtimeClassPollInterval
	persistentVolumePollInterval, podDisruptionBudgetReadyPollInterval, runtimeClassPollInterval = time.Millisecond, time.Millisecond, time.Millisecond

	t.Cleanup(func() {
		persistentVolumePollInterval, podDisruptionBudgetReadyPollInterval, runtimeClassPollInterval = oldPVCInterval, oldPDBInterval, oldRuntimeInterval
	})

	k8s, _ := coverageFixture(t)
	k8s.PrependReactor("create", "serviceaccounts", func(k8stesting.Action) (bool, runtime.Object, error) {
		return true, nil, errors.New("create serviceaccount denied")
	})
	require.Error(t, createServiceAccount(k8s.CoreV1(), "sa", "ns"))
	k8s, _ = coverageFixture(t, &corev1.ServiceAccount{ObjectMeta: metav1.ObjectMeta{Name: "sa", Namespace: "ns"}})
	k8s.PrependReactor("delete", "serviceaccounts", func(k8stesting.Action) (bool, runtime.Object, error) {
		return true, nil, errors.New("delete serviceaccount denied")
	})
	require.Error(t, deleteServiceAccount(k8s.CoreV1(), "sa", "ns"))
	serviceAccounts := serviceAccountGetError{ServiceAccountInterface: k8s.CoreV1().ServiceAccounts("ns"), err: errors.New("get serviceaccount failed")}
	coreClient := coreV1WithServiceAccountOverride{CoreV1Interface: k8s.CoreV1(), serviceAccounts: serviceAccounts}
	require.Error(t, deleteServiceAccount(coreClient, "sa", "ns"))

	k8s, _ = coverageFixture(t)
	k8s.PrependReactor("create", "roles", func(k8stesting.Action) (bool, runtime.Object, error) {
		return true, nil, errors.New("create role denied")
	})
	require.Error(t, createRole(k8s.RbacV1(), DefineRole("role", "ns")))
	k8s, _ = coverageFixture(t)
	k8s.PrependReactor("create", "roles", func(k8stesting.Action) (bool, runtime.Object, error) {
		return true, &rbacv1.Role{ObjectMeta: metav1.ObjectMeta{Name: "role", Namespace: "ns"}}, k8serrors.NewAlreadyExists(rbacv1.Resource("roles"), "role")
	})
	require.NoError(t, createRole(k8s.RbacV1(), DefineRole("role", "ns")))

	k8s, _ = coverageFixture(t)
	k8s.PrependReactor("create", "rolebindings", func(k8stesting.Action) (bool, runtime.Object, error) {
		return true, nil, k8serrors.NewAlreadyExists(rbacv1.Resource("rolebindings"), "binding")
	})
	require.NoError(t, createRoleBindingWithServiceAccountSubject(k8s.RbacV1(), "binding", "role", "sa", "ns", "ns"))

	k8s, _ = coverageFixture(t)
	k8s.PrependReactor("create", "persistentvolumes", func(k8stesting.Action) (bool, runtime.Object, error) {
		return true, nil, errors.New("create pv denied")
	})
	require.Error(t, createPersistentVolume(k8s.CoreV1(), &corev1.PersistentVolume{ObjectMeta: metav1.ObjectMeta{Name: "pv"}}))
	k8s, _ = coverageFixture(t)
	k8s.PrependReactor("create", "persistentvolumes", func(k8stesting.Action) (bool, runtime.Object, error) {
		return true, nil, k8serrors.NewAlreadyExists(corev1.Resource("persistentvolumes"), "pv")
	})
	require.NoError(t, createPersistentVolume(k8s.CoreV1(), &corev1.PersistentVolume{ObjectMeta: metav1.ObjectMeta{Name: "pv"}}))
	k8s, _ = coverageFixture(t)
	k8s.PrependReactor("delete", "persistentvolumes", func(k8stesting.Action) (bool, runtime.Object, error) {
		return true, nil, errors.New("delete pv denied")
	})
	require.Error(t, deletePersistentVolume(k8s.CoreV1(), "pv", time.Millisecond))
	k8s, _ = coverageFixture(t)
	k8s.PrependReactor("delete", "persistentvolumes", func(k8stesting.Action) (bool, runtime.Object, error) {
		return true, nil, k8serrors.NewNotFound(corev1.Resource("persistentvolumes"), "pv")
	})
	require.NoError(t, deletePersistentVolume(k8s.CoreV1(), "pv", time.Millisecond))

	pvc := &corev1.PersistentVolumeClaim{ObjectMeta: metav1.ObjectMeta{Name: "pvc", Namespace: "ns"}, Spec: corev1.PersistentVolumeClaimSpec{VolumeName: "pv"}, Status: corev1.PersistentVolumeClaimStatus{Phase: corev1.ClaimBound}}
	k8s, _ = coverageFixture(t)
	k8s.PrependReactor("create", "persistentvolumeclaims", func(k8stesting.Action) (bool, runtime.Object, error) {
		return true, nil, errors.New("create pvc denied")
	})
	require.Error(t, createAndWaitUntilPVCIsBound(k8s.CoreV1(), pvc, time.Millisecond, "pv"))
	k8s, _ = coverageFixture(t, pvc)
	k8s.PrependReactor("create", "persistentvolumeclaims", func(k8stesting.Action) (bool, runtime.Object, error) {
		return true, pvc.DeepCopy(), k8serrors.NewAlreadyExists(corev1.Resource("persistentvolumeclaims"), "pvc")
	})
	k8s.PrependReactor("get", "persistentvolumeclaims", func(k8stesting.Action) (bool, runtime.Object, error) { return true, pvc.DeepCopy(), nil })
	require.NoError(t, createAndWaitUntilPVCIsBound(k8s.CoreV1(), pvc, time.Second, "pv"))
	k8s, _ = coverageFixture(t, pvc)
	getOnce := true

	k8s.PrependReactor("get", "persistentvolumeclaims", func(k8stesting.Action) (bool, runtime.Object, error) {
		if getOnce {
			getOnce = false

			return true, nil, errors.New("transient pvc get failure")
		}

		return true, pvc.DeepCopy(), nil
	})
	require.NoError(t, createAndWaitUntilPVCIsBound(k8s.CoreV1(), pvc, time.Second, "pv"))
	k8s.PrependReactor("delete", "persistentvolumeclaims", func(k8stesting.Action) (bool, runtime.Object, error) {
		return true, nil, errors.New("delete pvc denied")
	})
	require.Error(t, deletePersistentVolumeClaim(k8s.CoreV1(), pvc))

	k8s, _ = coverageFixture(t)
	k8s.PrependReactor("create", "poddisruptionbudgets", func(k8stesting.Action) (bool, runtime.Object, error) {
		return true, nil, errors.New("create pdb denied")
	})
	require.Error(t, CreatePodDisruptionBudget(&policyv1.PodDisruptionBudget{ObjectMeta: metav1.ObjectMeta{Name: "pdb", Namespace: "ns"}}, time.Millisecond))
	k8s, _ = coverageFixture(t)
	k8s.PrependReactor("get", "poddisruptionbudgets", func(k8stesting.Action) (bool, runtime.Object, error) { return true, nil, errors.New("get pdb denied") })
	_, err := IsPodDisruptionBudgetCreated("pdb", "ns")
	require.Error(t, err)
	pdb := &policyv1.PodDisruptionBudget{ObjectMeta: metav1.ObjectMeta{Name: "pdb", Namespace: "ns", UID: "uid"}}
	k8s, _ = coverageFixture(t)
	k8s.PrependReactor("create", "poddisruptionbudgets", func(k8stesting.Action) (bool, runtime.Object, error) {
		return true, pdb.DeepCopy(), k8serrors.NewAlreadyExists(policyv1.Resource("poddisruptionbudgets"), "pdb")
	})
	pdbGetOnce := true

	k8s.PrependReactor("get", "poddisruptionbudgets", func(k8stesting.Action) (bool, runtime.Object, error) {
		if pdbGetOnce {
			pdbGetOnce = false

			return true, nil, errors.New("transient pdb get failure")
		}

		return true, pdb.DeepCopy(), nil
	})
	require.NoError(t, CreatePodDisruptionBudget(pdb, time.Second))

	k8s, _ = coverageFixture(t)
	k8s.PrependReactor("create", "runtimeclasses", func(k8stesting.Action) (bool, runtime.Object, error) {
		return true, nil, errors.New("create runtimeclass denied")
	})
	require.Error(t, createRunTimeClass(k8s, &nodev1.RuntimeClass{ObjectMeta: metav1.ObjectMeta{Name: "rtc"}}))
	rtc := &nodev1.RuntimeClass{ObjectMeta: metav1.ObjectMeta{Name: "rtc", UID: "uid"}}
	k8s, _ = coverageFixture(t)
	k8s.PrependReactor("create", "runtimeclasses", func(k8stesting.Action) (bool, runtime.Object, error) {
		return true, rtc.DeepCopy(), k8serrors.NewAlreadyExists(nodev1.Resource("runtimeclasses"), "rtc")
	})
	k8s.PrependReactor("get", "runtimeclasses", func(k8stesting.Action) (bool, runtime.Object, error) { return true, rtc.DeepCopy(), nil })
	require.NoError(t, createRunTimeClass(k8s, rtc))
	k8s, _ = coverageFixture(t)
	rtcGetOnce := true

	k8s.PrependReactor("get", "runtimeclasses", func(k8stesting.Action) (bool, runtime.Object, error) {
		if rtcGetOnce {
			rtcGetOnce = false

			return true, nil, errors.New("transient runtimeclass create check failure")
		}

		return true, rtc.DeepCopy(), nil
	})
	require.NoError(t, createRunTimeClass(k8s, rtc))
	k8s, _ = coverageFixture(t)
	k8s.PrependReactor("delete", "runtimeclasses", func(k8stesting.Action) (bool, runtime.Object, error) {
		return true, nil, errors.New("delete runtimeclass denied")
	})
	require.Error(t, deleteRunTimeClass(k8s, &nodev1.RuntimeClass{ObjectMeta: metav1.ObjectMeta{Name: "rtc"}}))
	k8s, _ = coverageFixture(t, rtc)
	rtcGetCount := 0

	k8s.PrependReactor("get", "runtimeclasses", func(k8stesting.Action) (bool, runtime.Object, error) {
		rtcGetCount++
		if rtcGetCount == 1 {
			return true, nil, errors.New("transient runtimeclass get failure")
		}

		return false, nil, nil
	})
	require.NoError(t, deleteRunTimeClass(k8s, rtc))

	k8s, _ = coverageFixture(t)
	k8s.PrependReactor("delete", "storageclasses", func(k8stesting.Action) (bool, runtime.Object, error) {
		return true, nil, errors.New("delete storageclass denied")
	})
	require.Error(t, deleteStorageClass(k8s.StorageV1(), "sc"))
	k8s, _ = coverageFixture(t)
	k8s.PrependReactor("create", "storageclasses", func(k8stesting.Action) (bool, runtime.Object, error) {
		return true, nil, errors.New("create storageclass denied")
	})
	require.Error(t, createStorageClass(k8s.StorageV1(), "sc", false))
	k8s, _ = coverageFixture(t)
	k8s.PrependReactor("create", "storageclasses", func(k8stesting.Action) (bool, runtime.Object, error) {
		return true, nil, k8serrors.NewAlreadyExists(corev1.Resource("storageclasses"), "sc")
	})
	require.NoError(t, createStorageClass(k8s.StorageV1(), "sc", false))
}

type coreV1WithServiceAccountOverride struct {
	typedcorev1.CoreV1Interface
	serviceAccounts typedcorev1.ServiceAccountInterface
}

//nolint:ireturn // This fake must implement the typed client interface.
func (c coreV1WithServiceAccountOverride) ServiceAccounts(string) typedcorev1.ServiceAccountInterface {
	return c.serviceAccounts
}

type serviceAccountGetError struct {
	typedcorev1.ServiceAccountInterface
	err error
}

func (c serviceAccountGetError) Get(context.Context, string, metav1.GetOptions) (*corev1.ServiceAccount, error) {
	return nil, c.err
}
