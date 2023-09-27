package rbac

import (
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// DefineClusterRoleBinding sets cluster ClusterRoleBinding for role and subject.
func DefineClusterRoleBinding(ref rbacv1.RoleRef, subjects []rbacv1.Subject) *rbacv1.ClusterRoleBinding {
	return &rbacv1.ClusterRoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name: "system:openshift:scc:privileged",
		},
		RoleRef:  ref,
		Subjects: subjects,
	}
}

// DefineRbacAuthorizationClusterRoleRef defines RoleRef struct.
func DefineRbacAuthorizationClusterRoleRef(roleRefName string) *rbacv1.RoleRef {
	return &rbacv1.RoleRef{
		Name:     roleRefName,
		Kind:     "ClusterRole",
		APIGroup: "rbac.authorization.k8s.io",
	}
}

// DefineRbacAuthorizationClusterGroupSubjects defines RBAC Subject list.
func DefineRbacAuthorizationClusterGroupSubjects(subjectNames []string) *[]rbacv1.Subject {
	var Subjects []rbacv1.Subject
	for _, subjectName := range subjectNames {
		Subjects = append(Subjects, rbacv1.Subject{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "Group",
			Name:     subjectName,
		})
	}

	return &Subjects
}

func DefineRbacAuthorizationClusterServiceAccountSubjects(clusterRoleBindingName, serviceAccountNamespace,
	serviceAccountName string) *rbacv1.ClusterRoleBinding {
	// Define the RoleBinding
	return &rbacv1.ClusterRoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name: clusterRoleBindingName,
		},
		Subjects: []rbacv1.Subject{
			{
				Kind:      "ServiceAccount",
				Name:      serviceAccountName,
				Namespace: serviceAccountNamespace,
			},
		},
		RoleRef: rbacv1.RoleRef{
			Kind:     "ClusterRole",
			Name:     "my-role",
			APIGroup: "rbac.authorization.k8s.io",
		},
	}
}
