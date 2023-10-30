package rbac

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefineClusterRoleBinding(t *testing.T) {
	// Define the RoleRef
	roleRef := DefineRbacAuthorizationClusterRoleRef("my-role")
	// Define the Subjects
	subjects := DefineRbacAuthorizationClusterGroupSubjects([]string{"my-group"})
	// Define the ClusterRoleBinding
	clusterRoleBinding := DefineClusterRoleBinding(*roleRef, *subjects)

	assert.Equal(t, "my-role", clusterRoleBinding.RoleRef.Name)
	assert.Equal(t, "my-group", clusterRoleBinding.Subjects[0].Name)
}

func TestDefineRbacAuthorizationClusterServiceAccountSubjects(t *testing.T) {
	// Define the ClusterRoleBinding
	clusterRoleBinding := DefineRbacAuthorizationClusterServiceAccountSubjects("my-role-binding", "my-namespace", "my-service-account")

	assert.Equal(t, clusterRoleBinding.Subjects[0].Name, "my-service-account")
	assert.Equal(t, clusterRoleBinding.Subjects[0].Namespace, "my-namespace")
}
