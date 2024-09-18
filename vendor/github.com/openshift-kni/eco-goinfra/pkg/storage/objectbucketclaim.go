package storage

import (
	"context"
	"fmt"

	"github.com/golang/glog"
	"github.com/openshift-kni/eco-goinfra/pkg/clients"
	"github.com/openshift-kni/eco-goinfra/pkg/msg"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	goclient "sigs.k8s.io/controller-runtime/pkg/client"

	noobaav1alpha1 "github.com/kube-object-storage/lib-bucket-provisioner/pkg/apis/objectbucket.io/v1alpha1"
)

// ObjectBucketClaimBuilder provides struct for the objectBucketClaim object.
type ObjectBucketClaimBuilder struct {
	// ObjectBucketClaim definition. Used to create objectBucketClaim object with minimum set of required elements.
	Definition *noobaav1alpha1.ObjectBucketClaim
	// Created objectBucketClaim object on the cluster.
	Object *noobaav1alpha1.ObjectBucketClaim
	// api client to interact with the cluster.
	apiClient goclient.Client
	// errorMsg is processed before objectBucketClaim object is created.
	errorMsg string
}

// NewObjectBucketClaimBuilder creates new instance of builder.
func NewObjectBucketClaimBuilder(
	apiClient *clients.Settings, name, nsname string) *ObjectBucketClaimBuilder {
	glog.V(100).Infof("Initializing new objectBucketClaim structure with the following params: "+
		"name: %s, namespace: %s", name, nsname)

	if apiClient == nil {
		glog.V(100).Infof("objectBucketClaim 'apiClient' cannot be empty")

		return nil
	}

	err := apiClient.AttachScheme(noobaav1alpha1.AddToScheme)
	if err != nil {
		glog.V(100).Infof("Failed to add objectbucket.io/v1alpha1 scheme to client schemes")

		return nil
	}

	builder := &ObjectBucketClaimBuilder{
		apiClient: apiClient.Client,
		Definition: &noobaav1alpha1.ObjectBucketClaim{
			ObjectMeta: metav1.ObjectMeta{
				Name:      name,
				Namespace: nsname,
			},
		},
	}

	if name == "" {
		glog.V(100).Infof("The name of the objectBucketClaim is empty")

		builder.errorMsg = "objectBucketClaim 'name' cannot be empty"
	}

	if nsname == "" {
		glog.V(100).Infof("The nsname of the objectBucketClaim is empty")

		builder.errorMsg = "objectBucketClaim 'nsname' cannot be empty"
	}

	return builder
}

// PullObjectBucketClaim retrieves an existing objectBucketClaim object from the cluster.
func PullObjectBucketClaim(apiClient *clients.Settings, name, nsname string) (*ObjectBucketClaimBuilder, error) {
	glog.V(100).Infof(
		"Pulling objectBucketClaim object name:%s in namespace: %s", name, nsname)

	if apiClient == nil {
		glog.V(100).Infof("The apiClient is empty")

		return nil, fmt.Errorf("objectBucketClaim 'apiClient' cannot be empty")
	}

	err := apiClient.AttachScheme(noobaav1alpha1.AddToScheme)
	if err != nil {
		glog.V(100).Infof("Failed to add objectbucket.io/v1alpha1 scheme to client schemes")

		return nil, err
	}

	builder := ObjectBucketClaimBuilder{
		apiClient: apiClient.Client,
		Definition: &noobaav1alpha1.ObjectBucketClaim{
			ObjectMeta: metav1.ObjectMeta{
				Name:      name,
				Namespace: nsname,
			},
		},
	}

	if name == "" {
		glog.V(100).Infof("The name of the objectBucketClaim is empty")

		return nil, fmt.Errorf("objectBucketClaim 'name' cannot be empty")
	}

	if nsname == "" {
		glog.V(100).Infof("The namespace of the objectBucketClaim is empty")

		return nil, fmt.Errorf("objectBucketClaim 'nsname' cannot be empty")
	}

	if !builder.Exists() {
		return nil, fmt.Errorf("objectBucketClaim object %s does not exist in namespace %s", name, nsname)
	}

	builder.Definition = builder.Object

	return &builder, nil
}

// Get returns objectBucketClaim object if found.
func (builder *ObjectBucketClaimBuilder) Get() (*noobaav1alpha1.ObjectBucketClaim, error) {
	if valid, err := builder.validate(); !valid {
		return nil, err
	}

	glog.V(100).Infof("Getting objectBucketClaim %s in namespace %s",
		builder.Definition.Name, builder.Definition.Namespace)

	objectBucketClaimObj := &noobaav1alpha1.ObjectBucketClaim{}
	err := builder.apiClient.Get(context.TODO(), goclient.ObjectKey{
		Name:      builder.Definition.Name,
		Namespace: builder.Definition.Namespace,
	}, objectBucketClaimObj)

	if err != nil {
		return nil, err
	}

	return objectBucketClaimObj, err
}

// Create makes a objectBucketClaim in the cluster and stores the created object in struct.
func (builder *ObjectBucketClaimBuilder) Create() (*ObjectBucketClaimBuilder, error) {
	if valid, err := builder.validate(); !valid {
		return builder, err
	}

	glog.V(100).Infof("Creating the objectBucketClaim %s in namespace %s",
		builder.Definition.Name, builder.Definition.Namespace)

	var err error
	if !builder.Exists() {
		err = builder.apiClient.Create(context.TODO(), builder.Definition)
		if err == nil {
			builder.Object = builder.Definition
		}
	}

	return builder, err
}

// Delete removes objectBucketClaim from a cluster.
func (builder *ObjectBucketClaimBuilder) Delete() error {
	if valid, err := builder.validate(); !valid {
		return err
	}

	glog.V(100).Infof("Deleting the objectBucketClaim %s from namespace %s",
		builder.Definition.Name, builder.Definition.Namespace)

	if !builder.Exists() {
		glog.V(100).Infof("objectBucketClaim %s in namespace %s cannot be deleted"+
			" because it does not exist",
			builder.Definition.Name, builder.Definition.Namespace)

		builder.Object = nil

		return nil
	}

	err := builder.apiClient.Delete(context.TODO(), builder.Definition)

	if err != nil {
		return fmt.Errorf("can not delete objectBucketClaim: %w", err)
	}

	builder.Object = nil

	return nil
}

// Exists checks whether the given objectBucketClaim exists.
func (builder *ObjectBucketClaimBuilder) Exists() bool {
	if valid, _ := builder.validate(); !valid {
		return false
	}

	glog.V(100).Infof("Checking if objectBucketClaim %s exists in namespace %s",
		builder.Definition.Name, builder.Definition.Namespace)

	var err error
	builder.Object, err = builder.Get()

	return err == nil || !k8serrors.IsNotFound(err)
}

// Update renovates the existing objectBucketClaim object with objectBucketClaim definition in builder.
func (builder *ObjectBucketClaimBuilder) Update() (*ObjectBucketClaimBuilder, error) {
	if valid, err := builder.validate(); !valid {
		return builder, err
	}

	glog.V(100).Info("Updating objectBucketClaim %s in namespace %s",
		builder.Definition.Name, builder.Definition.Namespace)

	err := builder.apiClient.Update(context.TODO(), builder.Definition)

	if err != nil {
		glog.V(100).Infof(
			msg.FailToUpdateError("objectBucketClaim", builder.Definition.Name, builder.Definition.Namespace))

		return nil, err
	}

	builder.Object = builder.Definition

	return builder, err
}

// WithStorageClassName sets the objectBucketClaim operator's storageClassName configuration.
func (builder *ObjectBucketClaimBuilder) WithStorageClassName(
	storageClassName string) *ObjectBucketClaimBuilder {
	if valid, _ := builder.validate(); !valid {
		return builder
	}

	glog.V(100).Infof(
		"Setting objectBucketClaim %s in namespace %s with the storageClassName config: %v",
		builder.Definition.Name, builder.Definition.Namespace, storageClassName)

	if storageClassName == "" {
		glog.V(100).Infof("'storageClassName' argument cannot be empty")

		builder.errorMsg = "'storageClassName' argument cannot be empty"

		return builder
	}

	builder.Definition.Spec.StorageClassName = storageClassName

	return builder
}

// WithGenerateBucketName sets the objectBucketClaim operator's generateBucketName configuration.
func (builder *ObjectBucketClaimBuilder) WithGenerateBucketName(
	generateBucketName string) *ObjectBucketClaimBuilder {
	if valid, _ := builder.validate(); !valid {
		return builder
	}

	glog.V(100).Infof(
		"Setting objectBucketClaim %s in namespace %s with the generateBucketName config: %v",
		builder.Definition.Name, builder.Definition.Namespace, generateBucketName)

	if generateBucketName == "" {
		glog.V(100).Infof("'generateBucketName' argument cannot be empty")

		builder.errorMsg = "'generateBucketName' argument cannot be empty"

		return builder
	}

	builder.Definition.Spec.GenerateBucketName = generateBucketName

	return builder
}

// validate will check that the builder and builder definition are properly initialized before
// accessing any member fields.
func (builder *ObjectBucketClaimBuilder) validate() (bool, error) {
	resourceCRD := "ObjectBucketClaim"

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
