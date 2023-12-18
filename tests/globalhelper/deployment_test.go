package globalhelper

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	k8sfake "k8s.io/client-go/kubernetes/fake"
)

func TestIsDeploymentReady(t *testing.T) {
	generateDeployment := func(availableReplicas, desiredReplicas, readyReplicas int32) *appsv1.Deployment {
		return &appsv1.Deployment{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-deployment",
				Namespace: "test-namespace",
			},
			Status: appsv1.DeploymentStatus{
				AvailableReplicas: availableReplicas,
				ReadyReplicas:     readyReplicas,
			},
			Spec: appsv1.DeploymentSpec{
				Replicas: &desiredReplicas,
			},
		}
	}

	testCases := []struct {
		deploymentExists  bool
		expectReady       bool
		expectedError     error
		availableReplicas int32
		desiredReplicas   int32
		readyReplicas     int32
	}{
		{
			deploymentExists:  true,
			expectReady:       true,
			expectedError:     nil,
			availableReplicas: 1,
			desiredReplicas:   1,
			readyReplicas:     1,
		},
		{
			deploymentExists:  false,
			expectReady:       false,
			expectedError:     errors.New("deployment not found"),
			availableReplicas: 0,
			desiredReplicas:   0,
			readyReplicas:     0,
		},
		{
			deploymentExists:  true,
			expectReady:       false,
			expectedError:     nil,
			availableReplicas: 0,
			desiredReplicas:   1,
			readyReplicas:     0,
		},
	}

	for _, testCase := range testCases {
		var runtimeObjects []runtime.Object
		if testCase.deploymentExists {
			runtimeObjects = append(runtimeObjects, generateDeployment(testCase.availableReplicas,
				testCase.desiredReplicas, testCase.readyReplicas))
		}

		client := k8sfake.NewSimpleClientset(runtimeObjects...)
		deployReady, err := IsDeploymentReady(client.AppsV1(), "test-namespace", "test-deployment")
		assert.Equal(t, testCase.expectReady, deployReady)

		if testCase.expectedError != nil {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err)
		}
	}
}

func TestCreateAndWaitUntilDeploymentIsReady(t *testing.T) {
	generateDeployment := func(availableReplicas, desiredReplicas, readyReplicas int32) *appsv1.Deployment {
		return &appsv1.Deployment{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-deployment",
				Namespace: "test-namespace",
			},
			Status: appsv1.DeploymentStatus{
				AvailableReplicas: availableReplicas,
				ReadyReplicas:     readyReplicas,
			},
			Spec: appsv1.DeploymentSpec{
				Replicas: &desiredReplicas,
			},
		}
	}

	testCases := []struct {
		expectedError     error
		availableReplicas int32
		desiredReplicas   int32
		readyReplicas     int32
	}{
		{
			expectedError:     nil,
			availableReplicas: 1,
			desiredReplicas:   1,
			readyReplicas:     1,
		},
		{
			expectedError:     nil,
			availableReplicas: 0,
			desiredReplicas:   0,
			readyReplicas:     0,
		},
	}

	for _, testCase := range testCases {
		// Create fake deployment
		var runtimeObjects []runtime.Object

		testDeployment := generateDeployment(testCase.availableReplicas, testCase.desiredReplicas, testCase.readyReplicas)
		runtimeObjects = append(runtimeObjects, testDeployment)
		client := k8sfake.NewSimpleClientset(runtimeObjects...)

		err := createAndWaitUntilDeploymentIsReady(client.AppsV1(),
			testDeployment, 5)
		assert.Equal(t, testCase.expectedError, err)
	}
}

func TestGetRunningDeployment(t *testing.T) {
	generateDeployment := func() *appsv1.Deployment {
		return &appsv1.Deployment{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-deployment",
				Namespace: "test-namespace",
			},
			Status: appsv1.DeploymentStatus{},
			Spec:   appsv1.DeploymentSpec{},
		}
	}

	var runtimeObjects []runtime.Object
	runtimeObjects = append(runtimeObjects, generateDeployment())
	client := k8sfake.NewSimpleClientset(runtimeObjects...)

	testDeployment, err := getRunningDeployment(client.AppsV1(), "test-namespace", "test-deployment")
	assert.Nil(t, err)
	assert.Equal(t, "test-deployment", testDeployment.Name)
	assert.Equal(t, "test-namespace", testDeployment.Namespace)

	// Test deployment not found
	testDeployment, err = getRunningDeployment(client.AppsV1(), "test-namespace", "test-deployment2")
	assert.NotNil(t, err)
	assert.Nil(t, testDeployment)
}

func TestDeleteDeployment(t *testing.T) {
	generateDeployment := func() *appsv1.Deployment {
		return &appsv1.Deployment{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-deployment",
				Namespace: "test-namespace",
			},
			Status: appsv1.DeploymentStatus{},
			Spec:   appsv1.DeploymentSpec{},
		}
	}

	var runtimeObjects []runtime.Object
	runtimeObjects = append(runtimeObjects, generateDeployment())
	client := k8sfake.NewSimpleClientset(runtimeObjects...)

	err := deleteDeployment(client.AppsV1(), "test-namespace", "test-deployment")
	assert.Nil(t, err)

	// Test deployment not found
	err = deleteDeployment(client.AppsV1(), "test-namespace", "test-deployment2")
	assert.Nil(t, err)
}
