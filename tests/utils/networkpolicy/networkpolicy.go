package networkpolicy

import (
	v1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func DefineDenyAllNetworkPolicy(name, ns string, policyTypes []v1.PolicyType, labels map[string]string) *v1.NetworkPolicy {
	return &v1.NetworkPolicy{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: ns,
		},
		Spec: v1.NetworkPolicySpec{
			PolicyTypes: policyTypes,
			PodSelector: metav1.LabelSelector{
				MatchLabels: labels,
			},
		},
	}
}

func DefinePolicyTypes(types []string) []v1.PolicyType {
	var policyTypes []v1.PolicyType
	for _, item := range types {
		policyTypes = append(policyTypes, v1.PolicyType(item))
	}

	return policyTypes
}
