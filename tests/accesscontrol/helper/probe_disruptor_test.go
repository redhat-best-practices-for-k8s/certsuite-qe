package helper

import (
	"errors"
	"testing"

	tsparams "github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/accesscontrol/parameters"
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/globalhelper"
	"github.com/stretchr/testify/assert"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8sruntime "k8s.io/apimachinery/pkg/runtime"
	k8sfake "k8s.io/client-go/kubernetes/fake"
	k8stesting "k8s.io/client-go/testing"
)

func TestProbeDisruptCommand(t *testing.T) {
	t.Parallel()

	cmd := probeDisruptCommand()
	assert.Len(t, cmd, 3)
	assert.Equal(t, "/bin/sh", cmd[0])
	assert.Equal(t, "-c", cmd[1])
	assert.Contains(t, cmd[2], "sleep 1")
	assert.Contains(t, cmd[2], "kill -9 1")
	assert.NotContains(t, cmd[2], "python3")
	assert.NotContains(t, cmd[2], "time.sleep")
	assert.NotContains(t, cmd[2], "dd if=")
}

func TestIsProbeDaemonSetReady(t *testing.T) {
	t.Parallel()

	assert.False(t, isProbeDaemonSetReady(nil))
	assert.False(t, isProbeDaemonSetReady(&appsv1.DaemonSetStatus{}))
	assert.False(t, isProbeDaemonSetReady(&appsv1.DaemonSetStatus{
		DesiredNumberScheduled: 3,
		CurrentNumberScheduled: 3,
		NumberAvailable:        3,
		NumberReady:            2,
	}))
	assert.False(t, isProbeDaemonSetReady(&appsv1.DaemonSetStatus{
		DesiredNumberScheduled: 3,
		CurrentNumberScheduled: 3,
		NumberAvailable:        3,
		NumberReady:            3,
		NumberMisscheduled:     1,
	}))
	assert.True(t, isProbeDaemonSetReady(&appsv1.DaemonSetStatus{
		DesiredNumberScheduled: 3,
		CurrentNumberScheduled: 3,
		NumberAvailable:        3,
		NumberReady:            3,
	}))
}

func TestProbeDaemonSetReady(t *testing.T) {
	t.Run("returns true when the daemonset is fully ready", func(t *testing.T) {
		client := k8sfake.NewClientset(testProbeDaemonSet(2, 2, 2, 2, 0))
		globalhelper.SetTestK8sAPIClient(client)
		t.Cleanup(globalhelper.UnsetTestK8sAPIClient)

		ready, err := probeDaemonSetReady(t.Context())

		assert.NoError(t, err)
		assert.True(t, ready)
	})

	t.Run("returns false when NumberReady is below Desired", func(t *testing.T) {
		client := k8sfake.NewClientset(testProbeDaemonSet(2, 1, 1, 2, 0))
		globalhelper.SetTestK8sAPIClient(client)
		t.Cleanup(globalhelper.UnsetTestK8sAPIClient)

		ready, err := probeDaemonSetReady(t.Context())

		assert.NoError(t, err)
		assert.False(t, ready)
	})

	t.Run("returns error if the daemonset is missing", func(t *testing.T) {
		globalhelper.SetTestK8sAPIClient(k8sfake.NewClientset())
		t.Cleanup(globalhelper.UnsetTestK8sAPIClient)

		ready, err := probeDaemonSetReady(t.Context())

		assert.Error(t, err)
		assert.False(t, ready)
	})
}

func TestFirstMatchingProbe(t *testing.T) {
	t.Parallel()

	pending := testProbePod("pending-probe", "worker-0", corev1.PodPending)
	running := testProbePod("running-probe", "worker-0", corev1.PodRunning)
	otherNode := testProbePod("other-node-probe", "worker-1", corev1.PodRunning)

	testCases := []struct {
		name     string
		pods     []corev1.Pod
		nodeName string
		wantName string
	}{
		{
			name:     "skips pending when looking for running",
			pods:     []corev1.Pod{*pending, *running},
			nodeName: "worker-0",
			wantName: "running-probe",
		},
		{
			name:     "ignores running probes on other nodes",
			pods:     []corev1.Pod{*otherNode, *running},
			nodeName: "worker-0",
			wantName: "running-probe",
		},
		{
			name:     "returns nil when nothing matches",
			pods:     []corev1.Pod{*pending, *otherNode},
			nodeName: "worker-0",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			got := firstMatchingProbe(testCase.pods, testCase.nodeName)
			if testCase.wantName == "" {
				assert.Nil(t, got)

				return
			}

			if assert.NotNil(t, got) {
				assert.Equal(t, testCase.wantName, got.Name)
			}
		})
	}
}

func TestFindProbePodOnNode(t *testing.T) {
	t.Run("returns running pod on the target node", func(t *testing.T) {
		client := k8sfake.NewClientset(
			testProbePod("pending-probe", "worker-0", corev1.PodPending),
			testProbePod("other-node-probe", "worker-1", corev1.PodRunning),
			testProbePod("running-probe", "worker-0", corev1.PodRunning),
		)
		globalhelper.SetTestK8sAPIClient(client)
		t.Cleanup(globalhelper.UnsetTestK8sAPIClient)

		pod, err := findProbePodOnNode(t.Context(), "worker-0")

		assert.NoError(t, err)

		if assert.NotNil(t, pod) {
			assert.Equal(t, "running-probe", pod.Name)
		}
	})

	t.Run("returns nil if no matching pod is found", func(t *testing.T) {
		globalhelper.SetTestK8sAPIClient(k8sfake.NewClientset())
		t.Cleanup(globalhelper.UnsetTestK8sAPIClient)

		pod, err := findProbePodOnNode(t.Context(), "worker-0")

		assert.NoError(t, err)
		assert.Nil(t, pod)
	})
}

func TestFindProbePodOnNodePropagatesListError(t *testing.T) {
	client := k8sfake.NewClientset()
	client.PrependReactor("list", "pods", func(_ k8stesting.Action) (bool, k8sruntime.Object, error) {
		return true, nil, errors.New("list failed")
	})
	globalhelper.SetTestK8sAPIClient(client)
	t.Cleanup(globalhelper.UnsetTestK8sAPIClient)

	pod, err := findProbePodOnNode(t.Context(), "worker-0")

	assert.Nil(t, pod)
	assert.EqualError(t, err, "list failed")
}

func testProbePod(name, node string, phase corev1.PodPhase) *corev1.Pod {
	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: tsparams.CertsuiteProbeDaemonSetNamespace,
			Labels:    map[string]string{"name": tsparams.CertsuiteProbeDaemonSetName},
		},
		Spec:   corev1.PodSpec{NodeName: node},
		Status: corev1.PodStatus{Phase: phase},
	}
}

func testProbeDaemonSet(desired, ready, available, current, misscheduled int32) *appsv1.DaemonSet {
	return &appsv1.DaemonSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      tsparams.CertsuiteProbeDaemonSetName,
			Namespace: tsparams.CertsuiteProbeDaemonSetNamespace,
		},
		Status: appsv1.DaemonSetStatus{
			DesiredNumberScheduled: desired,
			CurrentNumberScheduled: current,
			NumberReady:            ready,
			NumberAvailable:        available,
			NumberMisscheduled:     misscheduled,
		},
	}
}
