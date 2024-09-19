package globalhelper

import (
	"context"
	"testing"
	"time"

	egiClients "github.com/openshift-kni/eco-goinfra/pkg/clients"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	policyv1 "k8s.io/api/policy/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func generateNamespace(name string) *corev1.Namespace {
	return &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
	}
}

func TestDeleteNamespaceAndWait(t *testing.T) {
	testCases := []struct {
		name          string
		alreadyExists bool
	}{
		{
			name:          "test-namespace",
			alreadyExists: true,
		},
		{
			name:          "test-namespace",
			alreadyExists: false,
		},
	}

	for _, testCase := range testCases {
		var runtimeObjects []runtime.Object
		if testCase.alreadyExists {
			runtimeObjects = append(runtimeObjects, generateNamespace(testCase.name))
		}

		fakeClient := egiClients.GetTestClients(egiClients.TestClientParams{
			K8sMockObjects: runtimeObjects,
		})
		err := deleteNamespaceAndWait(fakeClient, testCase.name, 10*time.Second)
		assert.NoError(t, err)
	}
}

func TestCreateNamespace(t *testing.T) {
	testCases := []struct {
		name          string
		alreadyExists bool
	}{
		{
			name:          "test-namespace",
			alreadyExists: true,
		},
		{
			name:          "test-namespace",
			alreadyExists: false,
		},
	}

	for _, testCase := range testCases {
		var runtimeObjects []runtime.Object
		if testCase.alreadyExists {
			runtimeObjects = append(runtimeObjects, generateNamespace(testCase.name))
		}

		fakeClient := egiClients.GetTestClients(egiClients.TestClientParams{
			K8sMockObjects: runtimeObjects,
		})
		err := createNamespace(testCase.name, fakeClient)
		assert.NoError(t, err)
	}
}

func TestNamespaceExists(t *testing.T) {
	testCases := []struct {
		name          string
		alreadyExists bool
	}{
		{
			name:          "test-namespace",
			alreadyExists: true,
		},
		{
			name:          "test-namespace",
			alreadyExists: false,
		},
	}

	for _, testCase := range testCases {
		var runtimeObjects []runtime.Object
		if testCase.alreadyExists {
			runtimeObjects = append(runtimeObjects, generateNamespace(testCase.name))
		}

		fakeClient := egiClients.GetTestClients(egiClients.TestClientParams{
			K8sMockObjects: runtimeObjects,
		})
		exists, err := namespaceExists(testCase.name, fakeClient)
		assert.NoError(t, err)
		assert.Equal(t, testCase.alreadyExists, exists)
	}
}

func TestCleanPVCs(t *testing.T) {
	generatePVC := func(name, namespace string) *corev1.Pod {
		return &corev1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name:      name,
				Namespace: namespace,
			},
		}
	}

	testCases := []struct {
		name            string
		namespace       string
		namespaceExists bool
	}{
		{
			name:            "test-pvc",
			namespace:       "test-namespace",
			namespaceExists: true,
		},
		{
			name:            "test-pvc",
			namespace:       "test-namespace",
			namespaceExists: false,
		},
	}

	for _, testCase := range testCases {
		var runtimeObjects []runtime.Object
		if testCase.namespaceExists {
			runtimeObjects = append(runtimeObjects, generateNamespace(testCase.namespace))
		}
		runtimeObjects = append(runtimeObjects, generatePVC(testCase.name, testCase.namespace))

		fakeClient := egiClients.GetTestClients(egiClients.TestClientParams{
			K8sMockObjects: runtimeObjects,
		})
		err := cleanPVCs(testCase.namespace, fakeClient)
		assert.NoError(t, err)

		if testCase.namespaceExists {
			// check if namespace has no pvcs left
			pvcs, err := fakeClient.K8sClient.CoreV1().PersistentVolumeClaims(testCase.namespace).List(context.TODO(), metav1.ListOptions{})
			assert.NoError(t, err)
			assert.Equal(t, 0, len(pvcs.Items))
		}
	}
}

func TestCleanPodDisruptionBudget(t *testing.T) {
	generatePodDisruptionBudget := func(name, namespace string) *policyv1.PodDisruptionBudget {
		return &policyv1.PodDisruptionBudget{
			ObjectMeta: metav1.ObjectMeta{
				Name:      name,
				Namespace: namespace,
			},
		}
	}

	testCases := []struct {
		name            string
		namespace       string
		namespaceExists bool
	}{
		{
			name:            "test-pdb",
			namespace:       "test-namespace",
			namespaceExists: true,
		},
		{
			name:            "test-pdb",
			namespace:       "test-namespace",
			namespaceExists: false,
		},
	}

	for _, testCase := range testCases {
		var runtimeObjects []runtime.Object
		if testCase.namespaceExists {
			runtimeObjects = append(runtimeObjects, generateNamespace(testCase.namespace))
		}
		runtimeObjects = append(runtimeObjects, generatePodDisruptionBudget(testCase.name, testCase.namespace))

		fakeClient := egiClients.GetTestClients(egiClients.TestClientParams{
			K8sMockObjects: runtimeObjects,
		})
		err := cleanPodDisruptionBudget(testCase.namespace, fakeClient)
		assert.NoError(t, err)

		if testCase.namespaceExists {
			// check if namespace has no pod disruption budgets left
			pdbs, err := fakeClient.K8sClient.PolicyV1beta1().
				PodDisruptionBudgets(testCase.namespace).List(context.TODO(), metav1.ListOptions{})
			assert.NoError(t, err)
			assert.Equal(t, 0, len(pdbs.Items))
		}
	}
}

func TestCleanResourceQuotas(t *testing.T) {
	generateResourceQuota := func(name, namespace string) *corev1.ResourceQuota {
		return &corev1.ResourceQuota{
			ObjectMeta: metav1.ObjectMeta{
				Name:      name,
				Namespace: namespace,
			},
		}
	}

	testCases := []struct {
		name            string
		namespace       string
		namespaceExists bool
	}{
		{
			name:            "test-resourcequota",
			namespace:       "test-namespace",
			namespaceExists: true,
		},
		{
			name:            "test-resourcequota",
			namespace:       "test-namespace",
			namespaceExists: false,
		},
	}

	for _, testCase := range testCases {
		var runtimeObjects []runtime.Object
		if testCase.namespaceExists {
			runtimeObjects = append(runtimeObjects, generateNamespace(testCase.namespace))
		}
		runtimeObjects = append(runtimeObjects, generateResourceQuota(testCase.name, testCase.namespace))

		fakeClient := egiClients.GetTestClients(egiClients.TestClientParams{
			K8sMockObjects: runtimeObjects,
		})
		err := cleanResourceQuotas(testCase.namespace, fakeClient)
		assert.NoError(t, err)

		if testCase.namespaceExists {
			// check if namespace has no resource quotas left
			resourceQuotas, err := fakeClient.K8sClient.CoreV1().ResourceQuotas(testCase.namespace).List(context.TODO(), metav1.ListOptions{})
			assert.NoError(t, err)
			assert.Equal(t, 0, len(resourceQuotas.Items))
		}
	}
}

func TestCleanNetworkPolicies(t *testing.T) {
	generateNetworkPolicy := func(name, namespace string) *networkingv1.NetworkPolicy {
		return &networkingv1.NetworkPolicy{
			ObjectMeta: metav1.ObjectMeta{
				Name:      name,
				Namespace: namespace,
			},
		}
	}

	testCases := []struct {
		name            string
		namespace       string
		namespaceExists bool
	}{
		{
			name:            "test-networkpolicy",
			namespace:       "test-namespace",
			namespaceExists: true,
		},
		{
			name:            "test-networkpolicy",
			namespace:       "test-namespace",
			namespaceExists: false,
		},
	}

	for _, testCase := range testCases {
		var runtimeObjects []runtime.Object
		if testCase.namespaceExists {
			runtimeObjects = append(runtimeObjects, generateNamespace(testCase.namespace))
		}
		runtimeObjects = append(runtimeObjects, generateNetworkPolicy(testCase.name, testCase.namespace))

		fakeClient := egiClients.GetTestClients(egiClients.TestClientParams{
			K8sMockObjects: runtimeObjects,
		})
		err := cleanNetworkPolicies(testCase.namespace, fakeClient)
		assert.NoError(t, err)

		if testCase.namespaceExists {
			// check if namespace has no network policies left
			networkPolicies, err := fakeClient.K8sClient.
				NetworkingV1().NetworkPolicies(testCase.namespace).List(context.TODO(), metav1.ListOptions{})
			assert.NoError(t, err)
			assert.Equal(t, 0, len(networkPolicies.Items))
		}
	}
}
