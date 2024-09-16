package globalhelper

import (
	egiClients "github.com/openshift-kni/eco-goinfra/pkg/clients"
	egiRbac "github.com/openshift-kni/eco-goinfra/pkg/rbac"
	rbacv1 "k8s.io/api/rbac/v1"
)

func DeleteClusterRole(name string) error {
	return deleteClusterRole(egiClients.New(""), name)
}

func deleteClusterRole(client *egiClients.Settings, name string) error {
	return egiRbac.NewClusterRoleBuilder(client, name, rbacv1.PolicyRule{
		Verbs:     []string{"get"},
		APIGroups: []string{""},
	}).Delete()
}
