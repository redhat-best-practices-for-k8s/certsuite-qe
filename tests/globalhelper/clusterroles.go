package globalhelper

import (
	"context"
	"fmt"

	"github.com/golang/glog"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	typedrbacv1 "k8s.io/client-go/kubernetes/typed/rbac/v1"
)

func DeleteClusterRole(name string) error {
	return deleteClusterRole(GetAPIClient().K8sClient.RbacV1(), name)
}

func deleteClusterRole(client typedrbacv1.RbacV1Interface, name string) error {
	err := client.ClusterRoles().Delete(
		context.TODO(),
		name,
		metav1.DeleteOptions{},
	)
	if k8serrors.IsNotFound(err) {
		glog.V(5).Info(fmt.Sprintf("cluster-role %s is not found", name))

		return nil
	} else if err != nil {
		return fmt.Errorf("failed to delete cluster role %q: %w", name, err)
	}

	return nil
}
