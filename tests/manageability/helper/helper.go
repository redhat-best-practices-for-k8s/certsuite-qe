package helper

import (
	"k8s.io/utils/pointer"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func DefineManageabilityPod(podName, namespace, image string, label map[string]string) *corev1.Pod {
	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      podName,
			Namespace: namespace,
			Labels:    label},
		Spec: corev1.PodSpec{
			TerminationGracePeriodSeconds: pointer.Int64(0),
			Containers: []corev1.Container{
				{
					Name:  "httpd",
					Image: image,
					Ports: []corev1.ContainerPort{
						{
							Name:          "http",
							ContainerPort: *pointer.Int32(80),
						},
					},
				},
			},
		},
	}
}

func RedefinePodWithContainerPort(pod *corev1.Pod, containerIndex int, portName string) {
	totalContainers := len(pod.Spec.Containers)

	if containerIndex >= 0 && containerIndex < totalContainers {
		container := pod.Spec.Containers[containerIndex]
		container.Ports[0].Name = portName
	}
}
