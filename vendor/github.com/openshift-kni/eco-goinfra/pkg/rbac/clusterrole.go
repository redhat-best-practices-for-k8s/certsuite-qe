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

/*
ClusterRoleBuilder provides struct for clusterrole object

	containing connection to the cluster and the clusterrole definitions.
*/
type ClusterRoleBuilder struct {
	// Clusterrole definition. Used to create a clusterrole object.
	Definition *rbacv1.ClusterRole
	// Created clusterrole object
	Object *rbacv1.ClusterRole
	// Used in functions that define or mutate clusterrole definition. errorMsg is processed before clusterrole
	// object is created.
	errorMsg  string
	apiClient *clients.Settings
}

// ClusterRoleAdditionalOptions additional options for ClusterRole object.
type ClusterRoleAdditionalOptions func(builder *ClusterRoleBuilder) (*ClusterRoleBuilder, error)

// NewClusterRoleBuilder creates new instance of ClusterRoleBuilder.
func NewClusterRoleBuilder(apiClient *clients.Settings, name string, rule rbacv1.PolicyRule) *ClusterRoleBuilder {
	glog.V(100).Infof(
		"Initializing new clusterrole structure with the following params: "+
			"name: %s, policy rule: %v",
		name, rule)

	builder := ClusterRoleBuilder{
		apiClient: apiClient,
		Definition: &rbacv1.ClusterRole{
			ObjectMeta: metav1.ObjectMeta{
				Name: name,
			},
		},
	}

	if name == "" {
		glog.V(100).Infof("The name of the clusterrole is empty")

		builder.errorMsg = "clusterrole 'name' cannot be empty"
	}

	builder.WithRules([]rbacv1.PolicyRule{rule})

	return &builder
}

// WithRules appends additional rules to the clusterrole definition.
func (builder *ClusterRoleBuilder) WithRules(rules []rbacv1.PolicyRule) *ClusterRoleBuilder {
	if valid, _ := builder.validate(); !valid {
		return builder
	}

	glog.V(100).Infof("Appending to the definition of clusterrole %s these additional rules %v",
		builder.Definition.Name, rules)

	if len(rules) == 0 {
		glog.V(100).Infof("The list of rules is empty")

		builder.errorMsg = "cannot accept nil or empty slice as rules"
	}

	if builder.errorMsg != "" {
		return builder
	}

	for _, rule := range rules {
		if len(rule.Verbs) == 0 {
			glog.V(100).Infof("The clusterrole rule must contain at least one Verb entry")

			builder.errorMsg = "clusterrole rule must contain at least one Verb entry"
		}

		if builder.errorMsg != "" {
			return builder
		}
	}

	if builder.Definition.Rules == nil {
		builder.Definition.Rules = rules

		return builder
	}

	builder.Definition.Rules = append(builder.Definition.Rules, rules...)

	return builder
}

// WithOptions creates ClusterRole with generic mutation options.
func (builder *ClusterRoleBuilder) WithOptions(options ...ClusterRoleAdditionalOptions) *ClusterRoleBuilder {
	if valid, _ := builder.validate(); !valid {
		return builder
	}

	glog.V(100).Infof("Setting ClusterRole additional options")

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

// PullClusterRole pulls existing clusterrole from cluster.
func PullClusterRole(apiClient *clients.Settings, name string) (*ClusterRoleBuilder, error) {
	glog.V(100).Infof("Pulling existing clusterrole name %s from cluster", name)

	if apiClient == nil {
		glog.V(100).Infof("The apiClient cannot be nil")

		return nil, fmt.Errorf("the apiClient cannot be nil")
	}

	builder := ClusterRoleBuilder{
		apiClient: apiClient,
		Definition: &rbacv1.ClusterRole{
			ObjectMeta: metav1.ObjectMeta{
				Name: name,
			},
		},
	}

	if name == "" {
		glog.V(100).Infof("The name of the clusterrole is empty")

		return nil, fmt.Errorf("clusterrole 'name' cannot be empty")
	}

	if !builder.Exists() {
		return nil, fmt.Errorf("clusterrole object %s does not exist", name)
	}

	builder.Definition = builder.Object

	return &builder, nil
}

// Create generates a clusterrole in the cluster and stores the created object in struct.
func (builder *ClusterRoleBuilder) Create() (*ClusterRoleBuilder, error) {
	if valid, err := builder.validate(); !valid {
		return builder, err
	}

	glog.V(100).Infof("Creating clusterrole %s",
		builder.Definition.Name)

	var err error
	if !builder.Exists() {
		builder.Object, err = builder.apiClient.ClusterRoles().Create(
			context.TODO(), builder.Definition, metav1.CreateOptions{})
	}

	return builder, err
}

// Delete removes a clusterrole from the cluster.
func (builder *ClusterRoleBuilder) Delete() error {
	if valid, err := builder.validate(); !valid {
		return err
	}

	glog.V(100).Infof("Removing clusterrole %s",
		builder.Definition.Name)

	if !builder.Exists() {
		return nil
	}

	err := builder.apiClient.ClusterRoles().Delete(
		context.TODO(), builder.Definition.Name, metav1.DeleteOptions{})

	if err != nil {
		return err
	}

	builder.Object = nil

	return err
}

// Update modifies a clusterrole object in the cluster.
func (builder *ClusterRoleBuilder) Update() (*ClusterRoleBuilder, error) {
	if valid, err := builder.validate(); !valid {
		return builder, err
	}

	glog.V(100).Infof("Updating clusterrole %s",
		builder.Definition.Name)

	if !builder.Exists() {
		return nil, fmt.Errorf("clusterrole object %s does not exist, fail to update", builder.Definition.Name)
	}

	var err error
	builder.Object, err = builder.apiClient.ClusterRoles().Update(
		context.TODO(), builder.Definition, metav1.UpdateOptions{})

	return builder, err
}

// Exists checks if a clusterrole exists in the cluster.
func (builder *ClusterRoleBuilder) Exists() bool {
	if valid, _ := builder.validate(); !valid {
		return false
	}

	glog.V(100).Infof("Checking if clusterrole %s exists",
		builder.Definition.Name)

	var err error
	builder.Object, err = builder.apiClient.ClusterRoles().Get(
		context.TODO(), builder.Definition.Name, metav1.GetOptions{})

	return err == nil || !k8serrors.IsNotFound(err)
}

// validate will check that the builder and builder definition are properly initialized before
// accessing any member fields.
func (builder *ClusterRoleBuilder) validate() (bool, error) {
	resourceCRD := "clusterRole"

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
