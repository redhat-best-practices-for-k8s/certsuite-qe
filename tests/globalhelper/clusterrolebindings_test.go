package globalhelper

import (
	"testing"

	egiClients "github.com/openshift-kni/eco-goinfra/pkg/clients"
	"github.com/stretchr/testify/assert"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
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
				Subjects: []rbacv1.Subject{
					{
						Kind:      "User",
						Name:      "testUser",
						Namespace: "testNamespace",
						APIGroup:  "testAPIGroup",
					},
				},
			})
		}

		// Create a fake clientset
		fakeClient := egiClients.GetTestClients(egiClients.TestClientParams{
			K8sMockObjects: runtimeObjects,
		})
		assert.Nil(t, createClusterRoleBinding(fakeClient, &rbacv1.ClusterRoleBinding{
			ObjectMeta: metav1.ObjectMeta{
				Name: "testCRB",
			},
			Subjects: []rbacv1.Subject{
				{
					Kind:      "User",
					Name:      "testUser",
					Namespace: "testNamespace",
					APIGroup:  "testAPIGroup",
				},
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

	for _, testCase := range testCases {
		var runtimeObjects []runtime.Object

		testCRB := &rbacv1.ClusterRoleBinding{
			ObjectMeta: metav1.ObjectMeta{
				Name: "testCRB",
			},
			Subjects: []rbacv1.Subject{
				{
					Kind:      "User",
					Name:      "testUser",
					Namespace: "testNamespace",
					APIGroup:  "testAPIGroup",
				},
			},
		}

		if testCase.crbAlreadyExists {
			// Create a fake cluster role binding object
			runtimeObjects = append(runtimeObjects, testCRB)
		}

		// Create a fake clientset
		fakeClient := egiClients.GetTestClients(egiClients.TestClientParams{
			K8sMockObjects: runtimeObjects,
		})
		assert.Nil(t, deleteClusterRoleBinding(fakeClient, testCRB))
	}
}

func TestDeleteClusterRoleBindingByName(t *testing.T) {
	testCases := []struct {
		crbAlreadyExists bool
	}{
		{crbAlreadyExists: false},
		{crbAlreadyExists: true},
	}

	for _, testCase := range testCases {
		var runtimeObjects []runtime.Object

		testCRB := &rbacv1.ClusterRoleBinding{
			ObjectMeta: metav1.ObjectMeta{
				Name: "testCRB",
			},
			Subjects: []rbacv1.Subject{
				{
					Kind:      "User",
					Name:      "testUser",
					Namespace: "testNamespace",
					APIGroup:  "testAPIGroup",
				},
			},
		}

		if testCase.crbAlreadyExists {
			// Create a fake cluster role binding object
			runtimeObjects = append(runtimeObjects, testCRB)
		}

		// Create a fake clientset
		fakeClient := egiClients.GetTestClients(egiClients.TestClientParams{
			K8sMockObjects: runtimeObjects,
		})
		assert.Nil(t, deleteClusterRoleBindingByName(fakeClient, testCRB.Name))
	}
}
