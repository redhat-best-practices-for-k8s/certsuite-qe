package resourcequota

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func DefineResourceQuota(resourceQuotaName, cpuRequest, memoryRequest, cpuLimit, memoryLimit string) *corev1.ResourceQuota {
	return &corev1.ResourceQuota{
		ObjectMeta: metav1.ObjectMeta{
			Name: resourceQuotaName,
		},
		Spec: corev1.ResourceQuotaSpec{
			Hard: corev1.ResourceList{
				"requests.cpu":    resource.MustParse(cpuRequest),
				"requests.memory": resource.MustParse(memoryRequest),
				"limits.cpu":      resource.MustParse(cpuLimit),
				"limits.memory":   resource.MustParse(memoryLimit),
			},
		},
	}
}
