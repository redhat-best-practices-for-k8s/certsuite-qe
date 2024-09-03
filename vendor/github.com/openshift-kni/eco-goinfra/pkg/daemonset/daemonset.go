package daemonset

import (
	"context"
	"fmt"
	"time"

	"github.com/golang/glog"
	"github.com/openshift-kni/eco-goinfra/pkg/clients"
	"github.com/openshift-kni/eco-goinfra/pkg/msg"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	appsv1Typed "k8s.io/client-go/kubernetes/typed/apps/v1"
)

// Builder provides struct for daemonset object containing connection to the cluster and the daemonset definitions.
type Builder struct {
	// Daemonset definition. Used to create a daemonset object.
	Definition *appsv1.DaemonSet
	// Created daemonset object.
	Object *appsv1.DaemonSet
	// Used in functions that define or mutate daemonset definition. errorMsg is processed before the daemonset
	// object is created.
	errorMsg  string
	apiClient appsv1Typed.DaemonSetInterface
}

// AdditionalOptions additional options for daemonset object.
type AdditionalOptions func(builder *Builder) (*Builder, error)

var retryInterval = time.Second * 3

// NewBuilder creates a new instance of Builder.
func NewBuilder(
	apiClient *clients.Settings, name, nsname string, labels map[string]string, containerSpec corev1.Container) *Builder {
	glog.V(100).Infof(
		"Initializing new daemonset structure with the following params: "+
			"name: %s, namespace: %s, labels: %s, containerSpec %v",
		name, nsname, labels, containerSpec)

	if apiClient == nil {
		glog.V(100).Infof("The apiClient is nil")

		return nil
	}

	builder := &Builder{
		apiClient: apiClient.DaemonSets(nsname),
		Definition: &appsv1.DaemonSet{
			Spec: appsv1.DaemonSetSpec{
				Selector: &metav1.LabelSelector{
					MatchLabels: labels,
				},
				Template: corev1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Labels: labels,
					},
				},
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      name,
				Namespace: nsname,
			},
		},
	}

	builder.WithAdditionalContainerSpecs([]corev1.Container{containerSpec})

	if name == "" {
		glog.V(100).Infof("The name of the daemonset is empty")

		builder.errorMsg = "daemonset 'name' cannot be empty"

		return builder
	}

	if nsname == "" {
		glog.V(100).Infof("The namespace of the daemonset is empty")

		builder.errorMsg = "daemonset 'namespace' cannot be empty"

		return builder
	}

	if len(labels) == 0 {
		glog.V(100).Infof("There are no labels for the daemonset")

		builder.errorMsg = "daemonset 'labels' cannot be empty"

		return builder
	}

	return builder
}

// Pull loads an existing daemonSet into the Builder struct.
func Pull(apiClient *clients.Settings, name, nsname string) (*Builder, error) {
	glog.V(100).Infof("Pulling existing daemonset name:%s under namespace:%s", name, nsname)

	if apiClient == nil {
		glog.V(100).Infof("The apiClient is nil")

		return nil, fmt.Errorf("apiClient cannot be nil")
	}

	builder := &Builder{
		apiClient: apiClient.DaemonSets(nsname),
		Definition: &appsv1.DaemonSet{
			ObjectMeta: metav1.ObjectMeta{
				Name:      name,
				Namespace: nsname,
			},
		},
	}

	if name == "" {
		glog.V(100).Infof("The name of the daemonset is empty")

		return nil, fmt.Errorf("daemonset name cannot be empty")
	}

	if nsname == "" {
		glog.V(100).Infof("The namespace of the daemonset is empty")

		return nil, fmt.Errorf("daemonset namespace cannot be empty")
	}

	if !builder.Exists() {
		return nil, fmt.Errorf("daemonset object %s does not exist in namespace %s", name, nsname)
	}

	builder.Definition = builder.Object

	return builder, nil
}

// WithNodeSelector applies nodeSelector to the daemonset definition.
func (builder *Builder) WithNodeSelector(selector map[string]string) *Builder {
	if valid, _ := builder.validate(); !valid {
		return builder
	}

	glog.V(100).Infof("Applying nodeSelector %s to daemonset %s in namespace %s",
		selector, builder.Definition.Name, builder.Definition.Namespace)

	if len(selector) == 0 {
		glog.V(100).Infof("The nodeselector is empty")

		builder.errorMsg = "cannot accept empty map as nodeselector"

		return builder
	}

	builder.Definition.Spec.Template.Spec.NodeSelector = selector

	return builder
}

// WithHostNetwork applies HostNetwork to daemonset definition.
func (builder *Builder) WithHostNetwork() *Builder {
	if valid, _ := builder.validate(); !valid {
		return builder
	}

	glog.V(100).Infof("Enabling hostnetwork flag to daemonset %s in namespace %s",
		builder.Definition.Name, builder.Definition.Namespace)

	builder.Definition.Spec.Template.Spec.HostNetwork = true

	return builder
}

// WithVolume defines Volume of daemonset under PodTemplateSpec.
func (builder *Builder) WithVolume(dsVolume corev1.Volume) *Builder {
	if valid, _ := builder.validate(); !valid {
		return builder
	}

	if dsVolume.Name == "" {
		glog.V(100).Infof("The Volume name parameter is empty")

		builder.errorMsg = "Volume name parameter is empty"

		return builder
	}

	glog.V(100).Infof("Adding volume %s for daemonset %s pod template in namespace %s",
		dsVolume.Name, builder.Definition.Name, builder.Definition.Namespace)

	builder.Definition.Spec.Template.Spec.Volumes = append(
		builder.Definition.Spec.Template.Spec.Volumes,
		dsVolume)

	return builder
}

// WithAdditionalContainerSpecs appends a list of container specs to the daemonset definition.
func (builder *Builder) WithAdditionalContainerSpecs(specs []corev1.Container) *Builder {
	if valid, _ := builder.validate(); !valid {
		return builder
	}

	glog.V(100).Infof("Appending a list of container specs %v to daemonset %s in namespace %s",
		specs, builder.Definition.Name, builder.Definition.Namespace)

	if len(specs) == 0 {
		glog.V(100).Infof("The container specs are empty")

		builder.errorMsg = "cannot accept empty list as container specs"

		return builder
	}

	if builder.Definition.Spec.Template.Spec.Containers == nil {
		builder.Definition.Spec.Template.Spec.Containers = specs
	} else {
		builder.Definition.Spec.Template.Spec.Containers = append(builder.Definition.Spec.Template.Spec.Containers, specs...)
	}

	return builder
}

// WithOptions creates daemonset with generic mutation options.
func (builder *Builder) WithOptions(options ...AdditionalOptions) *Builder {
	if valid, _ := builder.validate(); !valid {
		return builder
	}

	glog.V(100).Infof("Setting daemonset additional options")

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

// Create builds daemonset in the cluster and stores the created object in struct.
func (builder *Builder) Create() (*Builder, error) {
	if valid, err := builder.validate(); !valid {
		return builder, err
	}

	glog.V(100).Infof("Creating daemonset %s in namespace %s", builder.Definition.Name, builder.Definition.Namespace)

	var err error
	if !builder.Exists() {
		builder.Object, err = builder.apiClient.Create(
			context.TODO(), builder.Definition, metav1.CreateOptions{})
	}

	return builder, err
}

// Update renovates the existing daemonset object with daemonset definition in builder.
func (builder *Builder) Update() (*Builder, error) {
	if valid, err := builder.validate(); !valid {
		return builder, err
	}

	glog.V(100).Infof("Updating daemonset %s in namespace %s", builder.Definition.Name, builder.Definition.Namespace)

	var err error
	builder.Object, err = builder.apiClient.Update(
		context.TODO(), builder.Definition, metav1.UpdateOptions{})

	return builder, err
}

// Delete removes the daemonset.
func (builder *Builder) Delete() error {
	if valid, err := builder.validate(); !valid {
		return err
	}

	glog.V(100).Infof("Deleting daemonset %s in namespace %s",
		builder.Definition.Name, builder.Definition.Namespace)

	if !builder.Exists() {
		builder.Object = nil

		return nil
	}

	err := builder.apiClient.Delete(
		context.TODO(), builder.Object.Name, metav1.DeleteOptions{})

	if err != nil {
		return err
	}

	builder.Object = nil

	return err
}

// CreateAndWaitUntilReady creates a daemonset in the cluster and waits until the daemonset is available.
func (builder *Builder) CreateAndWaitUntilReady(timeout time.Duration) (*Builder, error) {
	if valid, err := builder.validate(); !valid {
		return builder, err
	}

	glog.V(100).Infof("Creating daemonset %s in namespace %s and waiting for the defined period until it is ready",
		builder.Definition.Name, builder.Definition.Namespace)

	_, err := builder.Create()
	if err != nil {
		glog.V(100).Infof("Failed to create daemonset. Error is: '%s'", err.Error())

		return nil, fmt.Errorf(err.Error())
	}

	// Polls every retryInterval to determine if daemonset is available.
	err = wait.PollUntilContextTimeout(
		context.TODO(), retryInterval, timeout, true, func(ctx context.Context) (bool, error) {
			builder.Object, err = builder.apiClient.Get(
				context.TODO(), builder.Definition.Name, metav1.GetOptions{})

			if err != nil {
				return false, nil
			}

			for _, condition := range builder.Object.Status.Conditions {
				if condition.Type == "Available" {
					return condition.Status == "True", nil
				}
			}

			return false, err
		})

	if err == nil {
		return builder, nil
	}

	return nil, err
}

// DeleteAndWait deletes a daemonset and waits until it is removed from the cluster.
func (builder *Builder) DeleteAndWait(timeout time.Duration) error {
	if valid, err := builder.validate(); !valid {
		return err
	}

	glog.V(100).Infof("Deleting daemonset %s in namespace %s and waiting for the defined period until it is removed",
		builder.Definition.Name, builder.Definition.Namespace)

	if err := builder.Delete(); err != nil {
		return err
	}

	// Polls the daemonset every retryInterval until it is removed.
	return wait.PollUntilContextTimeout(
		context.TODO(), retryInterval, timeout, true, func(ctx context.Context) (bool, error) {
			_, err := builder.apiClient.Get(
				context.TODO(), builder.Definition.Name, metav1.GetOptions{})
			if k8serrors.IsNotFound(err) {
				return true, nil
			}

			return false, nil
		})
}

// Exists checks whether the given daemonset exists.
func (builder *Builder) Exists() bool {
	if valid, _ := builder.validate(); !valid {
		return false
	}

	glog.V(100).Infof("Checking if daemonset %s exists in namespace %s",
		builder.Definition.Name, builder.Definition.Namespace)

	var err error
	builder.Object, err = builder.apiClient.Get(
		context.TODO(), builder.Definition.Name, metav1.GetOptions{})

	return err == nil || !k8serrors.IsNotFound(err)
}

// IsReady waits for the daemonset to reach expected number of pods in Ready state.
func (builder *Builder) IsReady(timeout time.Duration) bool {
	if valid, _ := builder.validate(); !valid {
		return false
	}

	glog.V(100).Infof("Running periodic check until daemonset %s in namespace %s is ready or "+
		"timeout %s exceeded", builder.Definition.Name, builder.Definition.Namespace, timeout.String())

	// Polls every retryInterval to determine if daemonset is available.
	err := wait.PollUntilContextTimeout(
		context.TODO(), retryInterval, timeout, true, func(ctx context.Context) (bool, error) {
			if !builder.Exists() {
				return false, fmt.Errorf("daemonset %s is not present on cluster", builder.Object.Name)
			}

			var err error
			builder.Object, err = builder.apiClient.Get(
				context.TODO(), builder.Definition.Name, metav1.GetOptions{})

			if err != nil {
				glog.V(100).Infof("Failed to get daemonset from cluster. Error is: '%s'", err.Error())

				return false, nil
			}

			if builder.Object.Status.NumberReady == builder.Object.Status.DesiredNumberScheduled {
				return true, nil
			}

			if builder.Object.Status.NumberReady == builder.Object.Status.UpdatedNumberScheduled &&
				builder.Object.Status.UpdatedNumberScheduled != 0 {
				return true, nil
			}

			return false, err
		})

	return err == nil
}

// validate will check that the builder and builder definition are properly initialized before
// accessing any member fields.
func (builder *Builder) validate() (bool, error) {
	resourceCRD := "DaemonSet"

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
