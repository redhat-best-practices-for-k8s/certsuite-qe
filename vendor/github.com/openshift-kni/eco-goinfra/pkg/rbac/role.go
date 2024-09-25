package rbac

import (
	"context"
	"fmt"

	"github.com/golang/glog"
	"github.com/openshift-kni/eco-goinfra/pkg/clients"
	"github.com/openshift-kni/eco-goinfra/pkg/msg"
	rbacv1 "k8s.io/api/rbac/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// RoleBuilder provides a struct for role object containing connection to the cluster and the role definitions.
type RoleBuilder struct {
	// Role definition. Used to create a role object
	Definition *rbacv1.Role
	// Created role object
	Object *rbacv1.Role

	// Used in functions that define or mutate role definition. errorMsg is processed
	// before the role object is created
	errorMsg  string
	apiClient *clients.Settings
}

// RoleAdditionalOptions additional options for Role object.
type RoleAdditionalOptions func(builder *RoleBuilder) (*RoleBuilder, error)

// NewRoleBuilder create a new instance of RoleBuilder.
func NewRoleBuilder(apiClient *clients.Settings, name, nsname string, rule rbacv1.PolicyRule) *RoleBuilder {
	glog.V(100).Infof(
		"Initializing new role structure with the following params: "+
			"name: %s, namespace: %s, rule %v", name, nsname, rule)

	builder := &RoleBuilder{
		apiClient: apiClient,
		Definition: &rbacv1.Role{
			ObjectMeta: metav1.ObjectMeta{
				Name:      name,
				Namespace: nsname,
			},
		},
	}

	if name == "" {
		glog.V(100).Infof("The name of the role is empty")

		builder.errorMsg = "role 'name' cannot be empty"

		return builder
	}

	if nsname == "" {
		glog.V(100).Infof("The namespace of the role is empty")

		builder.errorMsg = "role 'nsname' cannot be empty"

		return builder
	}

	builder.WithRules([]rbacv1.PolicyRule{rule})

	return builder
}

// WithRules adds the specified PolicyRule to the Role.
func (builder *RoleBuilder) WithRules(rules []rbacv1.PolicyRule) *RoleBuilder {
	if valid, _ := builder.validate(); !valid {
		return builder
	}

	glog.V(100).Infof("Adding to role %s the following rules: %v",
		builder.Definition.Name, rules)

	if len(rules) == 0 {
		glog.V(100).Infof("The list of rules is empty")

		builder.errorMsg = "cannot create role with empty rule"
	}

	if builder.errorMsg != "" {
		return builder
	}

	for _, rule := range rules {
		if len(rule.Verbs) == 0 {
			glog.V(100).Infof("The role has no verbs")

			builder.errorMsg = "role must contain at least one Verb"
		}

		if len(rule.Resources) == 0 {
			glog.V(100).Infof("The role has no resources")

			builder.errorMsg = "role must contain at least one Resource"
		}

		if len(rule.APIGroups) == 0 {
			glog.V(100).Infof("The role has no apigroups")

			builder.errorMsg = "role must contain at least one APIGroup"
		}

		if builder.errorMsg != "" {
			return builder
		}
	}

	if builder.Definition.Rules == nil {
		builder.Definition.Rules = rules
	} else {
		builder.Definition.Rules = append(builder.Definition.Rules, rules...)
	}

	return builder
}

// WithOptions creates Role with generic mutation options.
func (builder *RoleBuilder) WithOptions(
	options ...RoleAdditionalOptions) *RoleBuilder {
	if valid, _ := builder.validate(); !valid {
		return builder
	}

	glog.V(100).Infof("Setting Role additional options")

	for _, option := range options {
		if option != nil {
			builder, err := option(builder)

			if err != nil {
				glog.V(100).Infof("Error occurred in mutation function")

				builder.errorMsg = err.Error()

				return builder
			}
		}
	}

	return builder
}

// PullRole pulls existing role from cluster.
func PullRole(apiClient *clients.Settings, name, nsname string) (*RoleBuilder, error) {
	glog.V(100).Infof("Pulling existing role name %s under namespace %s from cluster", name, nsname)

	if apiClient == nil {
		glog.V(100).Infof("The apiClient cannot be nil")

		return nil, fmt.Errorf("the apiClient cannot be nil")
	}

	builder := RoleBuilder{
		apiClient: apiClient,
		Definition: &rbacv1.Role{
			ObjectMeta: metav1.ObjectMeta{
				Name:      name,
				Namespace: nsname,
			},
		},
	}

	if name == "" {
		glog.V(100).Infof("The name of the role is empty")

		return nil, fmt.Errorf("role 'name' cannot be empty")
	}

	if nsname == "" {
		glog.V(100).Infof("The namespace of the role is empty")

		return nil, fmt.Errorf("role 'namespace' cannot be empty")
	}

	if !builder.Exists() {
		return nil, fmt.Errorf("role object %s does not exist in namespace %s", name, nsname)
	}

	builder.Definition = builder.Object

	return &builder, nil
}

// Create makes a Role in the cluster and stores the created object in struct.
func (builder *RoleBuilder) Create() (*RoleBuilder, error) {
	if valid, err := builder.validate(); !valid {
		return builder, err
	}

	glog.V(100).Infof("Creating role %s under namespace %s",
		builder.Definition.Name, builder.Definition.Namespace)

	var err error
	if !builder.Exists() {
		builder.Object, err = builder.apiClient.Roles(builder.Definition.Namespace).Create(
			context.TODO(), builder.Definition, metav1.CreateOptions{})
	}

	return builder, err
}

// Delete removes a Role.
func (builder *RoleBuilder) Delete() error {
	if valid, err := builder.validate(); !valid {
		return err
	}

	glog.V(100).Infof("Removing role %s under namespace %s",
		builder.Definition.Name, builder.Definition.Namespace)

	if !builder.Exists() {
		return nil
	}

	err := builder.apiClient.Roles(builder.Definition.Namespace).Delete(
		context.TODO(), builder.Definition.Name, metav1.DeleteOptions{})

	if err != nil {
		return err
	}

	builder.Object = nil

	return err
}

// Update modifies the existing Role object with role definition in builder.
func (builder *RoleBuilder) Update() (*RoleBuilder, error) {
	if valid, err := builder.validate(); !valid {
		return builder, err
	}

	glog.V(100).Infof("Updating role %s under namespace %s",
		builder.Definition.Name, builder.Definition.Namespace)

	if !builder.Exists() {
		return nil, fmt.Errorf("role object %s does not exist, fail to update", builder.Definition.Name)
	}

	var err error
	builder.Object, err = builder.apiClient.Roles(builder.Definition.Namespace).Update(
		context.TODO(), builder.Definition, metav1.UpdateOptions{})

	return builder, err
}

// Exists checks whether the given Role exists.
func (builder *RoleBuilder) Exists() bool {
	if valid, _ := builder.validate(); !valid {
		return false
	}

	glog.V(100).Infof("Checking if role %s exists under namespace %s",
		builder.Definition.Name, builder.Definition.Namespace)

	var err error
	builder.Object, err = builder.apiClient.Roles(builder.Definition.Namespace).Get(
		context.TODO(), builder.Definition.Name, metav1.GetOptions{})

	return err == nil || !k8serrors.IsNotFound(err)
}

// validate will check that the builder and builder definition are properly initialized before
// accessing any member fields.
func (builder *RoleBuilder) validate() (bool, error) {
	resourceCRD := "role"

	if builder == nil {
		glog.V(100).Infof("The %s builder is uninitialized", resourceCRD)

		return false, fmt.Errorf("error: received nil %s builder", resourceCRD)
	}

	if builder.Definition == nil {
		glog.V(100).Infof("The %s is undefined", resourceCRD)

		builder.errorMsg = msg.UndefinedCrdObjectErrString(resourceCRD)
	}

	if builder.apiClient == nil {
		glog.V(100).Infof("The %s builder apiclient is nil", resourceCRD)

		builder.errorMsg = fmt.Sprintf("%s builder cannot have nil apiClient", resourceCRD)
	}

	if builder.errorMsg != "" {
		glog.V(100).Infof("The %s builder has error message: %s", resourceCRD, builder.errorMsg)

		return false, fmt.Errorf(builder.errorMsg)
	}

	return true, nil
}
