package globalhelper

import (
	"testing"

	egiClients "github.com/openshift-kni/eco-goinfra/pkg/clients"
	"github.com/stretchr/testify/assert"
	rbacv1 "k8s.io/api/rbac/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func TestDeleteClusterRole(t *testing.T) {
	generateClusterRole := func(name string) *rbacv1.ClusterRole {
		return &rbacv1.ClusterRole{
			ObjectMeta: metav1.ObjectMeta{
				Name: name,
			},
		}
	}

	testCases := []struct {
		alreadyExists bool
	}{
		{alreadyExists: false},
		{alreadyExists: true},
	}

	for _, testCase := range testCases {
		var runtimeObjects []runtime.Object

		testCR := generateClusterRole("testCR")

		if testCase.alreadyExists {
			runtimeObjects = append(runtimeObjects, testCR)
		}

		fakeClient := egiClients.GetTestClients(egiClients.TestClientParams{
			K8sMockObjects: runtimeObjects,
		})
		assert.Nil(t, deleteClusterRole(fakeClient, testCR.Name))

		_, err := fakeClient.RbacV1Interface.ClusterRoles().Get(t.Context(), testCR.Name, metav1.GetOptions{})
		assert.NotNil(t, err)
		assert.True(t, k8serrors.IsNotFound(err))
	}
}
