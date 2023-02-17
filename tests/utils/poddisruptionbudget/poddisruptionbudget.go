package poddisruptionbudget

import (
	policyv1 "k8s.io/api/policy/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func DefinePodDisruptionBudgetMinAvailable(name string, namespace string, minAvailable intstr.IntOrString,
	label map[string]string) *policyv1.PodDisruptionBudget {
	return &policyv1.PodDisruptionBudget{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: policyv1.PodDisruptionBudgetSpec{
			MinAvailable: &minAvailable,
			Selector: &metav1.LabelSelector{
				MatchLabels: label,
			},
		},
	}
}

func DefinePodDisruptionBudgetMaxUnAvailable(name string, namespace string, maxUnAvailable intstr.IntOrString,
	label map[string]string) *policyv1.PodDisruptionBudget {
	return &policyv1.PodDisruptionBudget{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: policyv1.PodDisruptionBudgetSpec{
			MaxUnavailable: &maxUnAvailable,
			Selector: &metav1.LabelSelector{
				MatchLabels: label,
			},
		},
	}
}
