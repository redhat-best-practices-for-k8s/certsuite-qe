package helper

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/onsi/gomega"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	apiextv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	apiextclient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset/typed/apiextensions/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	k8sfake "k8s.io/client-go/kubernetes/fake"

	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/globalhelper"
	tsparams "github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/observability/parameters"
)

var errInjected = errors.New("injected failure")

var crdResource = schema.GroupResource{Group: "apiextensions.k8s.io", Resource: "customresourcedefinitions"}

// The vendored apiextensions clientset has no fake implementation, so the CRD
// client is faked by hand, mirroring the OLM fakes used by the other suites.
type fakeCrds struct {
	apiextclient.CustomResourceDefinitionInterface

	crds      map[string]*apiextv1.CustomResourceDefinition
	createErr error

	// establishedAfter keeps Get returning the CRD without the Established
	// condition for this many calls first; 0 reports it immediately.
	establishedAfter int

	// getFailsOnce makes only the first Get call fail, to exercise the retry loop.
	getFailsOnce bool
	getCalls     int
}

func (f *fakeCrds) Create(_ context.Context, crd *apiextv1.CustomResourceDefinition,
	_ metav1.CreateOptions) (*apiextv1.CustomResourceDefinition, error) {
	if f.createErr != nil {
		return nil, f.createErr
	}

	if _, exists := f.crds[crd.Name]; exists {
		return nil, k8serrors.NewAlreadyExists(crdResource, crd.Name)
	}

	f.crds[crd.Name] = crd

	return crd, nil
}

func (f *fakeCrds) Get(_ context.Context, name string,
	_ metav1.GetOptions) (*apiextv1.CustomResourceDefinition, error) {
	f.getCalls++

	if f.getFailsOnce && f.getCalls == 1 {
		return nil, errInjected
	}

	crd, exists := f.crds[name]
	if !exists {
		return nil, k8serrors.NewNotFound(crdResource, name)
	}

	if f.getCalls > f.establishedAfter {
		establishedCrd := crd.DeepCopy()
		establishedCrd.Status.Conditions = []apiextv1.CustomResourceDefinitionCondition{{
			Type: apiextv1.Established,
		}}

		return establishedCrd, nil
	}

	return crd, nil
}

func (f *fakeCrds) Delete(_ context.Context, name string, _ metav1.DeleteOptions) error {
	delete(f.crds, name)

	return nil
}

type fakeApiextensions struct {
	apiextclient.ApiextensionsV1Interface
	crds *fakeCrds
}

//nolint:ireturn // Must match the apiextensions typed client interface.
func (f *fakeApiextensions) CustomResourceDefinitions() apiextclient.CustomResourceDefinitionInterface {
	return f.crds
}

// withFakeCrds installs a fake CRD client, restoring the global client on cleanup.
func withFakeCrds(t *testing.T, crds ...*apiextv1.CustomResourceDefinition) *fakeCrds {
	t.Helper()

	globalhelper.SetTestK8sAPIClient(k8sfake.NewClientset())

	crdClient := &fakeCrds{crds: map[string]*apiextv1.CustomResourceDefinition{}}
	for _, aCrd := range crds {
		crdClient.crds[aCrd.Name] = aCrd
	}

	globalhelper.GetAPIClient().ApiextensionsV1Interface = &fakeApiextensions{crds: crdClient}

	t.Cleanup(globalhelper.UnsetTestK8sAPIClient)

	return crdClient
}

func assertContainerIdentity(t *testing.T, container corev1.Container, index int) {
	t.Helper()

	assert.Equal(t, fmt.Sprintf("%s-%d", tsparams.TestContainerBaseName, index), container.Name)
	assert.Equal(t, tsparams.SampleWorkloadImage, container.Image)
}

// assertStdoutContainers checks the container specs derived from stdout buffers;
// expected commands are the space-joined getContainerCommandWithStdout output.
func assertStdoutContainers(t *testing.T, containers []corev1.Container, expectedCommands []string) {
	t.Helper()

	require.Len(t, containers, len(expectedCommands))

	for index, expectedCommand := range expectedCommands {
		assertContainerIdentity(t, containers[index], index)
		assert.Equal(t, expectedCommand, strings.Join(containers[index].Command, " "))
	}
}

func assertTerminationMsgContainers(t *testing.T, containers []corev1.Container,
	expectedPolicies []corev1.TerminationMessagePolicy) {
	t.Helper()

	require.Len(t, containers, len(expectedPolicies))

	for index, expectedPolicy := range expectedPolicies {
		assertContainerIdentity(t, containers[index], index)
		assert.Equal(t, tsparams.TestContainerNormalCommand, containers[index].Command)
		assert.Equal(t, expectedPolicy, containers[index].TerminationMessagePolicy)
	}
}

func TestGetCertsuiteTargetPodLabelsSlice(t *testing.T) {
	assert.Equal(t, []string{tsparams.TestPodLabelKey + ":" + tsparams.TestPodLabelValue},
		GetCertsuiteTargetPodLabelsSlice())
}

func TestDefineDeploymentWithStdoutBuffers(t *testing.T) {
	dep := DefineDeploymentWithStdoutBuffers("obs-dp", "obs-ns", 3,
		[]string{"Hello 1\n", "Hello 2"})

	assert.Equal(t, "obs-dp", dep.Name)
	assert.Equal(t, "obs-ns", dep.Namespace)
	assert.Equal(t, tsparams.CertsuiteTargetPodLabels, dep.Labels)
	assert.Equal(t, tsparams.CertsuiteTargetPodLabels, dep.Spec.Selector.MatchLabels)
	assert.Equal(t, tsparams.CertsuiteTargetPodLabels, dep.Spec.Template.Labels)
	require.NotNil(t, dep.Spec.Replicas)
	assert.Equal(t, int32(3), *dep.Spec.Replicas)
	assertStdoutContainers(t, dep.Spec.Template.Spec.Containers, []string{
		`/bin/bash -c printf "Hello 1\n" && sleep INF`,
		`/bin/bash -c printf "Hello 2" && sleep INF`,
	})
}

func TestDefineStatefulSetWithStdoutBuffers(t *testing.T) {
	sts := DefineStatefulSetWithStdoutBuffers("obs-st", "obs-ns", 2,
		[]string{"Hello 1"})

	assert.Equal(t, "obs-st", sts.Name)
	assert.Equal(t, "obs-ns", sts.Namespace)
	assert.Equal(t, tsparams.CertsuiteTargetPodLabels, sts.Spec.Selector.MatchLabels)
	require.NotNil(t, sts.Spec.Replicas)
	assert.Equal(t, int32(2), *sts.Spec.Replicas)
	assertStdoutContainers(t, sts.Spec.Template.Spec.Containers, []string{
		`/bin/bash -c printf "Hello 1" && sleep INF`,
	})
}

func TestDefineDaemonSetWithStdoutBuffers(t *testing.T) {
	daemonSet := DefineDaemonSetWithStdoutBuffers("obs-ds", "obs-ns",
		[]string{"Hello 1\n", "Hello 2"})

	assert.Equal(t, "obs-ds", daemonSet.Name)
	assert.Equal(t, "obs-ns", daemonSet.Namespace)
	assert.Equal(t, tsparams.CertsuiteTargetPodLabels, daemonSet.Spec.Selector.MatchLabels)
	assertStdoutContainers(t, daemonSet.Spec.Template.Spec.Containers, []string{
		`/bin/bash -c printf "Hello 1\n" && sleep INF`,
		`/bin/bash -c printf "Hello 2" && sleep INF`,
	})
}

func TestDefinePodWithStdoutBuffer(t *testing.T) {
	pod := DefinePodWithStdoutBuffer("obs-pod", "obs-ns", "Hello pod")

	assert.Equal(t, "obs-pod", pod.Name)
	assert.Equal(t, "obs-ns", pod.Namespace)
	assert.Equal(t, tsparams.CertsuiteTargetPodLabels, pod.Labels)
	require.Len(t, pod.Spec.Containers, 1)
	assert.Equal(t, tsparams.SampleWorkloadImage, pod.Spec.Containers[0].Image)
	assert.Equal(t, `/bin/bash -c printf "Hello pod" && sleep INF`,
		strings.Join(pod.Spec.Containers[0].Command, " "))
}

func TestDefineDeploymentWithoutTargetLabels(t *testing.T) {
	fakeLabels := map[string]string{"fakeLabelKey": "fakeLabelValue"}

	dep := DefineDeploymentWithoutTargetLabels("obs-dp", "obs-ns")

	assert.Equal(t, "obs-dp", dep.Name)
	assert.Equal(t, "obs-ns", dep.Namespace)
	assert.Equal(t, fakeLabels, dep.Labels)
	assert.Equal(t, fakeLabels, dep.Spec.Selector.MatchLabels)
	require.Len(t, dep.Spec.Template.Spec.Containers, 1)
	assert.Equal(t, tsparams.SampleWorkloadImage, dep.Spec.Template.Spec.Containers[0].Image)
}

func TestDefineCrdWithStatusSubresource(t *testing.T) {
	crd := DefineCrdWithStatusSubresource("Memcached", tsparams.CrdSuffix1)

	require.Len(t, crd.Spec.Versions, 1)
	version := crd.Spec.Versions[0]

	assert.Equal(t, "memcacheds."+tsparams.CrdSuffix1, crd.Name)
	assert.Equal(t, tsparams.CrdSuffix1, crd.Spec.Group)
	assert.Equal(t, "Memcached", crd.Spec.Names.Kind)
	assert.Equal(t, "memcached", crd.Spec.Names.Singular)
	assert.Equal(t, "memcacheds", crd.Spec.Names.Plural)
	assert.Equal(t, "MemcachedList", crd.Spec.Names.ListKind)
	assert.True(t, version.Served)
	assert.True(t, version.Storage)
	require.NotNil(t, version.Subresources)
	require.NotNil(t, version.Subresources.Status)
	require.NotNil(t, version.Schema)
	require.NotNil(t, version.Schema.OpenAPIV3Schema)
	assert.Contains(t, version.Schema.OpenAPIV3Schema.Properties, "spec")
	assert.Contains(t, version.Schema.OpenAPIV3Schema.Properties, "status")
}

func TestDefineCrdWithoutStatusSubresource(t *testing.T) {
	crd := DefineCrdWithoutStatusSubresource("Memcached", tsparams.CrdSuffix2)

	require.Len(t, crd.Spec.Versions, 1)
	version := crd.Spec.Versions[0]

	assert.Equal(t, "memcacheds."+tsparams.CrdSuffix2, crd.Name)
	assert.Equal(t, tsparams.CrdSuffix2, crd.Spec.Group)
	assert.Nil(t, version.Subresources)
	require.NotNil(t, version.Schema)
	require.NotNil(t, version.Schema.OpenAPIV3Schema)
	assert.Contains(t, version.Schema.OpenAPIV3Schema.Properties, "spec")
	assert.NotContains(t, version.Schema.OpenAPIV3Schema.Properties, "status")
}

func TestDefineDeploymentWithTerminationMsgPolicies(t *testing.T) {
	policies := []corev1.TerminationMessagePolicy{
		corev1.TerminationMessageReadFile,
		corev1.TerminationMessageFallbackToLogsOnError,
		tsparams.UseDefaultTerminationMsgPolicy,
	}

	dep := DefineDeploymentWithTerminationMsgPolicies("obs-dp", "obs-ns", 2, policies)

	assert.Equal(t, "obs-dp", dep.Name)
	require.NotNil(t, dep.Spec.Replicas)
	assert.Equal(t, int32(2), *dep.Spec.Replicas)
	assertTerminationMsgContainers(t, dep.Spec.Template.Spec.Containers, policies)
}

func TestDefineDaemonSetWithTerminationMsgPolicies(t *testing.T) {
	policies := []corev1.TerminationMessagePolicy{corev1.TerminationMessageReadFile}

	daemonSet := DefineDaemonSetWithTerminationMsgPolicies("obs-ds", "obs-ns", policies)

	assert.Equal(t, "obs-ds", daemonSet.Name)
	assertTerminationMsgContainers(t, daemonSet.Spec.Template.Spec.Containers, policies)
}

func TestDefineStatefulSetWithTerminationMsgPolicies(t *testing.T) {
	policies := []corev1.TerminationMessagePolicy{
		corev1.TerminationMessageFallbackToLogsOnError,
		tsparams.UseDefaultTerminationMsgPolicy,
	}

	sts := DefineStatefulSetWithTerminationMsgPolicies("obs-st", "obs-ns", 1, policies)

	assert.Equal(t, "obs-st", sts.Name)
	require.NotNil(t, sts.Spec.Replicas)
	assert.Equal(t, int32(1), *sts.Spec.Replicas)
	assertTerminationMsgContainers(t, sts.Spec.Template.Spec.Containers, policies)
}

func TestCreateAndWaitUntilCrdIsReady(t *testing.T) {
	testCrd := DefineCrdWithStatusSubresource("Memcached", tsparams.CrdSuffix1)

	t.Run("waits until the crd reports the established condition", func(t *testing.T) {
		// Route gomega assertion failures from the code under test to the failing subtest.
		gomega.RegisterTestingT(t)

		crdClient := withFakeCrds(t)

		require.NoError(t, CreateAndWaitUntilCrdIsReady(testCrd, 3*time.Second))
		assert.Contains(t, crdClient.crds, testCrd.Name)
	})

	t.Run("tolerates a crd that already exists", func(t *testing.T) {
		gomega.RegisterTestingT(t)

		withFakeCrds(t, testCrd)

		require.NoError(t, CreateAndWaitUntilCrdIsReady(testCrd, 3*time.Second))
	})

	t.Run("wraps the create failure", func(t *testing.T) {
		gomega.RegisterTestingT(t)

		crdClient := withFakeCrds(t)
		crdClient.createErr = errInjected

		err := CreateAndWaitUntilCrdIsReady(testCrd, 3*time.Second)
		require.ErrorIs(t, err, errInjected)
		assert.Contains(t, err.Error(), "failed to create crd")
	})

	t.Run("retries when reading the crd fails", func(t *testing.T) {
		gomega.RegisterTestingT(t)

		crdClient := withFakeCrds(t)
		crdClient.getFailsOnce = true

		// The poll interval is tsparams.CrdRetryInterval, so the timeout must span one retry.
		require.NoError(t, CreateAndWaitUntilCrdIsReady(testCrd, 3*tsparams.CrdRetryInterval))
	})

	t.Run("retries until the crd reports the established condition", func(t *testing.T) {
		gomega.RegisterTestingT(t)

		crdClient := withFakeCrds(t)
		crdClient.establishedAfter = 1

		// The poll interval is tsparams.CrdRetryInterval, so the timeout must span one retry.
		require.NoError(t, CreateAndWaitUntilCrdIsReady(testCrd, 3*tsparams.CrdRetryInterval))
	})
}

func TestDeleteCrdAndWaitUntilIsRemoved(t *testing.T) {
	testCrd := DefineCrdWithStatusSubresource("Memcached", tsparams.CrdSuffix2)

	t.Run("deletes and waits until the crd is gone", func(t *testing.T) {
		gomega.RegisterTestingT(t)

		crdClient := withFakeCrds(t, testCrd)

		DeleteCrdAndWaitUntilIsRemoved(testCrd.Name, 3*time.Second)

		assert.NotContains(t, crdClient.crds, testCrd.Name)
	})
}
