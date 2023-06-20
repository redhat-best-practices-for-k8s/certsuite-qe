package globalhelper

import (
	"context"
	"fmt"

	"github.com/golang/glog"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func CreateServiceAccount(serviceAccountName, namespace string) error {
	serviceAccount := corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      serviceAccountName,
			Namespace: namespace}}

	_, err :=
		APIClient.CoreV1Interface.ServiceAccounts(namespace).Create(context.Background(), &serviceAccount, metav1.CreateOptions{})

	if k8serrors.IsAlreadyExists(err) {
		glog.V(5).Info(fmt.Sprintf("serviceaccount %s already installed", serviceAccountName))

		return nil
	}

	return err
}

func DeleteServiceAccount(serviceAccountName, namespace string) error {
	serviceAccount, err := APIClient.ServiceAccounts(namespace).Get(context.Background(), serviceAccountName, metav1.GetOptions{})

	if k8serrors.IsNotFound(err) {
		return nil
	}

	if serviceAccount == nil {
		return err
	}

	err = APIClient.ServiceAccounts(namespace).Delete(context.Background(), serviceAccountName,
		metav1.DeleteOptions{})

	if err != nil {
		return fmt.Errorf("failed to delete serviceaccount %w", err)
	}

	return nil
}

func CreateRole(roleName, namespace string) error {
	aRole := rbacv1.Role{
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

	_, err :=
		APIClient.RbacV1Interface.Roles(namespace).Create(context.Background(), &aRole, metav1.CreateOptions{})

	if k8serrors.IsAlreadyExists(err) {
		glog.V(5).Info(fmt.Sprintf("role %s already installed", roleName))

		return nil
	}

	return err
}

func DeleteRoleBinding(bindingName, namespace string) error {
	err :=
		APIClient.RbacV1Interface.RoleBindings(namespace).Delete(context.Background(), bindingName, metav1.DeleteOptions{})

	if k8serrors.IsNotFound(err) {
		glog.V(5).Info(fmt.Sprintf("rolebinding %s does not exist", bindingName))

		return nil
	}

	return err
}

func CreateRoleBindingWithServiceAccountSubject(bindingName, roleName, serviceAccountName, namespace string) error {
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
			Namespace: namespace,
		}},
		RoleRef: rbacv1.RoleRef{
			Kind:     "Role",
			Name:     roleName,
			APIGroup: "rbac.authorization.k8s.io",
		},
	}

	_, err :=
		APIClient.RbacV1Interface.RoleBindings(namespace).Create(context.Background(), &aRoleBinding, metav1.CreateOptions{})

	if k8serrors.IsAlreadyExists(err) {
		glog.V(5).Info(fmt.Sprintf("rolebinding %s already exists", bindingName))

		return nil
	}

	return err
}
