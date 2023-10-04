package globalhelper

import (
	"testing"

	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	k8sfake "k8s.io/client-go/kubernetes/fake"
)

func TestCreateAndWaitUntilDaemonSetIsReady(t *testing.T) {
	generateDaemonset := func(numAvailable, numScheduled, numUnavailable int) *appsv1.DaemonSet {
		return &appsv1.DaemonSet{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-daemonset",
				Namespace: "default",
			},
			Status: appsv1.DaemonSetStatus{
				NumberAvailable:        int32(numAvailable),
				DesiredNumberScheduled: int32(numScheduled),
				NumberUnavailable:      int32(numUnavailable),
			},
		}
	}

	testCases := []struct {
		name           string
		numAvailable   int
		numScheduled   int
		numUnavailable int
		wantErr        bool
	}{
		{
			name:           "Daemonset is ready",
			numAvailable:   1,
			numScheduled:   1,
			numUnavailable: 0,
			wantErr:        false,
		},
		{
			name:           "Daemonset is not ready",
			numAvailable:   0,
			numScheduled:   1,
			numUnavailable: 1,
			wantErr:        false,
		},
	}

	for _, testCase := range testCases {
		// Create fake daemonset
		var runtimeObjects []runtime.Object
		runtimeObjects = append(runtimeObjects, generateDaemonset(testCase.numAvailable, testCase.numScheduled, testCase.numUnavailable))
		client := k8sfake.NewSimpleClientset(runtimeObjects...)

		t.Run(testCase.name, func(t *testing.T) {
			err := createAndWaitUntilDaemonSetIsReady(client.AppsV1(),
				generateDaemonset(testCase.numAvailable, testCase.numScheduled, testCase.numUnavailable), 5)
			if (err != nil) != testCase.wantErr {
				t.Errorf("CreateAndWaitUntilDaemonSetIsReady() error = %v, wantErr %v", err, testCase.wantErr)

				return
			}
		})
	}
}

func TestIsDaemonsetReady(t *testing.T) {
	generateDaemonset := func(numAvailable, numScheduled, numUnavailable int) *appsv1.DaemonSet {
		return &appsv1.DaemonSet{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-daemonset",
				Namespace: "default",
			},
			Status: appsv1.DaemonSetStatus{
				NumberAvailable:        int32(numAvailable),
				DesiredNumberScheduled: int32(numScheduled),
				NumberUnavailable:      int32(numUnavailable),
			},
		}
	}

	testCases := []struct {
		name           string
		numAvailable   int
		numScheduled   int
		numUnavailable int
		wantErr        bool
		expectedOutput bool
	}{
		{
			name:           "Daemonset is ready",
			numAvailable:   1,
			numScheduled:   1,
			numUnavailable: 0,
			wantErr:        false,
			expectedOutput: true,
		},
		{
			name:           "Daemonset is not ready",
			numAvailable:   0,
			numScheduled:   1,
			numUnavailable: 1,
			wantErr:        false,
			expectedOutput: false,
		},
	}

	for _, testCase := range testCases {
		// Create fake daemonset
		var runtimeObjects []runtime.Object
		runtimeObjects = append(runtimeObjects, generateDaemonset(testCase.numAvailable, testCase.numScheduled, testCase.numUnavailable))
		client := k8sfake.NewSimpleClientset(runtimeObjects...)

		t.Run(testCase.name, func(t *testing.T) {
			got, err := isDaemonSetReady(client.AppsV1(), "default", "test-daemonset")
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
