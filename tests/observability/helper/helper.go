package helper

import (
	"fmt"

	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/globalparameters"
	params "github.com/test-network-function/cnfcert-tests-verification/tests/observability/parameters"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/daemonset"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/deployment"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/pod"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/statefulset"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

// RunXXX/ValidateXXX helpers to ease the reading of the TC.
func RunTnfPassingTestCase(tnfTcName string) {
	err := globalhelper.LaunchTests(
		tnfTcName,
		globalhelper.ConvertSpecNameToFileName(CurrentGinkgoTestDescription().FullTestText))
	Expect(err).ToNot(HaveOccurred())
}

func RunTnfFailingTestCase(tnfTcName string) {
	err := globalhelper.LaunchTests(
		tnfTcName,
		globalhelper.ConvertSpecNameToFileName(CurrentGinkgoTestDescription().FullTestText))
	Expect(err).To(HaveOccurred())
}

func ValidateTnfTcAsPassed(tnfTcName string) {
	err := globalhelper.ValidateIfReportsAreValid(
		tnfTcName,
		globalparameters.TestCasePassed)
	Expect(err).ToNot(HaveOccurred())
}

func ValidateTnfTcAsFailed(tnfTcName string) {
	err := globalhelper.ValidateIfReportsAreValid(
		tnfTcName,
		globalparameters.TestCaseFailed)
	Expect(err).ToNot(HaveOccurred())
}

func ValidateTnfTcAsSkipped(tnfTcName string) {
	err := globalhelper.ValidateIfReportsAreValid(
		tnfTcName,
		globalparameters.TestCaseSkipped)
	Expect(err).ToNot(HaveOccurred())
}

// For some reason, there's a function that expects labels' key/values separated
// by colon instead of the equal char.
func GetTnfTargetPodLabelsSlice() []string {
	return []string{params.QeTestPodLabelKey + ":" + params.QeTestPodLabelValue}
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
	newPod := pod.DefinePod(name, params.QeTestNamespace, globalhelper.Configuration.General.TestImage)
	// Add labels.
	newPod = pod.RedefinePodWithLabel(newPod, params.TnfTargetPodLabels)
	// Change command to use the stdout buffer.
	newPod.Spec.Containers[0].Command = getContainerCommandWithStdout(stdoutBuffer)

	return newPod
}

func DefineDeploymentWithoutTargetLabels(name string) *appsv1.Deployment {
	return deployment.DefineDeployment(name, params.QeTestNamespace,
		globalhelper.Configuration.General.TestImage,
		map[string]string{"fakeLabelKey": "fakeLabelValue"})
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
				Name:    fmt.Sprintf("%s-%d", params.QeTestContainerBaseName, index),
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
	dep := deployment.DefineDeployment(name, params.QeTestNamespace,
		globalhelper.Configuration.General.TestImage, params.TnfTargetPodLabels)

	// Customize its replicas and container specs.
	dep = deployment.RedefineWithReplicaNumber(dep, int32(replicas))
	dep = deployment.RedefineWithContainerSpecs(dep, containerSpecs)

	return dep
}

func defineStatefulSetWithContainerSpecs(name string, replicas int,
	containerSpecs []corev1.Container) *appsv1.StatefulSet {
	// Define base statefulSet
	sts := statefulset.DefineStatefulSet(name, params.QeTestNamespace,
		globalhelper.Configuration.General.TestImage, params.TnfTargetPodLabels)

	// Customize its replicas and container specs.
	sts = statefulset.RedefineWithReplicaNumber(sts, int32(replicas))
	sts = statefulset.RedefineWithContainerSpecs(sts, containerSpecs)

	return sts
}

func defineDaemonSetWithContainerSpecs(name string,
	containerSpecs []corev1.Container) *appsv1.DaemonSet {
	// Define base daemonSet
	daemonSet := daemonset.DefineDaemonSet(params.QeTestNamespace,
		globalhelper.Configuration.General.TestImage, params.TnfTargetPodLabels, name)

	// Customize its container specs.
	return daemonset.RedefineWithContainerSpecs(daemonSet, containerSpecs)
}
