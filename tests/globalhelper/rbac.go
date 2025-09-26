package globalhelper

import (
	"context"
	"fmt"
	"time"

	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	corev1Typed "k8s.io/client-go/kubernetes/typed/core/v1"
	rbacv1Typed "k8s.io/client-go/kubernetes/typed/rbac/v1"
	klog "k8s.io/klog/v2"
	"k8s.io/utils/ptr"
)

// CreateServiceAccount creates a service account.
func CreateServiceAccount(serviceAccountName, namespace string) error {
	return createServiceAccount(GetAPIClient().K8sClient.CoreV1(), serviceAccountName, namespace)
}

func createServiceAccount(client corev1Typed.CoreV1Interface, serviceAccountName, namespace string) error {
	serviceAccount := corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      serviceAccountName,
			Namespace: namespace}}

	_, err :=
		client.ServiceAccounts(namespace).Create(context.TODO(), &serviceAccount, metav1.CreateOptions{})

	if k8serrors.IsAlreadyExists(err) {
		klog.V(5).Info(fmt.Sprintf("serviceaccount %s already installed", serviceAccountName))

		return nil
	}

	return err
}

// DeleteServiceAccount deletes a service account.
func DeleteServiceAccount(serviceAccountName, namespace string) error {
	return deleteServiceAccount(GetAPIClient().K8sClient.CoreV1(), serviceAccountName, namespace)
}

func deleteServiceAccount(client corev1Typed.CoreV1Interface, serviceAccountName, namespace string) error {
	serviceAccount, err := client.ServiceAccounts(namespace).Get(context.TODO(), serviceAccountName, metav1.GetOptions{})

	if k8serrors.IsNotFound(err) {
		return nil
	}

	if serviceAccount == nil {
		return err
	}

	err = client.ServiceAccounts(namespace).Delete(context.TODO(), serviceAccountName,
		metav1.DeleteOptions{})

	if err != nil {
		return fmt.Errorf("failed to delete serviceaccount %w", err)
	}

	return nil
}

func DefineRole(roleName, namespace string) *rbacv1.Role {
	return &rbacv1.Role{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Role",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      roleName,
			Namespace: namespace,
		},
		Rules: []rbacv1.PolicyRule{{
			APIGroups: []string{"*"},
			Resources: []string{"*"},
			Verbs:     []string{"*"},
		},
		},
	}
}

func RedefineRoleWithAPIGroups(role *rbacv1.Role, newAPIGroups []string) {
	role.Rules[0].APIGroups = newAPIGroups
}

func RedefineRoleWithResources(role *rbacv1.Role, newResources []string) {
	role.Rules[0].Resources = newResources
}

// CreateRole creates a role.
func CreateRole(aRole *rbacv1.Role) error {
	return createRole(GetAPIClient().K8sClient.RbacV1(), aRole)
}

func createRole(client rbacv1Typed.RbacV1Interface, aRole *rbacv1.Role) error {
	_, err :=
		client.Roles(aRole.Namespace).Create(context.TODO(), aRole, metav1.CreateOptions{})

	if k8serrors.IsAlreadyExists(err) {
		klog.V(5).Info(fmt.Sprintf("role %s already installed", aRole.Name))

		return nil
	}

	return err
}

// DeleteRole deletes a role.
func DeleteRole(roleName, namespace string) error {
	return deleteRole(GetAPIClient().K8sClient.RbacV1(), roleName, namespace)
}

func deleteRole(client rbacv1Typed.RbacV1Interface, roleName, namespace string) error {
	err :=
		client.Roles(namespace).Delete(context.TODO(), roleName, metav1.DeleteOptions{})

	if k8serrors.IsNotFound(err) {
		klog.V(5).Info(fmt.Sprintf("role %s does not exist", roleName))

		return nil
	}

	return err
}

// DeleteRoleBinding deletes a role binding.
func DeleteRoleBinding(roleBindingName, namespace string) error {
	return deleteRoleBinding(GetAPIClient().K8sClient.RbacV1(), roleBindingName, namespace)
}

func deleteRoleBinding(client rbacv1Typed.RbacV1Interface, bindingName, namespace string) error {
	err :=
		client.RoleBindings(namespace).Delete(context.TODO(), bindingName, metav1.DeleteOptions{})

	if k8serrors.IsNotFound(err) {
		klog.V(5).Info(fmt.Sprintf("rolebinding %s does not exist", bindingName))

		return nil
	}

	return err
}

// CreateRoleBindingWithServiceAccountSubject creates a role binding with a service account subject.
func CreateRoleBindingWithServiceAccountSubject(bindingName, roleName, serviceAccountName,
	serviceAccountNamespace, namespace string) error {
	return createRoleBindingWithServiceAccountSubject(GetAPIClient().K8sClient.RbacV1(), bindingName, roleName,
		serviceAccountName, serviceAccountNamespace, namespace)
}

func createRoleBindingWithServiceAccountSubject(client rbacv1Typed.RbacV1Interface, bindingName, roleName, serviceAccountName,
	serviceAccountNamespace, namespace string) error {
	aRoleBinding := rbacv1.RoleBinding{
		TypeMeta: metav1.TypeMeta{
			Kind:       "RoleBinding",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      bindingName,
			Namespace: namespace,
		},
		Subjects: []rbacv1.Subject{{
			Kind:      "ServiceAccount",
			Name:      serviceAccountName,
			Namespace: serviceAccountNamespace,
		}},
		RoleRef: rbacv1.RoleRef{
			Kind:     "Role",
			Name:     roleName,
			APIGroup: "rbac.authorization.k8s.io",
		},
	}

	_, err :=
		client.RoleBindings(namespace).Create(context.TODO(), &aRoleBinding, metav1.CreateOptions{})

	if k8serrors.IsAlreadyExists(err) {
		klog.V(5).Info(fmt.Sprintf("rolebinding %s already exists", bindingName))

		return nil
	}

	return err
}

// CreatePersistentVolume creates a persistent volume.
func CreatePersistentVolume(persistentVolume *corev1.PersistentVolume) error {
	return createPersistentVolume(GetAPIClient().K8sClient.CoreV1(), persistentVolume)
}

func createPersistentVolume(client corev1Typed.CoreV1Interface, persistentVolume *corev1.PersistentVolume) error {
	_, err := client.PersistentVolumes().Create(context.TODO(), persistentVolume, metav1.CreateOptions{})
	if k8serrors.IsAlreadyExists(err) {
		klog.V(5).Info(fmt.Sprintf("persistent volume %s already created", persistentVolume.Name))
	} else if err != nil {
		return fmt.Errorf("failed to create persistent volume: %w", err)
	}

	return nil
}

// DeletePersistentVolume deletes a persistent volume.
func DeletePersistentVolume(persistentVolume string, timeout time.Duration) error {
	return deletePersistentVolume(GetAPIClient().K8sClient.CoreV1(), persistentVolume, timeout)
}

func deletePersistentVolume(client corev1Typed.CoreV1Interface, persistentVolume string, timeout time.Duration) error {
	err := client.PersistentVolumes().Delete(context.TODO(), persistentVolume, metav1.DeleteOptions{
		GracePeriodSeconds: ptr.To[int64](0),
	})
	if k8serrors.IsNotFound(err) {
		klog.V(5).Info(fmt.Sprintf("persistent volume %s does not exist", persistentVolume))

		return nil
	} else if err != nil {
		return fmt.Errorf("failed to delete persistent volume %w", err)
	}

	Eventually(func() bool {
		// if the pv was deleted, we will get an error.
		_, err := client.PersistentVolumes().Get(context.TODO(), persistentVolume, metav1.GetOptions{})

		return err != nil
	}, timeout, retryInterval*time.Second).Should(Equal(true), "PV is not removed yet.")

	return nil
}

func CreateAndWaitUntilPVCIsBound(pvc *corev1.PersistentVolumeClaim, timeout time.Duration, pvName string) error {
	return createAndWaitUntilPVCIsBound(GetAPIClient().K8sClient.CoreV1(), pvc, timeout, pvName)
}

func createAndWaitUntilPVCIsBound(client corev1Typed.CoreV1Interface, pvc *corev1.PersistentVolumeClaim,
	timeout time.Duration, pvName string) error {
	pvc, err := client.PersistentVolumeClaims(pvc.Namespace).Create(context.TODO(), pvc, metav1.CreateOptions{})
	if k8serrors.IsAlreadyExists(err) {
		klog.V(5).Info(fmt.Sprintf("persistent volume claim %s already created", pvc.Name))
	} else if err != nil {
		return fmt.Errorf("failed to create persistent volume claim: %w", err)
	}

	Eventually(func() bool {
		status, err := isPvcBound(client, pvc.Name, pvc.Namespace, pvName)
		if err != nil {
			klog.V(5).Info(fmt.Sprintf(
				"pvc %s is not bound, retry in %d seconds", pvc.Name, retryInterval))

			return false
		}

		return status
	}, timeout, retryInterval*time.Second).Should(Equal(true), "pvc is not bound")

	return nil
}

func isPvcBound(client corev1Typed.CoreV1Interface, pvcName string, namespace string, pvName string) (bool, error) {
	pvc, err := client.PersistentVolumeClaims(namespace).Get(context.TODO(), pvcName, metav1.GetOptions{})
	if err != nil {
		return false, err
	}

	return pvc.Status.Phase == corev1.ClaimBound && pvc.Spec.VolumeName == pvName, nil
}

func DeletePersistentVolumeClaim(pvc *corev1.PersistentVolumeClaim) error {
	return deletePersistentVolumeClaim(GetAPIClient().K8sClient.CoreV1(), pvc)
}

func deletePersistentVolumeClaim(client corev1Typed.CoreV1Interface, pvc *corev1.PersistentVolumeClaim) error {
	err := client.PersistentVolumeClaims(pvc.Namespace).Delete(context.TODO(), pvc.Name, metav1.DeleteOptions{
		GracePeriodSeconds: ptr.To[int64](0),
	})
	if k8serrors.IsNotFound(err) {
		klog.V(5).Info(fmt.Sprintf("persistent volume claim %s does not exist", pvc.Name))

		return nil
	} else if err != nil {
		return fmt.Errorf("failed to delete persistent volume claim %w", err)
	}

	return nil
}
