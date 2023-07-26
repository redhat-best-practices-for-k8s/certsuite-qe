package helper

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/golang/glog"
	. "github.com/onsi/gomega"

	"github.com/test-network-function/cnfcert-tests-verification/tests/globalhelper"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/crd"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/daemonset"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/deployment"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/pod"
	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/statefulset"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	apiextv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	tsparams "github.com/test-network-function/cnfcert-tests-verification/tests/observability/parameters"
)

// For some reason, there's a function that expects labels' key/values separated
// by colon instead of the equal char.
func GetTnfTargetPodLabelsSlice() []string {
	return []string{tsparams.TestPodLabelKey + ":" + tsparams.TestPodLabelValue}
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
	newPod := pod.DefinePod(name, tsparams.TestNamespace, globalhelper.GetConfiguration().General.TestImage, tsparams.TnfTargetPodLabels)

	// Change command to use the stdout buffer.
	newPod.Spec.Containers[0].Command = getContainerCommandWithStdout(stdoutBuffer)

	return newPod
}

func DefineDeploymentWithoutTargetLabels(name string) *appsv1.Deployment {
	return deployment.DefineDeployment(name, tsparams.TestNamespace,
		globalhelper.GetConfiguration().General.TestImage,
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

func CreateAndWaitUntilCrdIsReady(crd *apiextv1.CustomResourceDefinition, timeout time.Duration) error {
	_, err := globalhelper.GetAPIClient().CustomResourceDefinitions().Create(
		context.TODO(),
		crd,
		metav1.CreateOptions{},
	)
	if err != nil {
		return fmt.Errorf("failed to create crd: %w", err)
	}

	Eventually(func() bool {
		runningCrd, err := globalhelper.GetAPIClient().CustomResourceDefinitions().Get(
			context.TODO(),
			crd.Name,
			metav1.GetOptions{},
		)
		if err != nil {
			glog.V(5).Info(fmt.Sprintf(
				"crd %s is not ready, retry in 5 seconds", runningCrd.Name))

			return false
		}

		for _, condition := range runningCrd.Status.Conditions {
			if condition.Type == apiextv1.Established {
				return true
			}
		}

		return false
	}, timeout, tsparams.CrdRetryInterval).Should(Equal(true), "CRD is not ready")

	return nil
}

func DeleteCrdAndWaitUntilIsRemoved(crd string, timeout time.Duration) {
	err := globalhelper.GetAPIClient().CustomResourceDefinitions().Delete(
		context.TODO(),
		crd,
		metav1.DeleteOptions{})
	Expect(err).ToNot(HaveOccurred())

	Eventually(func() bool {
		_, err := globalhelper.GetAPIClient().CustomResourceDefinitions().Get(
			context.TODO(),
			crd,
			metav1.GetOptions{})

		// If the CRD was already removed, we'll get an error.
		return err != nil
	}, timeout, tsparams.CrdRetryInterval).Should(Equal(true), "CRD is not removed yet")
}

func DefineDeploymentWithTerminationMsgPolicies(name string, replicas int,
	policies []corev1.TerminationMessagePolicy) *appsv1.Deployment {
	// Create one container spec per policy.
	containerSpecs := createContainerSpecsFromTerminationMsgPolicies(policies)

	return defineDeploymentWithContainerSpecs(name, replicas, containerSpecs)
}

func DefineDaemonSetWithTerminationMsgPolicies(name string,
	policies []corev1.TerminationMessagePolicy) *appsv1.DaemonSet {
	// Create one container spec per policy.
	containerSpecs := createContainerSpecsFromTerminationMsgPolicies(policies)

	return defineDaemonSetWithContainerSpecs(name, containerSpecs)
}

func DefineStatefulSetWithTerminationMsgPolicies(name string, replicas int,
	policies []corev1.TerminationMessagePolicy) *appsv1.StatefulSet {
	// Create one container spec per policy.
	containerSpecs := createContainerSpecsFromTerminationMsgPolicies(policies)

	return defineStatefulSetWithContainerSpecs(name, replicas, containerSpecs)
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
				Name:    fmt.Sprintf("%s-%d", tsparams.TestContainerBaseName, index),
				Image:   globalhelper.GetConfiguration().General.TestImage,
				Command: getContainerCommandWithStdout(stdoutLines),
			},
		)
	}

	return containerSpecs
}

func createContainerSpecsFromTerminationMsgPolicies(policies []corev1.TerminationMessagePolicy) []corev1.Container {
	numContainers := len(policies)
	containerSpecs := []corev1.Container{}

	for index := 0; index < numContainers; index++ {
		container := corev1.Container{
			Name:    fmt.Sprintf("%s-%d", tsparams.TestContainerBaseName, index),
			Image:   globalhelper.GetConfiguration().General.TestImage,
			Command: tsparams.TestContainerNormalCommand,
		}

		// If the policy is an empty string, leave the field unset, and k8s will
		// set it as TerminationMessageReadFile
		if policies[index] != tsparams.UseDefaultTerminationMsgPolicy {
			container.TerminationMessagePolicy = policies[index]
		}
		containerSpecs = append(containerSpecs, container)
	}

	return containerSpecs
}

func defineDeploymentWithContainerSpecs(name string, replicas int,
	containerSpecs []corev1.Container) *appsv1.Deployment {
	// Define base deployment
	dep := deployment.DefineDeployment(name, tsparams.TestNamespace,
		globalhelper.GetConfiguration().General.TestImage, tsparams.TnfTargetPodLabels)

	// Customize its replicas and container specs.
	deployment.RedefineWithReplicaNumber(dep, int32(replicas))
	deployment.RedefineWithContainerSpecs(dep, containerSpecs)

	return dep
}

func defineStatefulSetWithContainerSpecs(name string, replicas int,
	containerSpecs []corev1.Container) *appsv1.StatefulSet {
	// Define base statefulSet
	sts := statefulset.DefineStatefulSet(name, tsparams.TestNamespace,
		globalhelper.GetConfiguration().General.TestImage, tsparams.TnfTargetPodLabels)

	// Customize its replicas and container specs.
	statefulset.RedefineWithReplicaNumber(sts, int32(replicas))
	statefulset.RedefineWithContainerSpecs(sts, containerSpecs)

	return sts
}

func defineDaemonSetWithContainerSpecs(name string,
	containerSpecs []corev1.Container) *appsv1.DaemonSet {
	// Define base daemonSet
	daemonSet := daemonset.DefineDaemonSet(tsparams.TestNamespace,
		globalhelper.GetConfiguration().General.TestImage, tsparams.TnfTargetPodLabels, name)

	// Customize its container specs.
	daemonset.RedefineWithContainerSpecs(daemonSet, containerSpecs)

	return daemonSet
}
