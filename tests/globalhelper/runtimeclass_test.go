package globalhelper

import (
	"testing"

	"github.com/stretchr/testify/assert"
	nodev1 "k8s.io/api/node/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	k8sfake "k8s.io/client-go/kubernetes/fake"
)

func TestIsRTCCreated(t *testing.T) {
	testCases := []struct {
		rtcAlreadyExists bool
	}{
		{rtcAlreadyExists: false},
		{rtcAlreadyExists: true},
	}

	for _, testCase := range testCases {
		var runtimeObjects []runtime.Object

		if testCase.rtcAlreadyExists {
			// Create a fake runtime class object
			runtimeObjects = append(runtimeObjects, &nodev1.RuntimeClass{
				ObjectMeta: metav1.ObjectMeta{
					Name: "testRTC",
					UID:  "testUID",
				},
			})
		}

		// Create a fake clientset
		client := k8sfake.NewClientset(runtimeObjects...)
		isCreated, err := isRtcCreated(client, &nodev1.RuntimeClass{
			ObjectMeta: metav1.ObjectMeta{
				Name: "testRTC",
				UID:  "testUID",
			},
		})

		assert.Equal(t, testCase.rtcAlreadyExists, isCreated)

		if testCase.rtcAlreadyExists {
			assert.Nil(t, err)
		} else {
			assert.True(t, k8serrors.IsNotFound(err))
		}
	}
}

func TestIsRtcDeleted(t *testing.T) {
	testCases := []struct {
		rtcAlreadyDeleted bool
	}{
		{rtcAlreadyDeleted: false},
		{rtcAlreadyDeleted: true},
	}

	for _, testCase := range testCases {
		var runtimeObjects []runtime.Object

		if !testCase.rtcAlreadyDeleted {
			// Create a fake runtime class object
			runtimeObjects = append(runtimeObjects, &nodev1.RuntimeClass{
				ObjectMeta: metav1.ObjectMeta{
					Name: "testRTC",
					UID:  "testUID",
				},
			})
		}

		// Create a fake clientset
		client := k8sfake.NewClientset(runtimeObjects...)
		isDeleted, err := isRtcDeleted(client, &nodev1.RuntimeClass{
			ObjectMeta: metav1.ObjectMeta{
				Name: "testRTC",
				UID:  "testUID",
			},
		})

		assert.Equal(t, testCase.rtcAlreadyDeleted, isDeleted)
		assert.Nil(t, err)
	}
}
