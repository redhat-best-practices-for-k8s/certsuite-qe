package olm

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"

	"github.com/golang/glog"
	"github.com/openshift-kni/eco-goinfra/pkg/clients"
	"github.com/openshift-kni/eco-goinfra/pkg/msg"
	operatorv1 "github.com/openshift-kni/eco-goinfra/pkg/schemes/olm/package-server/operators/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtimeClient "sigs.k8s.io/controller-runtime/pkg/client"
)

// PackageManifestBuilder provides a struct for PackageManifest object from the cluster
// and a PackageManifest definition.
type PackageManifestBuilder struct {
	// PackageManifest definition. Used to create
	// PackageManifest object with minimum set of required elements.
	Definition *operatorv1.PackageManifest
	// Created PackageManifest object on the cluster.
	Object *operatorv1.PackageManifest
	// api client to interact with the cluster.
	apiClient runtimeClient.Client
	// errorMsg is processed before PackageManifest object is created.
	errorMsg string
}

// PullPackageManifest loads an existing PackageManifest into Builder struct.
func PullPackageManifest(apiClient *clients.Settings, name, nsname string) (*PackageManifestBuilder, error) {
	glog.V(100).Infof("Pulling existing PackageManifest name %s in namespace %s", name, nsname)

	if apiClient == nil {
		glog.V(100).Infof("The apiClient cannot be nil")

		return nil, fmt.Errorf("packagemanifest 'apiClient' cannot be empty")
	}

	err := apiClient.AttachScheme(operatorv1.AddToScheme)

	if err != nil {
		glog.V(100).Infof("Failed to add operatorv1 scheme to client schemes")

		return nil, err
	}

	builder := &PackageManifestBuilder{
		apiClient: apiClient.Client,
		Definition: &operatorv1.PackageManifest{
			ObjectMeta: metav1.ObjectMeta{
				Name:      name,
				Namespace: nsname,
			},
		},
	}

	if name == "" {
		glog.V(100).Infof("The Name of the PackageManifest is empty")

		return nil, fmt.Errorf("packageManifest 'name' cannot be empty")
	}

	if nsname == "" {
		glog.V(100).Infof("The Namespace of the PackageManifest is empty")

		return nil, fmt.Errorf("packageManifest 'nsname' cannot be empty")
	}

	if !builder.Exists() {
		return nil, fmt.Errorf("packageManifest object %s does not exist in namespace %s", name, nsname)
	}

	builder.Definition = builder.Object

	return builder, nil
}

// PullPackageManifestByCatalog loads an existing PackageManifest from specified catalog into Builder struct.
func PullPackageManifestByCatalog(apiClient *clients.Settings, name, nsname,
	catalog string) (*PackageManifestBuilder, error) {
	glog.V(100).Infof("Pulling existing PackageManifest name %s in namespace %s and from catalog %s",
		name, nsname, catalog)

	if name == "" {
		glog.V(100).Infof("The Name of the PackageManifest is empty")

		return nil, fmt.Errorf("packageManifest 'name' cannot be empty")
	}

	fieldSelector, err := fields.ParseSelector(fmt.Sprintf("metadata.name=%s", name))

	if err != nil {
		glog.V(100).Infof("Failed to parse invalid packageManifest name %s", name)

		return nil, err
	}

	if catalog == "" {
		glog.V(100).Infof("The Catalog of the PackageManifest is empty")

		return nil, fmt.Errorf("packageManifest 'catalog' cannot be empty")
	}

	labelSelector, err := labels.Parse(fmt.Sprintf("catalog=%s", catalog))

	if err != nil {
		glog.V(100).Infof("Failed to parse invalid catalog name %s", catalog)

		return nil, err
	}

	packageManifests, err := ListPackageManifest(apiClient, nsname, runtimeClient.ListOptions{
		LabelSelector: labelSelector,
		FieldSelector: fieldSelector,
	})

	if err != nil {
		glog.V(100).Infof("Failed to list PackageManifests with name %s in namespace %s from catalog"+
			" %s due to %s", name, nsname, catalog, err.Error())

		return nil, err
	}

	if len(packageManifests) == 0 {
		glog.V(100).Infof("The list of matching PackageManifests is empty")

		return nil, fmt.Errorf("no matching PackageManifests were found")
	}

	if len(packageManifests) > 1 {
		glog.V(100).Infof("More than one matching PackageManifests were found")

		return nil, fmt.Errorf("more than one matching PackageManifests were found")
	}

	return packageManifests[0], nil
}

// Get returns PackageManifest object if found.
func (builder *PackageManifestBuilder) Get() (*operatorv1.PackageManifest, error) {
	if valid, err := builder.validate(); !valid {
		return nil, err
	}

	glog.V(100).Infof(
		"Collecting packageManifest object %s in namespace %s",
		builder.Definition.Name, builder.Definition.Namespace)

	packageManifest := &operatorv1.PackageManifest{}
	err := builder.apiClient.Get(context.TODO(),
		runtimeClient.ObjectKey{Name: builder.Definition.Name, Namespace: builder.Definition.Namespace},
		packageManifest)

	if err != nil {
		glog.V(100).Infof(
			"PackageManifest object %s does not exist in namespace %s",
			builder.Definition.Name, builder.Definition.Namespace)

		return nil, err
	}

	return packageManifest, nil
}

// Exists checks whether the given PackageManifest exists.
func (builder *PackageManifestBuilder) Exists() bool {
	if valid, _ := builder.validate(); !valid {
		return false
	}

	glog.V(100).Infof(
		"Checking if PackageManifest %s exists", builder.Definition.Name)

	var err error
	builder.Object, err = builder.Get()

	return err == nil || !k8serrors.IsNotFound(err)
}

// Delete removes a PackageManifest.
func (builder *PackageManifestBuilder) Delete() error {
	if valid, err := builder.validate(); !valid {
		return err
	}

	glog.V(100).Infof("Deleting PackageManifest %s in namespace %s", builder.Definition.Name,
		builder.Definition.Namespace)

	if !builder.Exists() {
		glog.V(100).Infof("PackageManifest object %s does not exist in namespace %s",
			builder.Definition.Name, builder.Definition.Namespace)

		builder.Object = nil

		return nil
	}

	err := builder.apiClient.Delete(context.TODO(), builder.Definition)

	if err != nil {
		return err
	}

	builder.Object = nil

	return nil
}

// validate will check that the builder and builder definition are properly initialized before
// accessing any member fields.
func (builder *PackageManifestBuilder) validate() (bool, error) {
	resourceCRD := "PackageManifest"

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
