package storage

import (
	"context"
	"fmt"

	goclient "sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/golang/glog"
	"github.com/openshift-kni/eco-goinfra/pkg/clients"
	"github.com/openshift-kni/eco-goinfra/pkg/msg"
	odfoperatorv1alpha1 "github.com/red-hat-storage/odf-operator/api/v1alpha1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// SystemODFBuilder provides struct for SystemODF object containing connection
// to the cluster and the SystemODF definitions.
type SystemODFBuilder struct {
	// SystemODF definition. Used to create a SystemODF object
	Definition *odfoperatorv1alpha1.StorageSystem
	// Created SystemODF object
	Object *odfoperatorv1alpha1.StorageSystem
	// api client to interact with the cluster.
	apiClient goclient.Client
	// Used in functions that define or mutate SystemODF definition. errorMsg is processed before the
	// SystemODF object is created.
	errorMsg string
}

// NewSystemODFBuilder creates a new instance of Builder.
func NewSystemODFBuilder(apiClient *clients.Settings, name, nsname string) *SystemODFBuilder {
	glog.V(100).Infof(
		"Initializing new SystemODF structure with the following params: %s, %s", name, nsname)

	if apiClient == nil {
		glog.V(100).Infof("SystemODF 'apiClient' cannot be empty")

		return nil
	}

	err := apiClient.AttachScheme(odfoperatorv1alpha1.AddToScheme)
	if err != nil {
		glog.V(100).Infof("Failed to add odf v1alpha1 scheme to client schemes")

		return nil
	}

	builder := &SystemODFBuilder{
		apiClient: apiClient.Client,
		Definition: &odfoperatorv1alpha1.StorageSystem{
			ObjectMeta: metav1.ObjectMeta{
				Name:      name,
				Namespace: nsname,
			},
		},
	}

	if name == "" {
		glog.V(100).Infof("The name of the SystemODF is empty")

		builder.errorMsg = "SystemODF 'name' cannot be empty"

		return builder
	}

	if nsname == "" {
		glog.V(100).Infof("The namespace of the SystemODF is empty")

		builder.errorMsg = "SystemODF 'nsname' cannot be empty"

		return builder
	}

	return builder
}

// PullSystemODF gets an existing SystemODF object from the cluster.
func PullSystemODF(apiClient *clients.Settings, name, namespace string) (*SystemODFBuilder, error) {
	glog.V(100).Infof("Pulling existing SystemODF object %s from namespace %s",
		name, namespace)

	if apiClient == nil {
		glog.V(100).Infof("The SystemODF's apiClient is empty")

		return nil, fmt.Errorf("SystemODF 'apiClient' cannot be empty")
	}

	err := apiClient.AttachScheme(odfoperatorv1alpha1.AddToScheme)
	if err != nil {
		glog.V(100).Infof("Failed to add odf v1alpha1 scheme to client schemes")

		return nil, err
	}

	builder := SystemODFBuilder{
		apiClient: apiClient.Client,
		Definition: &odfoperatorv1alpha1.StorageSystem{
			ObjectMeta: metav1.ObjectMeta{
				Name:      name,
				Namespace: namespace,
			},
		},
	}

	if name == "" {
		glog.V(100).Infof("The name of the SystemODF is empty")

		return nil, fmt.Errorf("SystemODF 'name' cannot be empty")
	}

	if namespace == "" {
		glog.V(100).Infof("The namespace of the SystemODF is empty")

		return nil, fmt.Errorf("SystemODF 'namespace' cannot be empty")
	}

	if !builder.Exists() {
		return nil, fmt.Errorf("SystemODF object %s does not exist in namespace %s",
			name, namespace)
	}

	builder.Definition = builder.Object

	return &builder, nil
}

// Get fetches existing SystemODF from cluster.
func (builder *SystemODFBuilder) Get() (*odfoperatorv1alpha1.StorageSystem, error) {
	if valid, err := builder.validate(); !valid {
		return nil, err
	}

	glog.V(100).Infof("Getting existing SystemODF with name %s from the namespace %s",
		builder.Definition.Name, builder.Definition.Namespace)

	storageSystemObj := &odfoperatorv1alpha1.StorageSystem{}
	err := builder.apiClient.Get(context.TODO(), goclient.ObjectKey{
		Name:      builder.Definition.Name,
		Namespace: builder.Definition.Namespace,
	}, storageSystemObj)

	if err != nil {
		glog.V(100).Infof("failed to find SystemODF object %s in namespace %s",
			builder.Definition.Name, builder.Definition.Namespace)

		return nil, err
	}

	return storageSystemObj, nil
}

// Exists checks whether the given SystemODF exists.
func (builder *SystemODFBuilder) Exists() bool {
	if valid, _ := builder.validate(); !valid {
		return false
	}

	glog.V(100).Infof("Checking if SystemODF %s exists in namespace %s",
		builder.Definition.Name, builder.Definition.Namespace)

	var err error
	builder.Object, err = builder.Get()

	return err == nil || !k8serrors.IsNotFound(err)
}

// Create makes a SystemODF in the cluster and stores the created object in struct.
func (builder *SystemODFBuilder) Create() (*SystemODFBuilder, error) {
	if valid, err := builder.validate(); !valid {
		return builder, err
	}

	glog.V(100).Infof("Creating the SystemODF %s in namespace %s",
		builder.Definition.Name, builder.Definition.Namespace,
	)

	var err error
	if !builder.Exists() {
		err = builder.apiClient.Create(context.TODO(), builder.Definition)
		if err == nil {
			builder.Object = builder.Definition
		}
	}

	return builder, err
}

// Delete removes SystemODF object from a cluster.
func (builder *SystemODFBuilder) Delete() error {
	if valid, err := builder.validate(); !valid {
		return err
	}

	glog.V(100).Infof("Deleting the SystemODF object %s in namespace %s",
		builder.Definition.Name, builder.Definition.Namespace)

	if !builder.Exists() {
		glog.V(100).Infof(" SystemODF %s in namespace %s cannot be deleted"+
			" because it does not exist",
			builder.Definition.Name, builder.Definition.Namespace)

		builder.Object = nil

		return nil
	}

	err := builder.apiClient.Delete(context.TODO(), builder.Definition)

	if err != nil {
		return fmt.Errorf("can not delete SystemODF: %w", err)
	}

	builder.Object = nil

	return nil
}

// WithSpec sets the SystemODF with storageCluster spec values.
func (builder *SystemODFBuilder) WithSpec(
	kind odfoperatorv1alpha1.StorageKind, name, nsname string) *SystemODFBuilder {
	if valid, _ := builder.validate(); !valid {
		return builder
	}

	glog.V(100).Infof(
		"Setting SystemODF %s in namespace %s with storageCluster spec; \n"+
			"kind: %v, name: %s, namespace %s",
		builder.Definition.Name, builder.Definition.Namespace, kind, name, nsname)

	if kind == "" {
		glog.V(100).Infof("The kind of the SystemODF spec is empty")

		builder.errorMsg = "SystemODF spec 'kind' cannot be empty"

		return builder
	}

	if name == "" {
		glog.V(100).Infof("The name of the storageCluster spec is empty")

		builder.errorMsg = "SystemODF spec 'name' cannot be empty"

		return builder
	}

	if nsname == "" {
		glog.V(100).Infof("The namespace of the SystemODF spec is empty")

		builder.errorMsg = "SystemODF spec 'nsname' cannot be empty"

		return builder
	}

	builder.Definition.Spec.Kind = kind
	builder.Definition.Spec.Name = name
	builder.Definition.Spec.Namespace = nsname

	return builder
}

// validate will check that the builder and builder definition are properly initialized before
// accessing any member fields.
func (builder *SystemODFBuilder) validate() (bool, error) {
	resourceCRD := "StorageSystem"

	if builder == nil {
		glog.V(100).Infof("The %s builder is uninitialized", resourceCRD)

		return false, fmt.Errorf("error: received nil %s builder", resourceCRD)
	}

	if builder.Definition == nil {
		glog.V(100).Infof("The %s is undefined", resourceCRD)

		return false, fmt.Errorf(msg.UndefinedCrdObjectErrString(resourceCRD))
	}

	if builder.apiClient == nil {
		glog.V(100).Infof("The %s builder apiclient is nil", resourceCRD)

		return false, fmt.Errorf("%s builder cannot have nil apiClient", resourceCRD)
	}

	if builder.errorMsg != "" {
		glog.V(100).Infof("The %s builder has error message: %s", resourceCRD, builder.errorMsg)

		return false, fmt.Errorf(builder.errorMsg)
	}

	return true, nil
}
