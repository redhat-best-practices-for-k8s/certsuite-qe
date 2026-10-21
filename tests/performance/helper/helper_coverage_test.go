package helper

import (
	"bytes"
	"context"
	"errors"
	"testing"

	egiClients "github.com/openshift-kni/eco-goinfra/pkg/clients"
	egipod "github.com/openshift-kni/eco-goinfra/pkg/pod"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	nodev1 "k8s.io/api/node/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	k8sfake "k8s.io/client-go/kubernetes/fake"
	k8stesting "k8s.io/client-go/testing"

	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/globalhelper"
	tsparams "github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/performance/parameters"
)

var errInjected = errors.New("injected failure")

// withStubbedEgiPod replaces the eco-goinfra seams with no-op fakes, restoring them on cleanup.
// Individual tests override the seam they need. Tests using it must not run in parallel.
func withStubbedEgiPod(t *testing.T) {
	t.Helper()

	oldClient, oldPull, oldExec := getEgiClient, egiPodPull, egiPodExec

	t.Cleanup(func() { getEgiClient, egiPodPull, egiPodExec = oldClient, oldPull, oldExec })

	getEgiClient = func() *egiClients.Settings { return nil }
	egiPodPull = func(*egiClients.Settings, string, string) (*egipod.Builder, error) {
		return &egipod.Builder{}, nil
	}
	egiPodExec = func(*egipod.Builder, []string, ...string) (bytes.Buffer, error) {
		return bytes.Buffer{}, nil
	}
}

// withFakeClient installs a fake k8s clientset, restoring the global client on cleanup.
func withFakeClient(t *testing.T, objects ...runtime.Object) *k8sfake.Clientset {
	t.Helper()

	client := k8sfake.NewClientset(objects...)
	globalhelper.SetTestK8sAPIClient(client)
	t.Cleanup(globalhelper.UnsetTestK8sAPIClient)

	return client
}

// failCreate makes the fake clientset fail create calls for the given resource.
func failCreate(client *k8sfake.Clientset, resource string) {
	client.PrependReactor("create", resource, func(k8stesting.Action) (bool, runtime.Object, error) {
		return true, nil, errInjected
	})
}

// failPodPull overrides the pull seam to fail; withStubbedEgiPod restores it on cleanup.
func failPodPull() {
	egiPodPull = func(*egiClients.Settings, string, string) (*egipod.Builder, error) {
		return nil, errInjected
	}
}

func testExecPod() *corev1.Pod {
	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{Name: "perf-pod", Namespace: "perf-ns"},
		Spec:       corev1.PodSpec{Containers: []corev1.Container{{Name: "rt-app"}}},
	}
}

func TestExecCommandContainer(t *testing.T) {
	pod := testExecPod()

	t.Run("returns stdout from the exec", func(t *testing.T) {
		withStubbedEgiPod(t)

		var pulledName, pulledNamespace string

		var executedCommand []string

		egiPodPull = func(_ *egiClients.Settings, name, namespace string) (*egipod.Builder, error) {
			pulledName, pulledNamespace = name, namespace

			return &egipod.Builder{}, nil
		}
		egiPodExec = func(_ *egipod.Builder, command []string, _ ...string) (bytes.Buffer, error) {
			executedCommand = command

			return *bytes.NewBufferString("policy: SCHED_FIFO"), nil
		}

		stdout, stderr, err := ExecCommandContainer(pod, "chrt -p 1")
		require.NoError(t, err)
		assert.Equal(t, "policy: SCHED_FIFO", stdout)
		assert.Empty(t, stderr)
		assert.Equal(t, "perf-pod", pulledName)
		assert.Equal(t, "perf-ns", pulledNamespace)
		assert.Equal(t, []string{"sh", "-c", "chrt -p 1"}, executedCommand)
	})

	t.Run("pull failure", func(t *testing.T) {
		withStubbedEgiPod(t)
		failPodPull()

		_, _, err := ExecCommandContainer(pod, "chrt -p 1")
		require.ErrorIs(t, err, errInjected)
	})

	t.Run("exec failure returns partial stdout", func(t *testing.T) {
		withStubbedEgiPod(t)
		egiPodExec = func(*egipod.Builder, []string, ...string) (bytes.Buffer, error) {
			return *bytes.NewBufferString("partial"), errInjected
		}

		stdout, _, err := ExecCommandContainer(pod, "chrt -p 1")
		require.ErrorIs(t, err, errInjected)
		assert.Equal(t, "partial", stdout)
	})

	t.Run("pod without containers", func(t *testing.T) {
		withStubbedEgiPod(t)

		pod := &corev1.Pod{ObjectMeta: metav1.ObjectMeta{Name: "empty-pod", Namespace: "perf-ns"}}

		_, _, err := ExecCommandContainer(pod, "chrt -p 1")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "pod empty-pod has no containers")
	})
}

func TestChangeSchedulingPolicy(t *testing.T) {
	pod := testExecPod()

	t.Run("passes through the exec output", func(t *testing.T) {
		withStubbedEgiPod(t)
		egiPodExec = func(*egipod.Builder, []string, ...string) (bytes.Buffer, error) {
			return *bytes.NewBufferString("fifo"), nil
		}

		outStr, errStr, err := ChangeSchedulingPolicy(pod, "chrt -f 0 1")
		require.NoError(t, err)
		assert.Equal(t, "fifo", outStr)
		assert.Empty(t, errStr)
	})

	t.Run("wraps the exec failure", func(t *testing.T) {
		withStubbedEgiPod(t)
		failPodPull()

		_, _, err := ChangeSchedulingPolicy(pod, "chrt -f 0 1")
		require.ErrorIs(t, err, errInjected)
		assert.Contains(t, err.Error(), "cannot execute command")
		assert.Contains(t, err.Error(), "on rt-app")
	})

	t.Run("pod without containers", func(t *testing.T) {
		withStubbedEgiPod(t)

		pod := &corev1.Pod{ObjectMeta: metav1.ObjectMeta{Name: "empty-pod", Namespace: "perf-ns"}}

		_, _, err := ChangeSchedulingPolicy(pod, "chrt -f 0 1")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "cannot execute command")
		assert.Contains(t, err.Error(), "empty-pod")
	})
}

func TestConfigurePrivilegedServiceAccount(t *testing.T) {
	namespace := "perf-ns"

	t.Run("creates the role, rolebinding and service account", func(t *testing.T) {
		client := withFakeClient(t)

		require.NoError(t, ConfigurePrivilegedServiceAccount(namespace))

		role, err := client.RbacV1().Roles(namespace).Get(
			context.TODO(), tsparams.PrivilegedRoleName, metav1.GetOptions{})
		require.NoError(t, err)
		require.Len(t, role.Rules, 1)
		assert.Equal(t, rbacv1.PolicyRule{
			APIGroups: []string{"*"}, Resources: []string{"*"}, Verbs: []string{"*"},
		}, role.Rules[0])

		binding, err := client.RbacV1().RoleBindings(namespace).Get(
			context.TODO(), tsparams.PrivilegedRoleName, metav1.GetOptions{})
		require.NoError(t, err)
		assert.Equal(t, tsparams.PrivilegedRoleName, binding.RoleRef.Name)
		require.Len(t, binding.Subjects, 1)
		assert.Equal(t, tsparams.PrivilegedRoleName, binding.Subjects[0].Name)
		assert.Equal(t, namespace, binding.Subjects[0].Namespace)

		serviceAccount, err := client.CoreV1().ServiceAccounts(namespace).Get(
			context.TODO(), tsparams.PrivilegedRoleName, metav1.GetOptions{})
		require.NoError(t, err)
		assert.Equal(t, tsparams.PrivilegedRoleName, serviceAccount.Name)
	})

	t.Run("tolerates objects that already exist", func(t *testing.T) {
		role, roleBinding, serviceAccount := getPrivilegedServiceAccountObjects(namespace)
		// Distinctive content proves the call does not overwrite what is already there.
		role.Rules = []rbacv1.PolicyRule{{Verbs: []string{"keep"}}}
		roleBinding.Subjects = []rbacv1.Subject{{Kind: "User", Name: "keep"}}

		client := withFakeClient(t, &role, &roleBinding, &serviceAccount)

		require.NoError(t, ConfigurePrivilegedServiceAccount(namespace))

		gotRole, err := client.RbacV1().Roles(namespace).Get(
			context.TODO(), tsparams.PrivilegedRoleName, metav1.GetOptions{})
		require.NoError(t, err)
		assert.Equal(t, role.Rules, gotRole.Rules)

		gotBinding, err := client.RbacV1().RoleBindings(namespace).Get(
			context.TODO(), tsparams.PrivilegedRoleName, metav1.GetOptions{})
		require.NoError(t, err)
		assert.Equal(t, roleBinding.Subjects, gotBinding.Subjects)

		gotServiceAccount, err := client.CoreV1().ServiceAccounts(namespace).Get(
			context.TODO(), tsparams.PrivilegedRoleName, metav1.GetOptions{})
		require.NoError(t, err)
		assert.Equal(t, serviceAccount.Name, gotServiceAccount.Name)
	})

	createFailures := []struct {
		resource string
		wantMsg  string
	}{
		{resource: "roles", wantMsg: "error creating role"},
		{resource: "rolebindings", wantMsg: "error creating rolebinding"},
		{resource: "serviceaccounts", wantMsg: "error creating service account"},
	}

	for _, failure := range createFailures {
		t.Run(failure.resource+" create failure", func(t *testing.T) {
			client := withFakeClient(t)
			failCreate(client, failure.resource)

			err := ConfigurePrivilegedServiceAccount(namespace)
			require.ErrorIs(t, err, errInjected)
			assert.Contains(t, err.Error(), failure.wantMsg)
		})
	}
}

func TestDefineExclusivePod(t *testing.T) {
	labels := map[string]string{"app": "test"}

	pod := DefineExclusivePod("exclusive-pod", "perf-ns", "quay.io/test/image:1", labels)

	assert.Equal(t, "exclusive-pod", pod.Name)
	assert.Equal(t, "perf-ns", pod.Namespace)
	assert.Equal(t, labels, pod.Labels)
	assert.Equal(t, tsparams.PrivilegedRoleName, pod.Spec.ServiceAccountName)
	require.NotNil(t, pod.Spec.TerminationGracePeriodSeconds)
	assert.Equal(t, int64(0), *pod.Spec.TerminationGracePeriodSeconds)
	require.Len(t, pod.Spec.Containers, 2)
	assert.Equal(t, "shared", pod.Spec.Containers[0].Name)
	assert.Equal(t, "exclusive", pod.Spec.Containers[1].Name)

	for _, container := range pod.Spec.Containers {
		assert.Equal(t, "quay.io/test/image:1", container.Image)
		assert.Equal(t, []string{"/bin/bash", "-c", "sleep INF"}, container.Command)
		assert.Equal(t, resource.MustParse("1"), container.Resources.Limits[corev1.ResourceCPU])
		assert.Equal(t, resource.MustParse("112Mi"), container.Resources.Limits[corev1.ResourceMemory])
		assert.Equal(t, resource.MustParse("1"), container.Resources.Requests[corev1.ResourceCPU])
		assert.Equal(t, resource.MustParse("112Mi"), container.Resources.Requests[corev1.ResourceMemory])
	}

	require.NotNil(t, pod.Spec.SecurityContext)
	assert.Equal(t, int64(1000), *pod.Spec.SecurityContext.RunAsUser)
	assert.Equal(t, int64(1000), *pod.Spec.SecurityContext.RunAsGroup)
	assert.True(t, *pod.Spec.SecurityContext.RunAsNonRoot)
}

func TestDefineRtPod(t *testing.T) {
	labels := map[string]string{"app": "test"}

	pod := DefineRtPod("rt-pod", "perf-ns", "quay.io/test/image:1", labels)

	assert.Equal(t, "rt-pod", pod.Name)
	assert.Equal(t, "perf-ns", pod.Namespace)
	assert.Equal(t, labels, pod.Labels)
	assert.Equal(t, tsparams.PrivilegedRoleName, pod.Spec.ServiceAccountName)
	require.Len(t, pod.Spec.Containers, 1)

	container := pod.Spec.Containers[0]
	assert.Equal(t, "rt-app", container.Name)
	assert.Equal(t, "quay.io/test/image:1", container.Image)
	assert.Equal(t, []string{"/bin/bash", "-c", "sleep INF"}, container.Command)
	assert.Equal(t, resource.MustParse("1"), container.Resources.Limits[corev1.ResourceCPU])
	assert.Equal(t, resource.MustParse("112Mi"), container.Resources.Limits[corev1.ResourceMemory])
	require.NotNil(t, container.SecurityContext)
	assert.True(t, *container.SecurityContext.Privileged)
	assert.Equal(t, int64(0), *container.SecurityContext.RunAsUser)
}

func TestRedefinePodWithSharedContainer(t *testing.T) {
	pod := DefineExclusivePod("exclusive-pod", "perf-ns", "quay.io/test/image:1", nil)

	RedefinePodWithSharedContainer(pod, 0)
	assert.Equal(t, resource.MustParse("1"), pod.Spec.Containers[0].Resources.Limits[corev1.ResourceCPU])
	assert.Equal(t, resource.MustParse("150m"), pod.Spec.Containers[0].Resources.Requests[corev1.ResourceCPU])
	assert.Equal(t, resource.MustParse("1"), pod.Spec.Containers[1].Resources.Requests[corev1.ResourceCPU])

	RedefinePodWithSharedContainer(pod, 1)
	assert.Equal(t, resource.MustParse("150m"), pod.Spec.Containers[1].Resources.Requests[corev1.ResourceCPU])

	// Out-of-range indices leave the pod untouched.
	pod = DefineExclusivePod("exclusive-pod", "perf-ns", "quay.io/test/image:1", nil)
	RedefinePodWithSharedContainer(pod, -1)
	RedefinePodWithSharedContainer(pod, 2)
	assert.Equal(t, resource.MustParse("1"), pod.Spec.Containers[0].Resources.Requests[corev1.ResourceCPU])
}

func TestDefineRtPodInIsolatedCPUPool(t *testing.T) {
	rtc := &nodev1.RuntimeClass{ObjectMeta: metav1.ObjectMeta{Name: tsparams.CertsuiteRunTimeClass}}

	testPod, err := DefineRtPodInIsolatedCPUPool("perf-ns", rtc)
	require.NoError(t, err)
	assert.Equal(t, tsparams.TestPodName, testPod.Name)
	assert.Equal(t, tsparams.CertsuiteTargetPodLabels, testPod.Labels)

	require.NotNil(t, testPod.Spec.RuntimeClassName)
	assert.Equal(t, tsparams.CertsuiteRunTimeClass, *testPod.Spec.RuntimeClassName)
	assert.Equal(t, tsparams.DisableStr, testPod.Annotations["cpu-load-balancing.crio.io"])
	assert.Equal(t, tsparams.DisableStr, testPod.Annotations["irq-load-balancing.crio.io"])

	require.Len(t, testPod.Spec.Containers, 1)
	assert.Equal(t, resource.MustParse("1"), testPod.Spec.Containers[0].Resources.Limits[corev1.ResourceCPU])
	assert.Equal(t, resource.MustParse("1"), testPod.Spec.Containers[0].Resources.Requests[corev1.ResourceCPU])
}

func TestDefineDpdkPod(t *testing.T) {
	pod := DefineDpdkPod("dpdk-pod", "perf-ns")

	assert.Equal(t, "dpdk-pod", pod.Name)
	assert.Equal(t, "perf-ns", pod.Namespace)
	assert.Equal(t, tsparams.CertsuiteTargetPodLabels, pod.Labels)
	assert.Equal(t, "sriovnet1", pod.Annotations["k8s.v1.cni.cncf.io/networks"])

	require.Len(t, pod.Spec.Containers, 1)

	container := pod.Spec.Containers[0]
	assert.Equal(t, "app-container", container.Name)
	assert.Equal(t, "registry.redhat.io/openshift4/dpdk-base-rhel8:v4.9", container.Image)
	assert.Equal(t, resource.MustParse("512Mi"), container.Resources.Limits[corev1.ResourceMemory])
	require.NotNil(t, container.SecurityContext)
	assert.True(t, *container.SecurityContext.Privileged)
	require.NotNil(t, container.SecurityContext.Capabilities)
	assert.Equal(t, []corev1.Capability{"IPC_LOCK", "SYS_RESOURCE", "NET_RAW"}, container.SecurityContext.Capabilities.Add)

	require.Len(t, container.VolumeMounts, 1)
	assert.Equal(t, "hugepage", container.VolumeMounts[0].Name)
	assert.Equal(t, "/dev/hugepages", container.VolumeMounts[0].MountPath)
	require.Len(t, pod.Spec.Volumes, 1)
	require.NotNil(t, pod.Spec.Volumes[0].EmptyDir)
	assert.Equal(t, corev1.StorageMediumHugePages, pod.Spec.Volumes[0].EmptyDir.Medium)
}
