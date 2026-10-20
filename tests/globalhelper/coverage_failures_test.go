//nolint:lll // Inline fake API fixtures and reactor callbacks are intentionally compact.
package globalhelper

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/onsi/gomega"
	egiClients "github.com/openshift-kni/eco-goinfra/pkg/clients"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	k8stesting "k8s.io/client-go/testing"
	"k8s.io/utils/ptr"
)

func TestCoverageWorkloadFailurePaths(t *testing.T) {
	gomega.RegisterTestingT(t)

	var err error
	old := []time.Duration{deploymentPollInterval, deploymentReadyProbeTimeout, podReadyPollInterval,
		daemonSetReadyPollInterval, replicaSetReadyPollInterval, statefulSetReadyPollInterval}
	deploymentPollInterval, deploymentReadyProbeTimeout, podReadyPollInterval = time.Millisecond, time.Millisecond, time.Millisecond
	daemonSetReadyPollInterval, replicaSetReadyPollInterval, statefulSetReadyPollInterval = time.Millisecond, time.Millisecond, time.Millisecond

	t.Cleanup(func() {
		deploymentPollInterval, deploymentReadyProbeTimeout, podReadyPollInterval = old[0], old[1], old[2]
		daemonSetReadyPollInterval, replicaSetReadyPollInterval, statefulSetReadyPollInterval = old[3], old[4], old[5]
	})

	dep := &appsv1.Deployment{ObjectMeta: metav1.ObjectMeta{Name: "dep", Namespace: "ns"}, Spec: appsv1.DeploymentSpec{Replicas: ptr.To[int32](1)}}
	k8s, _ := coverageFixture(t)
	k8s.PrependReactor("create", "deployments", func(k8stesting.Action) (bool, runtime.Object, error) { return true, nil, errors.New("create denied") })
	require.Error(t, CreateAndWaitUntilDeploymentIsReady(dep, time.Millisecond))
	coverageFixture(t)
	require.Error(t, createAndWaitUntilDeploymentIsReady(egiClient, dep, time.Millisecond))
	k8s, _ = coverageFixture(t)
	k8s.PrependReactor("get", "deployments", func(k8stesting.Action) (bool, runtime.Object, error) { return true, nil, errors.New("get denied") })
	require.Error(t, createAndWaitUntilDeploymentIsReady(egiClient, dep, 4*time.Millisecond))
	failedDeployment := dep.DeepCopy()
	failedDeployment.Status.Conditions = []appsv1.DeploymentCondition{{Type: appsv1.DeploymentReplicaFailure, Status: corev1.ConditionTrue}}
	k8s, _ = coverageFixture(t)
	k8s.PrependReactor("get", "deployments", func(k8stesting.Action) (bool, runtime.Object, error) { return true, failedDeployment, nil })
	require.Error(t, createAndWaitUntilDeploymentIsReady(egiClient, dep, time.Second))
	readyButUnavailable := dep.DeepCopy()
	readyButUnavailable.Status.ReadyReplicas = 1
	readyButUnavailable.Status.AvailableReplicas = 0
	k8s, _ = coverageFixture(t)
	k8s.PrependReactor("get", "deployments", func(k8stesting.Action) (bool, runtime.Object, error) {
		return true, readyButUnavailable.DeepCopy(), nil
	})
	require.Error(t, createAndWaitUntilDeploymentIsReady(egiClient, dep, 5*time.Millisecond))

	daemonSet := &appsv1.DaemonSet{ObjectMeta: metav1.ObjectMeta{Name: "ds", Namespace: "ns"}, Spec: appsv1.DaemonSetSpec{Template: corev1.PodTemplateSpec{Spec: corev1.PodSpec{Containers: []corev1.Container{{Name: "c"}}}}}}
	k8s, _ = coverageFixture(t)
	k8s.PrependReactor("create", "daemonsets", func(k8stesting.Action) (bool, runtime.Object, error) { return true, nil, errors.New("create denied") })
	require.Error(t, CreateAndWaitUntilDaemonSetIsReady(daemonSet, time.Millisecond))
	coverageFixture(t)
	broken := failingControllerClient{Client: egiClient.Client, getErr: errors.New("get denied")}
	egiClient.Client, apiclient.Client = broken, broken
	_, err = isDaemonSetReady(egiClient, "ns", "missing")
	require.Error(t, err)
	_, err = getDaemonSetPullPolicy(daemonSet, egiClient)
	require.Error(t, err)
	_, err = getRunningDaemonset(daemonSet, egiClient)
	require.Error(t, err)

	rs := &appsv1.ReplicaSet{ObjectMeta: metav1.ObjectMeta{Name: "rs", Namespace: "ns"}, Spec: appsv1.ReplicaSetSpec{Replicas: ptr.To[int32](1)}}
	k8s, _ = coverageFixture(t)
	k8s.PrependReactor("create", "replicasets", func(k8stesting.Action) (bool, runtime.Object, error) { return true, nil, errors.New("create denied") })
	require.Error(t, CreateAndWaitUntilReplicaSetIsReady(rs, time.Millisecond))
	k8s, _ = coverageFixture(t)
	k8s.PrependReactor("get", "replicasets", func(k8stesting.Action) (bool, runtime.Object, error) { return true, nil, errors.New("get denied") })
	_, err = isReplicaSetReady("ns", "missing")
	require.Error(t, err)

	sts := &appsv1.StatefulSet{ObjectMeta: metav1.ObjectMeta{Name: "sts", Namespace: "ns"}, Spec: appsv1.StatefulSetSpec{Replicas: ptr.To[int32](1)}}
	k8s, _ = coverageFixture(t)
	k8s.PrependReactor("create", "statefulsets", func(k8stesting.Action) (bool, runtime.Object, error) { return true, nil, errors.New("create denied") })
	require.Error(t, CreateAndWaitUntilStatefulSetIsReady(sts, time.Millisecond))
	coverageFixture(t)
	broken = failingControllerClient{Client: egiClient.Client, getErr: errors.New("get denied")}
	egiClient.Client, apiclient.Client = broken, broken
	_, err = isStatefulSetReady(egiClient, "ns", "missing")
	require.Error(t, err)
	_, err = getRunningStatefulSet(egiClient, "ns", "missing")
	require.Error(t, err)

	k8s, _ = coverageFixture(t)
	k8s.PrependReactor("get", "pods", func(k8stesting.Action) (bool, runtime.Object, error) { return true, nil, errors.New("pod get denied") })
	_, err = isPodReady("ns", "missing")
	require.Error(t, err)
	_ = k8s
}

func TestCoverageWorkloadRetryBranches(t *testing.T) {
	gomega.RegisterTestingT(t)
	old := []time.Duration{podReadyPollInterval, daemonSetReadyPollInterval, replicaSetReadyPollInterval, statefulSetReadyPollInterval, deploymentPollInterval}
	oldDSChecker, oldSTSChecker, oldDeploymentChecker := daemonSetReadyChecker, statefulSetReadyChecker, deploymentReadyChecker
	podReadyPollInterval, daemonSetReadyPollInterval, replicaSetReadyPollInterval, statefulSetReadyPollInterval = time.Millisecond, time.Millisecond, time.Millisecond, time.Millisecond
	deploymentPollInterval = time.Millisecond

	t.Cleanup(func() {
		podReadyPollInterval, daemonSetReadyPollInterval, replicaSetReadyPollInterval, statefulSetReadyPollInterval = old[0], old[1], old[2], old[3]
		deploymentPollInterval = old[4]
		daemonSetReadyChecker, statefulSetReadyChecker, deploymentReadyChecker = oldDSChecker, oldSTSChecker, oldDeploymentChecker
	})

	pod := &corev1.Pod{ObjectMeta: metav1.ObjectMeta{Name: "pod", Namespace: "ns"}, Spec: corev1.PodSpec{Containers: []corev1.Container{{Name: "c"}}}, Status: corev1.PodStatus{Phase: corev1.PodRunning, ContainerStatuses: []corev1.ContainerStatus{{Name: "c", Ready: true}}}}
	k8s, _ := coverageFixture(t)
	firstPodGet := true

	k8s.PrependReactor("get", "pods", func(k8stesting.Action) (bool, runtime.Object, error) {
		if firstPodGet {
			firstPodGet = false

			return true, nil, errors.New("transient pod get failure")
		}

		return false, nil, nil
	})
	require.NoError(t, CreateAndWaitUntilPodIsReady(pod, time.Second))
	k8s, _ = coverageFixture(t)
	podGetCalls := 0

	k8s.PrependReactor("get", "pods", func(k8stesting.Action) (bool, runtime.Object, error) {
		podGetCalls++
		if podGetCalls == 2 {
			return true, nil, errors.New("transient readiness get failure")
		}

		return true, pod.DeepCopy(), nil
	})
	require.NoError(t, CreateAndWaitUntilPodIsReady(pod, time.Second))

	daemonSet := &appsv1.DaemonSet{ObjectMeta: metav1.ObjectMeta{Name: "ds", Namespace: "ns"}}

	coverageFixture(t)
	firstReadyError := true
	daemonSetReadyChecker = func(*egiClients.Settings, string, string) (bool, error) {
		if firstReadyError {
			firstReadyError = false

			return false, errors.New("transient daemonset readiness error")
		}

		return true, nil
	}

	require.NoError(t, CreateAndWaitUntilDaemonSetIsReady(daemonSet, time.Second))

	replicaSet := &appsv1.ReplicaSet{ObjectMeta: metav1.ObjectMeta{Name: "rs", Namespace: "ns"}, Spec: appsv1.ReplicaSetSpec{Replicas: ptr.To[int32](1)}, Status: appsv1.ReplicaSetStatus{AvailableReplicas: 1}}
	k8s, _ = coverageFixture(t)
	firstGet := true

	k8s.PrependReactor("get", "replicasets", func(k8stesting.Action) (bool, runtime.Object, error) {
		if firstGet {
			firstGet = false

			return true, nil, errors.New("transient replicaset get failure")
		}

		return false, nil, nil
	})
	require.NoError(t, CreateAndWaitUntilReplicaSetIsReady(replicaSet, time.Second))

	sts := &appsv1.StatefulSet{ObjectMeta: metav1.ObjectMeta{Name: "sts", Namespace: "ns"}}

	coverageFixture(t)
	firstReadyError = true
	statefulSetReadyChecker = func(*egiClients.Settings, string, string) (bool, error) {
		if firstReadyError {
			firstReadyError = false

			return false, errors.New("transient statefulset readiness error")
		}

		return true, nil
	}

	require.NoError(t, CreateAndWaitUntilStatefulSetIsReady(sts, time.Second))

	dep := &appsv1.Deployment{ObjectMeta: metav1.ObjectMeta{Name: "dep", Namespace: "ns"}, Spec: appsv1.DeploymentSpec{Replicas: ptr.To[int32](1)}, Status: appsv1.DeploymentStatus{ReadyReplicas: 1}}

	coverageFixture(t)
	firstReadyError = true
	deploymentReadyChecker = func(*egiClients.Settings, string, string) (bool, error) {
		if firstReadyError {
			firstReadyError = false

			return false, errors.New("transient deployment readiness error")
		}

		return true, nil
	}

	require.NoError(t, CreateAndWaitUntilDeploymentIsReady(dep, time.Second))
}
