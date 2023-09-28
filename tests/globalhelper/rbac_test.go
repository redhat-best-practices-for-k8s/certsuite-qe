package globalhelper

import (
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	k8sfake "k8s.io/client-go/kubernetes/fake"
)

func TestCreateServiceAccount(t *testing.T) {
	testCases := []struct {
		saAlreadyExists bool
	}{
		{saAlreadyExists: false},
		{saAlreadyExists: true},
	}

	for _, tc := range testCases {
		var runtimeObjects []runtime.Object

		if tc.saAlreadyExists {
			// Create a fake service account object
			runtimeObjects = append(runtimeObjects, &corev1.ServiceAccount{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "testSA",
					Namespace: "default",
				},
			})
		}

		// Create a fake clientset
		client := k8sfake.NewSimpleClientset(runtimeObjects...)
		assert.Nil(t, CreateServiceAccount(client.CoreV1(), "testSA", "default"))
	}
}

func TestDeleteServiceAccount(t *testing.T) {
	testCases := []struct {
		saAlreadyExists bool
	}{
		{saAlreadyExists: false},
		{saAlreadyExists: true},
	}

	for _, tc := range testCases {
		var runtimeObjects []runtime.Object

		if tc.saAlreadyExists {
			// Create a fake service account object
			runtimeObjects = append(runtimeObjects, &corev1.ServiceAccount{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "testSA",
					Namespace: "default",
				},
			})
		}

		// Create a fake clientset
		client := k8sfake.NewSimpleClientset(runtimeObjects...)
		assert.Nil(t, DeleteServiceAccount(client.CoreV1(), "testSA", "default"))
	}
}

func TestDefineRole(t *testing.T) {
	testCases := []struct {
		roleName  string
		namespace string
	}{
		{roleName: "testRole", namespace: "default"},
	}

	for _, tc := range testCases {
		assert.Equal(t, tc.roleName, DefineRole(tc.roleName, tc.namespace).Name)
		assert.Equal(t, tc.namespace, DefineRole(tc.roleName, tc.namespace).Namespace)
	}
}

func TestRedefineRoleWithAPIGroups(t *testing.T) {
	testCases := []struct {
		roleName  string
		namespace string
	}{
		{roleName: "testRole", namespace: "default"},
	}

	for _, tc := range testCases {
		role := DefineRole(tc.roleName, tc.namespace)
		role.Rules[0].APIGroups = []string{"testAPIGroup"}
		assert.Equal(t, []string{"testAPIGroup"}, role.Rules[0].APIGroups)
	}
}

func TestRedefineRoleWithResources(t *testing.T) {
	testCases := []struct {
		roleName  string
		namespace string
	}{
		{roleName: "testRole", namespace: "default"},
	}

	for _, tc := range testCases {
		role := DefineRole(tc.roleName, tc.namespace)
		role.Rules[0].Resources = []string{"testResource"}
		assert.Equal(t, []string{"testResource"}, role.Rules[0].Resources)
	}
}

func TestCreateRole(t *testing.T) {
	testCases := []struct {
		roleAlreadyExists bool
	}{
		{roleAlreadyExists: false},
		{roleAlreadyExists: true},
	}

	for _, tc := range testCases {
		var runtimeObjects []runtime.Object

		if tc.roleAlreadyExists {
			// Create a fake role object
			runtimeObjects = append(runtimeObjects, &corev1.ServiceAccount{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "testRole",
					Namespace: "default",
				},
			})
		}

		// Create a fake clientset
		client := k8sfake.NewSimpleClientset(runtimeObjects...)
		assert.Nil(t, CreateRole(client.RbacV1(), DefineRole("testRole", "default")))
	}
}

func TestDeleteRole(t *testing.T) {
	testCases := []struct {
		roleAlreadyExists bool
	}{
		{roleAlreadyExists: false},
		{roleAlreadyExists: true},
	}

	for _, tc := range testCases {
		var runtimeObjects []runtime.Object

		if tc.roleAlreadyExists {
			// Create a fake role object
			runtimeObjects = append(runtimeObjects, &corev1.ServiceAccount{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "testRole",
					Namespace: "default",
				},
			})
		}

		// Create a fake clientset
		client := k8sfake.NewSimpleClientset(runtimeObjects...)
		assert.Nil(t, DeleteRole(client.RbacV1(), "testRole", "default"))
	}
}

func TestDeleteRoleBinding(t *testing.T) {
	testCases := []struct {
		roleBindingAlreadyExists bool
	}{
		{roleBindingAlreadyExists: false},
		{roleBindingAlreadyExists: true},
	}

	for _, tc := range testCases {
		var runtimeObjects []runtime.Object

		if tc.roleBindingAlreadyExists {
			// Create a fake role binding object
			runtimeObjects = append(runtimeObjects, &corev1.ServiceAccount{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "testRoleBinding",
					Namespace: "default",
				},
			})
		}

		// Create a fake clientset
		client := k8sfake.NewSimpleClientset(runtimeObjects...)
		assert.Nil(t, DeleteRoleBinding(client.RbacV1(), "testRoleBinding", "default"))
	}
}

func TestCreateRoleBindingWithServiceAccountSubject(t *testing.T) {
	testCases := []struct {
		roleBindingAlreadyExists bool
	}{
		{roleBindingAlreadyExists: false},
		{roleBindingAlreadyExists: true},
	}

	for _, tc := range testCases {
		var runtimeObjects []runtime.Object

		if tc.roleBindingAlreadyExists {
			// Create a fake role binding object
			runtimeObjects = append(runtimeObjects, &corev1.ServiceAccount{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "testRoleBinding",
					Namespace: "default",
				},
			})
		}

		// Create a fake clientset
		client := k8sfake.NewSimpleClientset(runtimeObjects...)
		assert.Nil(t, CreateRoleBindingWithServiceAccountSubject(client.RbacV1(), "testRoleBinding",
			"my-role", "testSA", "default", "default"))
	}
}
