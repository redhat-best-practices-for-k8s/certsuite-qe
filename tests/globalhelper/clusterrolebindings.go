package globalhelper

import (
	"context"
	"fmt"

	"github.com/golang/glog"
	rbacv1 "k8s.io/api/rbac/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	typedrbacv1 "k8s.io/client-go/kubernetes/typed/rbac/v1"
)

// CreateClusterRoleBinding creates a cluster role binding.
func CreateClusterRoleBinding(clusterRoleBinding *rbacv1.ClusterRoleBinding) error {
	return createClusterRoleBinding(GetAPIClient().K8sClient.RbacV1(), clusterRoleBinding)
}

func createClusterRoleBinding(client typedrbacv1.RbacV1Interface, clusterRoleBinding *rbacv1.ClusterRoleBinding) error {
	_, err := client.ClusterRoleBindings().Create(context.TODO(), clusterRoleBinding, metav1.CreateOptions{})
	if k8serrors.IsAlreadyExists(err) {
		glog.V(5).Info(fmt.Sprintf("cluster role binding %s already exists", clusterRoleBinding.Name))
	} else if err != nil {
		return fmt.Errorf("failed to create cluster role binding: %w", err)
	}

	return nil
}

// DeleteClusterRoleBinding deletes a cluster role binding.
func DeleteClusterRoleBinding(clusterRoleBinding *rbacv1.ClusterRoleBinding) error {
	return deleteClusterRoleBinding(GetAPIClient().K8sClient.RbacV1(), clusterRoleBinding)
}

func deleteClusterRoleBinding(client typedrbacv1.RbacV1Interface, clusterRoleBinding *rbacv1.ClusterRoleBinding) error {
	// Exit if the cluster role binding does not exist
	_, err := client.ClusterRoleBindings().Get(context.TODO(), clusterRoleBinding.Name, metav1.GetOptions{})
	if k8serrors.IsNotFound(err) {
		return nil
	}

	if err := client.ClusterRoleBindings().Delete(context.TODO(), clusterRoleBinding.Name, metav1.DeleteOptions{}); err != nil {
		return fmt.Errorf("failed to delete cluster role binding: %w", err)
	}

	return nil
}

func DeleteClusterRoleBindingByName(name string) error {
	return deleteClusterRoleBindingByName(GetAPIClient().K8sClient.RbacV1(), name)
}

func deleteClusterRoleBindingByName(client typedrbacv1.RbacV1Interface, name string) error {
	// Exit if the cluster role binding does not exist
	_, err := client.ClusterRoleBindings().Get(context.TODO(), name, metav1.GetOptions{})
	if k8serrors.IsNotFound(err) {
		return nil
	}

	if err := client.ClusterRoleBindings().Delete(context.TODO(), name, metav1.DeleteOptions{}); err != nil {
		return fmt.Errorf("failed to delete cluster role binding: %w", err)
	}

	return nil
}
