package globalhelper

import (
	"context"
	"testing"
	"time"

	"github.com/onsi/gomega"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	networkingv1 "k8s.io/api/networking/v1"
	policyv1 "k8s.io/api/policy/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	k8sfake "k8s.io/client-go/kubernetes/fake"
	"k8s.io/utils/ptr"
)

func TestSetTestK8sAPIClientWiresPolicyAndNetworking(t *testing.T) {
	gomega.RegisterTestingT(t)

	oldPDBInterval := podDisruptionBudgetReadyPollInterval
	podDisruptionBudgetReadyPollInterval = time.Millisecond

	t.Cleanup(func() {
		podDisruptionBudgetReadyPollInterval = oldPDBInterval

		UnsetTestK8sAPIClient()
	})

	UnsetTestK8sAPIClient()
	SetTestK8sAPIClient(k8sfake.NewClientset())

	require.NotNil(t, GetAPIClient().PolicyV1Interface)
	require.NotNil(t, GetAPIClient().NetworkingV1Interface)

	pdb := &policyv1.PodDisruptionBudget{
		ObjectMeta: metav1.ObjectMeta{Name: "pdb", Namespace: "ns", UID: "pdb-uid"},
		Spec: policyv1.PodDisruptionBudgetSpec{
			MinAvailable: ptr.To(intstr.FromInt32(1)),
		},
	}

	require.NoError(t, CreatePodDisruptionBudget(pdb, time.Second))

	created, err := IsPodDisruptionBudgetCreated(pdb.Name, pdb.Namespace)
	require.NoError(t, err)
	assert.True(t, created)

	networkPolicy := &networkingv1.NetworkPolicy{
		ObjectMeta: metav1.ObjectMeta{Name: "np", Namespace: "ns"},
	}

	_, err = GetAPIClient().NetworkPolicies(networkPolicy.Namespace).Create(context.TODO(), networkPolicy, metav1.CreateOptions{})
	require.NoError(t, err)

	list, err := GetAPIClient().NetworkPolicies(networkPolicy.Namespace).List(context.TODO(), metav1.ListOptions{})
	require.NoError(t, err)
	assert.Len(t, list.Items, 1)
}
