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

func CreateClusterRoleBinding(client typedrbacv1.RbacV1Interface, clusterRoleBinding *rbacv1.ClusterRoleBinding) error {
	_, err := client.ClusterRoleBindings().Create(context.TODO(), clusterRoleBinding, metav1.CreateOptions{})
	if k8serrors.IsAlreadyExists(err) {
		glog.V(5).Info(fmt.Sprintf("cluster role binding %s already exists", clusterRoleBinding.Name))
	} else if err != nil {
		return fmt.Errorf("failed to create cluster role binding: %w", err)
	}

	return nil
}

func DeleteClusterRoleBinding(client typedrbacv1.RbacV1Interface, name string) error {
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
