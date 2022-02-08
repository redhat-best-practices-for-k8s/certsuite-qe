package pod

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"
)

func DefinePod(podName string, namespace string, image string, label map[string]string) *corev1.Pod {
	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      podName,
			Namespace: namespace,
			Labels:    label},
		Spec: corev1.PodSpec{
			TerminationGracePeriodSeconds: pointer.Int64Ptr(0),
			Containers: []corev1.Container{
				{
					Name:    "test",
					Image:   image,
					Command: []string{"/bin/bash", "-c", "sleep INF"}}}}}
}
