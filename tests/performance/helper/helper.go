package helper

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/utils/pointer"

	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func DefineExclusivePod(podName string, namespace string, image string, label map[string]string) *corev1.Pod {
	cpuLimit := "1"
	memoryLimit := "512Mi"
	containerCommand := []string{"/bin/bash", "-c", "sleep INF"}

	containerResource := corev1.ResourceRequirements{
		Limits: corev1.ResourceList{
			corev1.ResourceCPU:    resource.MustParse(cpuLimit),
			corev1.ResourceMemory: resource.MustParse(memoryLimit),
		},
		Requests: corev1.ResourceList{
			corev1.ResourceCPU:    resource.MustParse(cpuLimit),
			corev1.ResourceMemory: resource.MustParse(memoryLimit),
		},
	}

	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      podName,
			Namespace: namespace,
			Labels:    label},
		Spec: corev1.PodSpec{
			TerminationGracePeriodSeconds: pointer.Int64(0),
			SecurityContext: &corev1.PodSecurityContext{
				RunAsUser:    pointer.Int64(1000),
				RunAsGroup:   pointer.Int64(1000),
				RunAsNonRoot: pointer.Bool(true)},
			Containers: []corev1.Container{
				{
					Name:      "shared",
					Image:     image,
					Command:   containerCommand,
					Resources: containerResource},
				{
					Name:      "exclusive",
					Image:     image,
					Command:   containerCommand,
					Resources: containerResource},
			},
		},
	}
}

func RedefinePodWithSharedContainer(pod *corev1.Pod, containerIndex int) {
	totalContainers := len(pod.Spec.Containers)
	limit := "1"
	req := "0.75"

	if containerIndex >= 0 && containerIndex < totalContainers {
		pod.Spec.Containers[containerIndex].Resources = corev1.ResourceRequirements{
			Limits: corev1.ResourceList{
				corev1.ResourceCPU: resource.MustParse(limit),
			},
			Requests: corev1.ResourceList{
				corev1.ResourceCPU: resource.MustParse(req),
			},
		}
	}
}
