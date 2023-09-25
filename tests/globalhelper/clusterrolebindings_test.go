package globalhelper

import (
	"testing"

	"github.com/stretchr/testify/assert"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	k8sfake "k8s.io/client-go/kubernetes/fake"
)

func TestCreateClusterRoleBinding(t *testing.T) {
	testCases := []struct {
		crbAlreadyExists bool
	}{
		{crbAlreadyExists: false},
		{crbAlreadyExists: true},
	}

	for _, tc := range testCases {
		var runtimeObjects []runtime.Object

		if tc.crbAlreadyExists {
			// Create a fake cluster role binding object
			runtimeObjects = append(runtimeObjects, &rbacv1.ClusterRoleBinding{
				ObjectMeta: metav1.ObjectMeta{
					Name: "testCRB",
				},
			})
		}

		// Create a fake clientset
		client := k8sfake.NewSimpleClientset(runtimeObjects...)
		assert.Nil(t, CreateClusterRoleBinding(client.RbacV1(), &rbacv1.ClusterRoleBinding{
			ObjectMeta: metav1.ObjectMeta{
				Name: "testCRB",
			},
		}))
	}
}

func TestDeleteClusterRoleBinding(t *testing.T) {
	testCases := []struct {
		crbAlreadyExists bool
	}{
		{crbAlreadyExists: false},
		{crbAlreadyExists: true},
	}

	for _, tc := range testCases {
		var runtimeObjects []runtime.Object

		if tc.crbAlreadyExists {
			// Create a fake cluster role binding object
			runtimeObjects = append(runtimeObjects, &rbacv1.ClusterRoleBinding{
				ObjectMeta: metav1.ObjectMeta{
					Name: "testCRB",
				},
			})
		}

		// Create a fake clientset
		client := k8sfake.NewSimpleClientset(runtimeObjects...)
		assert.Nil(t, DeleteClusterRoleBinding(client.RbacV1(), "testCRB"))
	}
}
