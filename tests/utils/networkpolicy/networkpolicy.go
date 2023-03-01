package networkpolicy

import (
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func DefineDenyAllNetworkPolicy(name, ns string, policyTypes []networkingv1.PolicyType,
	labels map[string]string) *networkingv1.NetworkPolicy {
	return &networkingv1.NetworkPolicy{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: ns,
		},
		Spec: networkingv1.NetworkPolicySpec{
			PolicyTypes: policyTypes,
			PodSelector: metav1.LabelSelector{
				MatchLabels: labels,
			},
		},
	}
}

func DefinePolicyTypes(types []string) []networkingv1.PolicyType {
	var policyTypes []networkingv1.PolicyType
	for _, item := range types {
		policyTypes = append(policyTypes, networkingv1.PolicyType(item))
	}

	return policyTypes
}
