package globalhelper

import (
	"context"
	"fmt"

	"github.com/golang/glog"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	corev1Typed "k8s.io/client-go/kubernetes/typed/core/v1"
	rbacv1Typed "k8s.io/client-go/kubernetes/typed/rbac/v1"
)

func CreateServiceAccount(client corev1Typed.CoreV1Interface, serviceAccountName, namespace string) error {
	serviceAccount := corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      serviceAccountName,
			Namespace: namespace}}

	_, err :=
		client.ServiceAccounts(namespace).Create(context.TODO(), &serviceAccount, metav1.CreateOptions{})

	if k8serrors.IsAlreadyExists(err) {
		glog.V(5).Info(fmt.Sprintf("serviceaccount %s already installed", serviceAccountName))

		return nil
	}

	return err
}

func DeleteServiceAccount(client corev1Typed.CoreV1Interface, serviceAccountName, namespace string) error {
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

func DefineRole(roleName, namespace string) rbacv1.Role {
	return rbacv1.Role{
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

func RedefineRoleWithAPIGroups(role rbacv1.Role, newAPIGroups []string) {
	role.Rules[0].APIGroups = newAPIGroups
}

func RedefineRoleWithResources(role rbacv1.Role, newResources []string) {
	role.Rules[0].Resources = newResources
}

func CreateRole(client rbacv1Typed.RbacV1Interface, aRole rbacv1.Role) error {
	_, err :=
		client.Roles(aRole.Namespace).Create(context.TODO(), &aRole, metav1.CreateOptions{})

	if k8serrors.IsAlreadyExists(err) {
		glog.V(5).Info(fmt.Sprintf("role %s already installed", aRole.Name))

		return nil
	}

	return err
}

func DeleteRole(client rbacv1Typed.RbacV1Interface, roleName, namespace string) error {
	err :=
		client.Roles(namespace).Delete(context.TODO(), roleName, metav1.DeleteOptions{})

	if k8serrors.IsNotFound(err) {
		glog.V(5).Info(fmt.Sprintf("role %s does not exist", roleName))

		return nil
	}

	return err
}

func DeleteRoleBinding(client rbacv1Typed.RbacV1Interface, bindingName, namespace string) error {
	err :=
		client.RoleBindings(namespace).Delete(context.TODO(), bindingName, metav1.DeleteOptions{})

	if k8serrors.IsNotFound(err) {
		glog.V(5).Info(fmt.Sprintf("rolebinding %s does not exist", bindingName))

		return nil
	}

	return err
}

func CreateRoleBindingWithServiceAccountSubject(client rbacv1Typed.RbacV1Interface, bindingName, roleName, serviceAccountName,
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
		glog.V(5).Info(fmt.Sprintf("rolebinding %s already exists", bindingName))

		return nil
	}

	return err
}
