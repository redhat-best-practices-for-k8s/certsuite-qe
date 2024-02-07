package globalhelper

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	k8sfake "k8s.io/client-go/kubernetes/fake"
)

func TestCreateAndWaitUntilDaemonSetIsReady(t *testing.T) {
	generateDaemonset := func(numAvailable, numScheduled, numUnavailable, numReady int) *appsv1.DaemonSet {
		return &appsv1.DaemonSet{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-daemonset",
				Namespace: "default",
			},
			Status: appsv1.DaemonSetStatus{
				NumberAvailable:        int32(numAvailable),
				DesiredNumberScheduled: int32(numScheduled),
				NumberUnavailable:      int32(numUnavailable),
				NumberReady:            int32(numReady),
			},
		}
	}

	testCases := []struct {
		name           string
		numAvailable   int
		numScheduled   int
		numUnavailable int
		numReady       int
		wantErr        bool
	}{
		{
			name:           "Daemonset is ready",
			numAvailable:   1,
			numScheduled:   1,
			numReady:       1,
			numUnavailable: 0,
			wantErr:        false,
		},
		{
			name:           "Daemonset is not ready",
			numAvailable:   0,
			numReady:       0,
			numScheduled:   1,
			numUnavailable: 1,
			wantErr:        false,
		},
	}

	for _, testCase := range testCases {
		// Create fake daemonset
		var runtimeObjects []runtime.Object
		runtimeObjects = append(runtimeObjects, generateDaemonset(testCase.numAvailable,
			testCase.numScheduled, testCase.numUnavailable, testCase.numReady))
		client := k8sfake.NewSimpleClientset(runtimeObjects...)

		t.Run(testCase.name, func(t *testing.T) {
			err := createAndWaitUntilDaemonSetIsReady(client.AppsV1(), client.CoreV1(),
				generateDaemonset(testCase.numAvailable, testCase.numScheduled, testCase.numUnavailable, testCase.numReady), 5)
			if (err != nil) != testCase.wantErr {
				t.Errorf("CreateAndWaitUntilDaemonSetIsReady() error = %v, wantErr %v", err, testCase.wantErr)

				return
			}
		})
	}
}

func TestIsDaemonsetReady(t *testing.T) {
	generateDaemonset := func(numAvailable, numScheduled, numUnavailable, numReady int) *appsv1.DaemonSet {
		return &appsv1.DaemonSet{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-daemonset",
				Namespace: "default",
			},
			Status: appsv1.DaemonSetStatus{
				NumberAvailable:        int32(numAvailable),
				DesiredNumberScheduled: int32(numScheduled),
				NumberUnavailable:      int32(numUnavailable),
				NumberReady:            int32(numReady),
			},
		}
	}

	testCases := []struct {
		name           string
		numAvailable   int
		numScheduled   int
		numUnavailable int
		numReady       int
		wantErr        bool
		expectedOutput bool
	}{
		{
			name:           "Daemonset is ready",
			numAvailable:   1,
			numReady:       1,
			numScheduled:   1,
			numUnavailable: 0,
			wantErr:        false,
			expectedOutput: true,
		},
		{
			name:           "Daemonset is not ready",
			numAvailable:   0,
			numReady:       0,
			numScheduled:   1,
			numUnavailable: 1,
			wantErr:        false,
			expectedOutput: false,
		},
	}

	for _, testCase := range testCases {
		// Create fake daemonset
		var runtimeObjects []runtime.Object
		runtimeObjects = append(runtimeObjects, generateDaemonset(testCase.numAvailable,
			testCase.numScheduled, testCase.numUnavailable, testCase.numReady))
		runtimeObjects = append(runtimeObjects, &corev1.Node{
			ObjectMeta: metav1.ObjectMeta{
				Name: "test-node",
				Labels: map[string]string{
					"node-role.kubernetes.io/worker-cnf": "",
				},
			},
		})

		client := k8sfake.NewSimpleClientset(runtimeObjects...)

		t.Run(testCase.name, func(t *testing.T) {
			got, err := isDaemonSetReady(client.AppsV1(), client.CoreV1(), "default", "test-daemonset")
			if (err != nil) != testCase.wantErr {
				t.Errorf("isDaemonsetReady() error = %v, wantErr %v", err, testCase.wantErr)

				return
			}

			if got != testCase.expectedOutput {
				t.Errorf("isDaemonsetReady() = %v, want %v", got, testCase.expectedOutput)
			}
		})
	}
}

func TestGetDaemonSetPullPolicy(t *testing.T) {
	generateDaemonset := func(imagePullPolicy corev1.PullPolicy) *appsv1.DaemonSet {
		return &appsv1.DaemonSet{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-daemonset",
				Namespace: "test-namespace",
			},
			Spec: appsv1.DaemonSetSpec{
				Template: corev1.PodTemplateSpec{
					Spec: corev1.PodSpec{
						Containers: []corev1.Container{
							{
								ImagePullPolicy: imagePullPolicy,
							},
						},
					},
				},
			},
		}
	}

	testCases := []struct {
		pullPolicy corev1.PullPolicy
	}{
		{
			pullPolicy: corev1.PullAlways,
		},
		{
			pullPolicy: corev1.PullIfNotPresent,
		},
	}

	for _, testCase := range testCases {
		// Create fake daemonset
		var runtimeObjects []runtime.Object

		testDaemonset := generateDaemonset(testCase.pullPolicy)
		runtimeObjects = append(runtimeObjects, testDaemonset)
		client := k8sfake.NewSimpleClientset(runtimeObjects...)

		pullPolicy, err := getDaemonSetPullPolicy(testDaemonset, client.AppsV1())
		assert.Equal(t, testCase.pullPolicy, pullPolicy)
		assert.Nil(t, err)
	}
}

func TestGetRunningDaemonset(t *testing.T) {
	generateDaemonset := func() *appsv1.DaemonSet {
		return &appsv1.DaemonSet{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-daemonset",
				Namespace: "test-namespace",
			},
		}
	}

	testCases := []struct {
		alreadyExists bool
		expectedError error
	}{
		{
			alreadyExists: true,
			expectedError: nil,
		},
		{
			alreadyExists: false,
			expectedError: errors.New("daemonset not found"),
		},
	}

	for _, testCase := range testCases {
		// Create fake daemonset
		var runtimeObjects []runtime.Object

		if testCase.alreadyExists {
			testDaemonset := generateDaemonset()
			runtimeObjects = append(runtimeObjects, testDaemonset)
		}

		client := k8sfake.NewSimpleClientset(runtimeObjects...)
		daemonset, err := getRunningDaemonset(generateDaemonset(), client.AppsV1())

		if testCase.expectedError != nil {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err)
			assert.Equal(t, "test-daemonset", daemonset.Name)
		}
	}
}
