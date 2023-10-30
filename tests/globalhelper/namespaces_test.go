package globalhelper

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	policyv1 "k8s.io/api/policy/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	k8sfake "k8s.io/client-go/kubernetes/fake"
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

		client := k8sfake.NewSimpleClientset(runtimeObjects...)
		err := deleteNamespaceAndWait(client.CoreV1(), testCase.name, 10*time.Second)
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

		client := k8sfake.NewSimpleClientset(runtimeObjects...)
		err := createNamespace(testCase.name, client.CoreV1())
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

		client := k8sfake.NewSimpleClientset(runtimeObjects...)
		exists, err := namespaceExists(testCase.name, client)
		assert.NoError(t, err)
		assert.Equal(t, testCase.alreadyExists, exists)
	}
}

func TestCleanPods(t *testing.T) {
	generatePod := func(name, namespace string) *corev1.Pod {
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
			name:            "test-pod",
			namespace:       "test-namespace",
			namespaceExists: true,
		},
		{
			name:            "test-pod",
			namespace:       "test-namespace",
			namespaceExists: false,
		},
	}

	for _, testCase := range testCases {
		var runtimeObjects []runtime.Object
		if testCase.namespaceExists {
			runtimeObjects = append(runtimeObjects, generateNamespace(testCase.namespace))
		}
		runtimeObjects = append(runtimeObjects, generatePod(testCase.name, testCase.namespace))

		client := k8sfake.NewSimpleClientset(runtimeObjects...)
		err := cleanPods(testCase.namespace, client)
		assert.NoError(t, err)

		if testCase.namespaceExists {
			// check if namespace has no pods left
			pods, err := client.CoreV1().Pods(testCase.namespace).List(context.TODO(), metav1.ListOptions{})
			assert.NoError(t, err)
			assert.Equal(t, 0, len(pods.Items))
		}
	}
}

func TestCleanDeployments(t *testing.T) {
	generateDeployment := func(name, namespace string) *appsv1.Deployment {
		return &appsv1.Deployment{
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
			name:            "test-deployment",
			namespace:       "test-namespace",
			namespaceExists: true,
		},
		{
			name:            "test-deployment",
			namespace:       "test-namespace",
			namespaceExists: false,
		},
	}

	for _, testCase := range testCases {
		var runtimeObjects []runtime.Object
		if testCase.namespaceExists {
			runtimeObjects = append(runtimeObjects, generateNamespace(testCase.namespace))
		}
		runtimeObjects = append(runtimeObjects, generateDeployment(testCase.name, testCase.namespace))

		client := k8sfake.NewSimpleClientset(runtimeObjects...)
		err := cleanDeployments(testCase.namespace, client)
		assert.NoError(t, err)

		if testCase.namespaceExists {
			// check if namespace has no deployments left
			deployments, err := client.AppsV1().Deployments(testCase.namespace).List(context.TODO(), metav1.ListOptions{})
			assert.NoError(t, err)
			assert.Equal(t, 0, len(deployments.Items))
		}
	}
}

func TestCleanDaemonsets(t *testing.T) {
	generateDaemonset := func(name, namespace string) *appsv1.DaemonSet {
		return &appsv1.DaemonSet{
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
			name:            "test-daemonset",
			namespace:       "test-namespace",
			namespaceExists: true,
		},
		{
			name:            "test-daemonset",
			namespace:       "test-namespace",
			namespaceExists: false,
		},
	}

	for _, testCase := range testCases {
		var runtimeObjects []runtime.Object
		if testCase.namespaceExists {
			runtimeObjects = append(runtimeObjects, generateNamespace(testCase.namespace))
		}
		runtimeObjects = append(runtimeObjects, generateDaemonset(testCase.name, testCase.namespace))

		client := k8sfake.NewSimpleClientset(runtimeObjects...)
		err := cleanDaemonSets(testCase.namespace, client)
		assert.NoError(t, err)

		if testCase.namespaceExists {
			// check if namespace has no daemonsets left
			daemonsets, err := client.AppsV1().DaemonSets(testCase.namespace).List(context.TODO(), metav1.ListOptions{})
			assert.NoError(t, err)
			assert.Equal(t, 0, len(daemonsets.Items))
		}
	}
}

func TestCleanReplicaSets(t *testing.T) {
	generateReplicaSet := func(name, namespace string) *appsv1.ReplicaSet {
		return &appsv1.ReplicaSet{
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
			name:            "test-replicaset",
			namespace:       "test-namespace",
			namespaceExists: true,
		},
		{
			name:            "test-replicaset",
			namespace:       "test-namespace",
			namespaceExists: false,
		},
	}

	for _, testCase := range testCases {
		var runtimeObjects []runtime.Object
		if testCase.namespaceExists {
			runtimeObjects = append(runtimeObjects, generateNamespace(testCase.namespace))
		}
		runtimeObjects = append(runtimeObjects, generateReplicaSet(testCase.name, testCase.namespace))

		client := k8sfake.NewSimpleClientset(runtimeObjects...)
		err := cleanReplicaSets(testCase.namespace, client)
		assert.NoError(t, err)

		if testCase.namespaceExists {
			// check if namespace has no replicasets left
			replicasets, err := client.AppsV1().ReplicaSets(testCase.namespace).List(context.TODO(), metav1.ListOptions{})
			assert.NoError(t, err)
			assert.Equal(t, 0, len(replicasets.Items))
		}
	}
}

func TestCleanServices(t *testing.T) {
	generateService := func(name, namespace string) *corev1.Service {
		return &corev1.Service{
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
			name:            "test-service",
			namespace:       "test-namespace",
			namespaceExists: true,
		},
		{
			name:            "test-service",
			namespace:       "test-namespace",
			namespaceExists: false,
		},
	}

	for _, testCase := range testCases {
		var runtimeObjects []runtime.Object
		if testCase.namespaceExists {
			runtimeObjects = append(runtimeObjects, generateNamespace(testCase.namespace))
		}
		runtimeObjects = append(runtimeObjects, generateService(testCase.name, testCase.namespace))

		client := k8sfake.NewSimpleClientset(runtimeObjects...)

		// get list of services
		services, err := client.CoreV1().Services(testCase.namespace).List(context.TODO(), metav1.ListOptions{})
		assert.NoError(t, err)
		assert.Equal(t, 1, len(services.Items))

		err = cleanServices(testCase.namespace, client)
		assert.NoError(t, err)

		if testCase.namespaceExists {
			// check if namespace has no services left
			services, err := client.CoreV1().Services(testCase.namespace).List(context.TODO(), metav1.ListOptions{})
			assert.NoError(t, err)
			assert.Equal(t, 0, len(services.Items))
		}
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

		client := k8sfake.NewSimpleClientset(runtimeObjects...)
		err := cleanPVCs(testCase.namespace, client)
		assert.NoError(t, err)

		if testCase.namespaceExists {
			// check if namespace has no pvcs left
			pvcs, err := client.CoreV1().PersistentVolumeClaims(testCase.namespace).List(context.TODO(), metav1.ListOptions{})
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

		client := k8sfake.NewSimpleClientset(runtimeObjects...)
		err := cleanPodDisruptionBudget(testCase.namespace, client)
		assert.NoError(t, err)

		if testCase.namespaceExists {
			// check if namespace has no pod disruption budgets left
			pdbs, err := client.PolicyV1beta1().PodDisruptionBudgets(testCase.namespace).List(context.TODO(), metav1.ListOptions{})
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

		client := k8sfake.NewSimpleClientset(runtimeObjects...)
		err := cleanResourceQuotas(testCase.namespace, client)
		assert.NoError(t, err)

		if testCase.namespaceExists {
			// check if namespace has no resource quotas left
			resourceQuotas, err := client.CoreV1().ResourceQuotas(testCase.namespace).List(context.TODO(), metav1.ListOptions{})
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

		client := k8sfake.NewSimpleClientset(runtimeObjects...)
		err := cleanNetworkPolicies(testCase.namespace, client)
		assert.NoError(t, err)

		if testCase.namespaceExists {
			// check if namespace has no network policies left
			networkPolicies, err := client.NetworkingV1().NetworkPolicies(testCase.namespace).List(context.TODO(), metav1.ListOptions{})
			assert.NoError(t, err)
			assert.Equal(t, 0, len(networkPolicies.Items))
		}
	}
}
