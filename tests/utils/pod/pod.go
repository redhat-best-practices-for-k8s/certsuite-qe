package pod

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"
)

func DefinePod(podName string, namespace string, image string) *corev1.Pod {
	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      podName,
			Namespace: namespace},
		Spec: corev1.PodSpec{
			TerminationGracePeriodSeconds: pointer.Int64Ptr(0),
			Containers: []corev1.Container{
				{
					Name:    "test",
					Image:   image,
					Command: []string{"/bin/bash", "-c", "sleep INF"}}}}}
}

func RedefinePodWithLabel(pod *corev1.Pod, label map[string]string) *corev1.Pod {
	pod.ObjectMeta.Labels = label

	return pod
}

func RedefineWithLivenessProbe(pod *corev1.Pod) *corev1.Pod {
	for index := range pod.Spec.Containers {
		pod.Spec.Containers[index].LivenessProbe = &corev1.Probe{
			Handler: corev1.Handler{
				Exec: &corev1.ExecAction{
					Command: []string{"ls"},
				},
			},
		}
	}

	return pod
}
