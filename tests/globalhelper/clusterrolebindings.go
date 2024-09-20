package globalhelper

import (
	egiClients "github.com/openshift-kni/eco-goinfra/pkg/clients"
	egiRbac "github.com/openshift-kni/eco-goinfra/pkg/rbac"
	rbacv1 "k8s.io/api/rbac/v1"
)

// CreateClusterRoleBinding creates a cluster role binding.
func CreateClusterRoleBinding(clusterRoleBinding *rbacv1.ClusterRoleBinding) error {
	return createClusterRoleBinding(egiClients.New(""), clusterRoleBinding)
}

func createClusterRoleBinding(client *egiClients.Settings, clusterRoleBinding *rbacv1.ClusterRoleBinding) error {
	_, err := egiRbac.NewClusterRoleBindingBuilder(client,
		clusterRoleBinding.Name, clusterRoleBinding.RoleRef.Name, rbacv1.Subject{
			Kind:      clusterRoleBinding.Subjects[0].Kind,
			Name:      clusterRoleBinding.Subjects[0].Name,
			Namespace: clusterRoleBinding.Subjects[0].Namespace,
			APIGroup:  clusterRoleBinding.Subjects[0].APIGroup,
		}).Create()

	return err
}

// DeleteClusterRoleBinding deletes a cluster role binding.
func DeleteClusterRoleBinding(clusterRoleBinding *rbacv1.ClusterRoleBinding) error {
	return deleteClusterRoleBinding(egiClients.New(""), clusterRoleBinding)
}

func deleteClusterRoleBinding(client *egiClients.Settings, clusterRoleBinding *rbacv1.ClusterRoleBinding) error {
	return egiRbac.NewClusterRoleBindingBuilder(client,
		clusterRoleBinding.Name, clusterRoleBinding.RoleRef.Name, rbacv1.Subject{
			Kind:      clusterRoleBinding.Subjects[0].Kind,
			Name:      clusterRoleBinding.Subjects[0].Name,
			Namespace: clusterRoleBinding.Subjects[0].Namespace,
			APIGroup:  clusterRoleBinding.Subjects[0].APIGroup,
		}).Delete()
}

func DeleteClusterRoleBindingByName(name string) error {
	return deleteClusterRoleBindingByName(egiClients.New(""), name)
}

func deleteClusterRoleBindingByName(client *egiClients.Settings, name string) error {
	return egiRbac.NewClusterRoleBindingBuilder(client, name, "notempty", rbacv1.Subject{
		Kind:      "User",
		Name:      "notempty",
		Namespace: "notempty",
		APIGroup:  "notempty",
	}).Delete()
}
