package observabilityhelper

import (
	"fmt"
	"strings"

	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/observability/observabilityparameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/crd"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/daemonset"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/deployment"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/pod"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/statefulset"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	apiextv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
)

// For some reason, there's a function that expects labels' key/values separated
// by colon instead of the equal char.
func GetTnfTargetPodLabelsSlice() []string {
	return []string{observabilityparameters.TestPodLabelKey + ":" + observabilityparameters.TestPodLabelValue}
}

// DefineDeploymentWithStdoutBuffers defines a deployment with a given name and replicas number, creating
// a container spec for each entry in the stdoutBuffers slice. The number of containers is
// not needed as a parameter, as it will match the number of entries in stdoutBuffers.
// There are equivalent functions for statefulSets and daemonSets.
func DefineDeploymentWithStdoutBuffers(name string, replicas int, stdoutBuffers []string) *appsv1.Deployment {
	// Get each container spec based on the desired stdout buffer per container.
	containerSpecs := createContainerSpecsFromStdoutBuffers(stdoutBuffers)

	return defineDeploymentWithContainerSpecs(name, replicas, containerSpecs)
}

func DefineStatefulSetWithStdoutBuffers(name string, replicas int, stdoutBuffers []string) *appsv1.StatefulSet {
	// Get each container spec based on the desired stdout buffer per container.
	containerSpecs := createContainerSpecsFromStdoutBuffers(stdoutBuffers)

	return defineStatefulSetWithContainerSpecs(name, replicas, containerSpecs)
}

func DefineDaemonSetWithStdoutBuffers(name string, stdoutBuffers []string) *appsv1.DaemonSet {
	// Get each container spec based on the desired stdout buffer per container.
	containerSpecs := createContainerSpecsFromStdoutBuffers(stdoutBuffers)

	// Define base daemonset
	return defineDaemonSetWithContainerSpecs(name, containerSpecs)
}

func DefinePodWithStdoutBuffer(name string, stdoutBuffer string) *corev1.Pod {
	newPod := pod.DefinePod(name, observabilityparameters.TestNamespace, globalhelper.Configuration.General.TestImage)
	// Add labels.
	newPod = pod.RedefinePodWithLabel(newPod, observabilityparameters.TnfTargetPodLabels)
	// Change command to use the stdout buffer.
	newPod.Spec.Containers[0].Command = getContainerCommandWithStdout(stdoutBuffer)

	return newPod
}

func DefineDeploymentWithoutTargetLabels(name string) *appsv1.Deployment {
	return deployment.DefineDeployment(name, observabilityparameters.TestNamespace,
		globalhelper.Configuration.General.TestImage,
		map[string]string{"fakeLabelKey": "fakeLabelValue"})
}

func DefineCrdWithStatusSubresource(kind, group string) *apiextv1.CustomResourceDefinition {
	return crd.DefineCustomResourceDefinition(apiextv1.CustomResourceDefinitionNames{
		Kind:     kind,
		Singular: strings.ToLower(kind),
		Plural:   strings.ToLower(kind) + "s",
		ListKind: kind + "List",
	}, group, true)
}

func DefineCrdWithoutStatusSubresource(kind, group string) *apiextv1.CustomResourceDefinition {
	return crd.DefineCustomResourceDefinition(apiextv1.CustomResourceDefinitionNames{
		Kind:     kind,
		Singular: strings.ToLower(kind),
		Plural:   strings.ToLower(kind) + "s",
		ListKind: kind + "List",
	}, group, false)
}

// getContainerCommandWithStdout is a helper function that will return the command slice
// to be used in a container spec. The command will call bash to execute printf <stdout>
// followed by an infinite sleep. The text should be whatever the TC needs to print in the
// container output to pass/fail the TC.
func getContainerCommandWithStdout(stdout string) []string {
	return []string{"/bin/bash", "-c", fmt.Sprintf("printf %q && sleep INF", stdout)}
}

// createContainerSpecsFromStdoutBuffers is a helper function that creates a container spec for
// each entry in the stdoutBuffers slice.
func createContainerSpecsFromStdoutBuffers(stdoutBuffers []string) []corev1.Container {
	numContainers := len(stdoutBuffers)
	containerSpecs := []corev1.Container{}

	for index := 0; index < numContainers; index++ {
		stdoutLines := stdoutBuffers[index]

		containerSpecs = append(containerSpecs,
			corev1.Container{
				Name:    fmt.Sprintf("%s-%d", observabilityparameters.TestContainerBaseName, index),
				Image:   globalhelper.Configuration.General.TestImage,
				Command: getContainerCommandWithStdout(stdoutLines),
			},
		)
	}

	return containerSpecs
}

func defineDeploymentWithContainerSpecs(name string, replicas int,
	containerSpecs []corev1.Container) *appsv1.Deployment {
	// Define base deployment
	dep := deployment.DefineDeployment(name, observabilityparameters.TestNamespace,
		globalhelper.Configuration.General.TestImage, observabilityparameters.TnfTargetPodLabels)

	// Customize its replicas and container specs.
	dep = deployment.RedefineWithReplicaNumber(dep, int32(replicas))
	dep = deployment.RedefineWithContainerSpecs(dep, containerSpecs)

	return dep
}

func defineStatefulSetWithContainerSpecs(name string, replicas int,
	containerSpecs []corev1.Container) *appsv1.StatefulSet {
	// Define base statefulSet
	sts := statefulset.DefineStatefulSet(name, observabilityparameters.TestNamespace,
		globalhelper.Configuration.General.TestImage, observabilityparameters.TnfTargetPodLabels)

	// Customize its replicas and container specs.
	sts = statefulset.RedefineWithReplicaNumber(sts, int32(replicas))
	sts = statefulset.RedefineWithContainerSpecs(sts, containerSpecs)

	return sts
}

func defineDaemonSetWithContainerSpecs(name string,
	containerSpecs []corev1.Container) *appsv1.DaemonSet {
	// Define base daemonSet
	daemonSet := daemonset.DefineDaemonSet(observabilityparameters.TestNamespace,
		globalhelper.Configuration.General.TestImage, observabilityparameters.TnfTargetPodLabels, name)

	// Customize its container specs.
	return daemonset.RedefineWithContainerSpecs(daemonSet, containerSpecs)
}
