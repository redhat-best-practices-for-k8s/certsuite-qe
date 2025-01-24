package olm

import (
	"context"
	"fmt"

	"github.com/golang/glog"
	"github.com/openshift-kni/eco-goinfra/pkg/clients"
	"github.com/openshift-kni/eco-goinfra/pkg/msg"
	operatorsV1alpha1 "github.com/openshift-kni/eco-goinfra/pkg/schemes/olm/operators/v1alpha1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtimeClient "sigs.k8s.io/controller-runtime/pkg/client"
)

// SubscriptionBuilder provides a struct for Subscription object containing connection to the
// cluster and the Subscription definition.
type SubscriptionBuilder struct {
	// Subscription definition. Used to create Subscription object with minimum set of required elements.
	Definition *operatorsV1alpha1.Subscription
	// Created Subscription object on the cluster.
	Object *operatorsV1alpha1.Subscription
	// api client to interact with the cluster.
	apiClient runtimeClient.Client
	// errorMsg is processed before Subscription object is created.
	errorMsg string
}

// NewSubscriptionBuilder returns a SubscriptionBuilder.
func NewSubscriptionBuilder(apiClient *clients.Settings, subName, subNamespace, catalogSource, catalogSourceNamespace,
	packageName string) *SubscriptionBuilder {
	glog.V(100).Infof(
		"Initializing new SubscriptionBuilder structure with the following params, subName: %s, "+
			"subNamespace: %s, catalogSource: %s, catalogSourceNamespace: %s, packageName: %s ",
		subName, subNamespace, catalogSource, catalogSourceNamespace, packageName)

	if apiClient == nil {
		glog.V(100).Infof("The apiClient cannot be nil")

		return nil
	}

	err := apiClient.AttachScheme(operatorsV1alpha1.AddToScheme)
	if err != nil {
		glog.V(100).Infof("Failed to add operatorsV1alpha1 scheme to client schemes")

		return nil
	}

	builder := &SubscriptionBuilder{
		apiClient: apiClient.Client,
		Definition: &operatorsV1alpha1.Subscription{
			ObjectMeta: metav1.ObjectMeta{
				Name:      subName,
				Namespace: subNamespace,
			},
			Spec: &operatorsV1alpha1.SubscriptionSpec{
				CatalogSource:          catalogSource,
				CatalogSourceNamespace: catalogSourceNamespace,
				Package:                packageName,
			},
		},
	}

	if subName == "" {
		glog.V(100).Infof("The Name of the Subscription is empty")

		builder.errorMsg = "subscription 'subName' cannot be empty"

		return builder
	}

	if subNamespace == "" {
		glog.V(100).Infof("The Namespace of the Subscription is empty")

		builder.errorMsg = "subscription 'subNamespace' cannot be empty"

		return builder
	}

	if catalogSource == "" {
		glog.V(100).Infof("The Catalogsource of the Subscription is empty")

		builder.errorMsg = "subscription 'catalogSource' cannot be empty"

		return builder
	}

	if catalogSourceNamespace == "" {
		glog.V(100).Infof("The Catalogsource namespace of the Subscription is empty")

		builder.errorMsg = "subscription 'catalogSourceNamespace' cannot be empty"

		return builder
	}

	if packageName == "" {
		glog.V(100).Infof("The Package name of the Subscription is empty")

		builder.errorMsg = "subscription 'packageName' cannot be empty"

		return builder
	}

	return builder
}

// WithChannel adds the specific channel to the Subscription.
func (builder *SubscriptionBuilder) WithChannel(channel string) *SubscriptionBuilder {
	if valid, _ := builder.validate(); !valid {
		return builder
	}

	glog.V(100).Infof("Defining Subscription builder object with channel: %s", channel)

	if channel == "" {
		builder.errorMsg = "can not redefine subscription with empty channel"

		return builder
	}

	builder.Definition.Spec.Channel = channel

	return builder
}

// WithStartingCSV adds the specific startingCSV to the Subscription.
func (builder *SubscriptionBuilder) WithStartingCSV(startingCSV string) *SubscriptionBuilder {
	if valid, _ := builder.validate(); !valid {
		return builder
	}

	glog.V(100).Infof("Defining Subscription builder object with startingCSV: %s",
		startingCSV)

	if startingCSV == "" {
		builder.errorMsg = "can not redefine subscription with empty startingCSV"

		return builder
	}

	builder.Definition.Spec.StartingCSV = startingCSV

	return builder
}

// WithInstallPlanApproval adds the specific installPlanApproval to the Subscription.
func (builder *SubscriptionBuilder) WithInstallPlanApproval(
	installPlanApproval operatorsV1alpha1.Approval) *SubscriptionBuilder {
	if valid, _ := builder.validate(); !valid {
		return builder
	}

	glog.V(100).Infof("Defining Subscription builder object with "+
		"installPlanApproval: %s", installPlanApproval)

	if !(installPlanApproval == "Automatic" || installPlanApproval == "Manual") {
		glog.V(100).Infof("The InstallPlanApproval of the Subscription must be either \"Automatic\" " +
			"or \"Manual\"")

		builder.errorMsg = "Subscription 'installPlanApproval' must be either \"Automatic\" or \"Manual\""

		return builder
	}

	builder.Definition.Spec.InstallPlanApproval = installPlanApproval

	return builder
}

// Get returns Subscription object if found.
func (builder *SubscriptionBuilder) Get() (*operatorsV1alpha1.Subscription, error) {
	if valid, err := builder.validate(); !valid {
		return nil, err
	}

	glog.V(100).Infof(
		"Collecting Subscription object %s in namespace %s",
		builder.Definition.Name, builder.Definition.Namespace)

	subscription := &operatorsV1alpha1.Subscription{}
	err := builder.apiClient.Get(context.TODO(),
		runtimeClient.ObjectKey{Name: builder.Definition.Name, Namespace: builder.Definition.Namespace},
		subscription)

	if err != nil {
		glog.V(100).Infof(
			"Subscription object %s does not exist in namespace %s",
			builder.Definition.Name, builder.Definition.Namespace)

		return nil, err
	}

	return subscription, nil
}

// Create makes an Subscription in cluster and stores the created object in struct.
func (builder *SubscriptionBuilder) Create() (*SubscriptionBuilder, error) {
	if valid, err := builder.validate(); !valid {
		return builder, err
	}

	glog.V(100).Infof("Creating the Subscription %s in namespace %s",
		builder.Definition.Name, builder.Definition.Namespace)

	if builder.Exists() {
		return builder, nil
	}

	err := builder.apiClient.Create(context.TODO(), builder.Definition)
	if err != nil {
		return builder, err
	}

	builder.Object = builder.Definition

	return builder, nil
}

// Exists checks whether the given Subscription exists.
func (builder *SubscriptionBuilder) Exists() bool {
	if valid, _ := builder.validate(); !valid {
		return false
	}

	glog.V(100).Infof(
		"Checking if Subscription %s exists",
		builder.Definition.Name)

	var err error
	builder.Object, err = builder.Get()

	return err == nil || !k8serrors.IsNotFound(err)
}

// Delete removes a Subscription.
func (builder *SubscriptionBuilder) Delete() error {
	if valid, err := builder.validate(); !valid {
		return err
	}

	glog.V(100).Infof("Deleting Subscription %s in namespace %s", builder.Definition.Name,
		builder.Definition.Namespace)

	if !builder.Exists() {
		glog.V(100).Infof("Subscription object %s does not exist in namespace %s",
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

// Update modifies the existing Subscription with the Subscription definition in SubscriptionBuilder.
func (builder *SubscriptionBuilder) Update() (*SubscriptionBuilder, error) {
	if valid, err := builder.validate(); !valid {
		return builder, err
	}

	glog.V(100).Infof("Updating Subscription %s in namespace %s",
		builder.Definition.Name, builder.Definition.Namespace)

	if !builder.Exists() {
		return nil, fmt.Errorf("subscription named %s in namespace %s does not exist",
			builder.Definition.Name, builder.Definition.Namespace)
	}

	err := builder.apiClient.Update(context.TODO(), builder.Definition)

	if err == nil {
		builder.Object = builder.Definition
	}

	return builder, err
}

// PullSubscription loads existing Subscription from cluster into the SubscriptionBuilder struct.
func PullSubscription(apiClient *clients.Settings, subName, subNamespace string) (*SubscriptionBuilder, error) {
	glog.V(100).Infof("Pulling existing Subscription %s from cluster in namespace %s",
		subName, subNamespace)

	if apiClient == nil {
		glog.V(100).Infof("The apiClient cannot be nil")

		return nil, fmt.Errorf("subscription 'apiClient' cannot be empty")
	}

	err := apiClient.AttachScheme(operatorsV1alpha1.AddToScheme)
	if err != nil {
		glog.V(100).Infof("Failed to add operatorsV1alpha1 scheme to client schemes")

		return nil, err
	}

	builder := &SubscriptionBuilder{
		apiClient: apiClient.Client,
		Definition: &operatorsV1alpha1.Subscription{
			ObjectMeta: metav1.ObjectMeta{
				Name:      subName,
				Namespace: subNamespace,
			},
		},
	}

	if subName == "" {
		glog.V(100).Infof("The name of the Subscription is empty")

		return nil, fmt.Errorf("subscription 'subName' cannot be empty")
	}

	if subNamespace == "" {
		glog.V(100).Infof("The namespace of the Subscription is empty")

		return nil, fmt.Errorf("subscription 'subNamespace' cannot be empty")
	}

	if !builder.Exists() {
		return nil, fmt.Errorf(
			"subscription object named %s does not exist in namespace %s", subName, subNamespace)
	}

	builder.Definition = builder.Object

	return builder, nil
}

// validate will check that the builder and builder definition are properly initialized before
// accessing any member fields.
func (builder *SubscriptionBuilder) validate() (bool, error) {
	resourceCRD := "Subscription"

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
