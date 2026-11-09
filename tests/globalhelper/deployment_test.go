package globalhelper

import (
	"errors"
	"testing"
	"time"

	egiClients "github.com/openshift-kni/eco-goinfra/pkg/clients"
	"github.com/stretchr/testify/assert"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
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
				Replicas:          desiredReplicas,
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

		fakeClient := egiClients.GetTestClients(egiClients.TestClientParams{
			K8sMockObjects: runtimeObjects,
		})
		deployReady, err := IsDeploymentReady(fakeClient, "test-namespace", "test-deployment")
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

		fakeClient := egiClients.GetTestClients(egiClients.TestClientParams{
			K8sMockObjects: runtimeObjects,
		})

		err := createAndWaitUntilDeploymentIsReady(fakeClient,
			testDeployment, 5)
		assert.Equal(t, testCase.expectedError, err)
	}
}

func TestCreateAndWaitUntilDeploymentIsReadyFailsFast(t *testing.T) {
	generatePendingDeployment := func(conditions []appsv1.DeploymentCondition) *appsv1.Deployment {
		replicas := int32(1)

		return &appsv1.Deployment{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "pending-deployment",
				Namespace: "test-namespace",
			},
			Spec: appsv1.DeploymentSpec{Replicas: &replicas},
			Status: appsv1.DeploymentStatus{
				Conditions: conditions,
			},
		}
	}

	oldInterval := deploymentPollInterval

	deploymentPollInterval = time.Millisecond

	t.Cleanup(func() { deploymentPollInterval = oldInterval })

	testCases := []struct {
		name              string
		conditions        []appsv1.DeploymentCondition
		expectedErrorText string
	}{
		{
			name: "replica failure is fatal",
			conditions: []appsv1.DeploymentCondition{
				{Type: appsv1.DeploymentProgressing, Status: v1.ConditionTrue},
				{Type: appsv1.DeploymentReplicaFailure, Status: v1.ConditionTrue},
			},
			expectedErrorText: "not schedulable",
		},
		{
			name:              "a deployment that never runs times out",
			expectedErrorText: "timed out after",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			fakeClient := egiClients.GetTestClients(egiClients.TestClientParams{})
			deployment := generatePendingDeployment(testCase.conditions)

			err := createAndWaitUntilDeploymentIsReady(fakeClient, deployment, 20*time.Millisecond)
			assert.ErrorContains(t, err, testCase.expectedErrorText)
		})
	}
}

func TestHasReplicaFailure(t *testing.T) {
	testCases := []struct {
		name       string
		conditions []appsv1.DeploymentCondition
		expected   bool
	}{
		{
			name:     "no conditions",
			expected: false,
		},
		{
			name:       "replica failure reported as true",
			conditions: []appsv1.DeploymentCondition{{Type: appsv1.DeploymentReplicaFailure, Status: v1.ConditionTrue}},
			expected:   true,
		},
		{
			name:       "replica failure reported as false",
			conditions: []appsv1.DeploymentCondition{{Type: appsv1.DeploymentReplicaFailure, Status: v1.ConditionFalse}},
			expected:   false,
		},
		{
			name:       "unrelated condition reported as true",
			conditions: []appsv1.DeploymentCondition{{Type: appsv1.DeploymentProgressing, Status: v1.ConditionTrue}},
			expected:   false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			deployment := &appsv1.Deployment{
				Status: appsv1.DeploymentStatus{Conditions: testCase.conditions},
			}

			assert.Equal(t, testCase.expected, hasReplicaFailure(deployment))
		})
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

	fakeClient := egiClients.GetTestClients(egiClients.TestClientParams{
		K8sMockObjects: runtimeObjects,
	})

	testDeployment, err := getRunningDeployment(fakeClient, "test-namespace", "test-deployment")
	assert.Nil(t, err)
	assert.Equal(t, "test-deployment", testDeployment.Name)
	assert.Equal(t, "test-namespace", testDeployment.Namespace)

	// Test deployment not found
	testDeployment, err = getRunningDeployment(fakeClient, "test-namespace", "test-deployment2")
	assert.NotNil(t, err)
	assert.Nil(t, testDeployment)
}

func TestDeleteDeployment(t *testing.T) {
	generateDeployment := func() *appsv1.Deployment {
		return &appsv1.Deployment{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test-deployment",
				Namespace: "test-namespace",
				Labels:    map[string]string{"app": "test-app"},
			},
			Status: appsv1.DeploymentStatus{},
			Spec: appsv1.DeploymentSpec{
				Template: v1.PodTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Labels: map[string]string{"app": "test-app"},
					},
				},
			},
		}
	}

	var runtimeObjects []runtime.Object
	runtimeObjects = append(runtimeObjects, generateDeployment())

	fakeClient := egiClients.GetTestClients(egiClients.TestClientParams{
		K8sMockObjects: runtimeObjects,
	})

	err := deleteDeployment(fakeClient, "test-namespace", "test-deployment")
	assert.Nil(t, err)

	// Test deployment not found
	err = deleteDeployment(fakeClient, "test-namespace", "test-deployment2")
	assert.Nil(t, err)
}
