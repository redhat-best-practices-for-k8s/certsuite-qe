package pod

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/utils/ptr"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	HugePages2Mi = "2Mi"
	HugePages1Gi = "1Gi"
)

// DefinePod defines pod manifest based on given params.
func DefinePod(podName string, namespace string, image string, label map[string]string) *corev1.Pod {
	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      podName,
			Namespace: namespace,
			Labels:    label},
		Spec: corev1.PodSpec{
			TerminationGracePeriodSeconds: ptr.To[int64](0),
			SecurityContext: &corev1.PodSecurityContext{
				RunAsUser:    ptr.To[int64](1000),
				RunAsGroup:   ptr.To[int64](1000),
				RunAsNonRoot: ptr.To[bool](true),
			},
			Containers: []corev1.Container{
				{
					Name:    "test",
					Image:   image,
					Command: []string{"/bin/bash", "-c", "sleep INF"}}}}}
}

func RedefineWithServiceAccount(pod *corev1.Pod, serviceAccountName string) {
	pod.Spec.ServiceAccountName = serviceAccountName
}

// RedefineWithReadinessProbe adds readiness probe to given pod manifest.
func RedefineWithReadinessProbe(pod *corev1.Pod) {
	for index := range pod.Spec.Containers {
		pod.Spec.Containers[index].ReadinessProbe = &corev1.Probe{
			ProbeHandler: corev1.ProbeHandler{
				Exec: &corev1.ExecAction{
					Command: []string{"ls"},
				},
			},
		}
	}
}

func RedefinePodContainerWithLivenessProbeCommand(pod *corev1.Pod, index int, commands []string) {
	pod.Spec.Containers[index].LivenessProbe = &corev1.Probe{
		ProbeHandler: corev1.ProbeHandler{
			Exec: &corev1.ExecAction{
				Command: commands,
			},
		},
		InitialDelaySeconds: 5,
		PeriodSeconds:       5,
	}
}

// RedefineWithLivenessProbe adds liveness probe to pod manifest.
func RedefineWithLivenessProbe(pod *corev1.Pod) {
	commands := []string{"ls"}

	for index := range pod.Spec.Containers {
		RedefinePodContainerWithLivenessProbeCommand(pod, index, commands)
	}
}

// RedefineWithStartUpProbe adds startup probe to pod manifest.
func RedefineWithStartUpProbe(pod *corev1.Pod) {
	for index := range pod.Spec.Containers {
		pod.Spec.Containers[index].StartupProbe = &corev1.Probe{
			ProbeHandler: corev1.ProbeHandler{
				Exec: &corev1.ExecAction{
					Command: []string{"ls"},
				},
			},
		}
	}
}

func RedefineWithPVC(pod *corev1.Pod, volumeName string, claimName string) {
	pod.Spec.Volumes = []corev1.Volume{
		{
			Name: volumeName,
			VolumeSource: corev1.VolumeSource{
				PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
					ClaimName: claimName,
				},
			},
		},
	}
}

func RedefineWithCPUResources(pod *corev1.Pod, limit string, req string) {
	for i := range pod.Spec.Containers {
		containerResources := &pod.Spec.Containers[i].Resources

		if containerResources.Requests == nil {
			containerResources.Requests = corev1.ResourceList{}
		}

		if containerResources.Limits == nil {
			containerResources.Limits = corev1.ResourceList{}
		}

		containerResources.Requests[corev1.ResourceCPU] = resource.MustParse(req)
		containerResources.Limits[corev1.ResourceCPU] = resource.MustParse(limit)
	}
}

func RedefineWithMemoryResources(pod *corev1.Pod, limit string, req string) {
	for i := range pod.Spec.Containers {
		containerResources := &pod.Spec.Containers[i].Resources

		if containerResources.Requests == nil {
			containerResources.Requests = corev1.ResourceList{}
		}

		if containerResources.Limits == nil {
			containerResources.Limits = corev1.ResourceList{}
		}

		containerResources.Requests[corev1.ResourceMemory] = resource.MustParse(req)
		containerResources.Limits[corev1.ResourceMemory] = resource.MustParse(limit)
	}
}

func RedefineWithRunTimeClass(pod *corev1.Pod, rtcName string) {
	pod.Spec.RuntimeClassName = ptr.To[string](rtcName)
}

// RedefineWithNodeAffinity redefines pod with nodeAffinity spec.
func RedefineWithNodeAffinity(pod *corev1.Pod, key string) {
	pod.Spec.Affinity = &corev1.Affinity{
		NodeAffinity: &corev1.NodeAffinity{
			RequiredDuringSchedulingIgnoredDuringExecution: &corev1.NodeSelector{
				NodeSelectorTerms: []corev1.NodeSelectorTerm{
					{
						MatchExpressions: []corev1.NodeSelectorRequirement{
							{
								Key:      key,
								Operator: corev1.NodeSelectorOpExists,
							},
						},
					},
				},
			},
		}}
}

// RedefineWithPodAffinity redefines pod with podAffinity spec.
func RedefineWithPodAffinity(put *corev1.Pod, label map[string]string) {
	put.Spec.Affinity = &corev1.Affinity{
		PodAffinity: &corev1.PodAffinity{
			RequiredDuringSchedulingIgnoredDuringExecution: []corev1.PodAffinityTerm{
				{
					LabelSelector: &metav1.LabelSelector{
						MatchLabels: label,
					},
					TopologyKey: "kubernetes.io/hostname",
				},
			},
		}}
}

// RedefineWithPodAntiAffinity redefines pod with podAntiAffinity spec.
func RedefineWithPodAntiAffinity(put *corev1.Pod, label map[string]string) {
	put.Spec.Affinity = &corev1.Affinity{
		PodAntiAffinity: &corev1.PodAntiAffinity{
			RequiredDuringSchedulingIgnoredDuringExecution: []corev1.PodAffinityTerm{
				{
					LabelSelector: &metav1.LabelSelector{
						MatchLabels: label,
					},
					TopologyKey: "kubernetes.io/hostname",
				},
			},
		}}
}

func RedefineWith2MiHugepages(pod *corev1.Pod, hugepages int) {
	hugepagesVal := resource.MustParse(fmt.Sprintf("%d%s", hugepages, "Mi"))

	for i := range pod.Spec.Containers {
		pod.Spec.Containers[i].Resources.Requests[corev1.ResourceHugePagesPrefix+HugePages2Mi] = hugepagesVal
		pod.Spec.Containers[i].Resources.Limits[corev1.ResourceHugePagesPrefix+HugePages2Mi] = hugepagesVal
	}
}

func RedefineWith1GiHugepages(pod *corev1.Pod, hugepages int) {
	hugepagesVal := resource.MustParse(fmt.Sprintf("%d%s", hugepages, "Gi"))

	for i := range pod.Spec.Containers {
		pod.Spec.Containers[i].Resources.Requests[corev1.ResourceHugePagesPrefix+HugePages1Gi] = hugepagesVal
		pod.Spec.Containers[i].Resources.Limits[corev1.ResourceHugePagesPrefix+HugePages1Gi] = hugepagesVal
	}
}

func RedefineFirstContainerWith2MiHugepages(pod *corev1.Pod, hugepages int) error {
	hugepagesVal := resource.MustParse(fmt.Sprintf("%d%s", hugepages, "Mi"))

	if len(pod.Spec.Containers) > 0 {
		pod.Spec.Containers[0].Resources.Requests[corev1.ResourceHugePagesPrefix+HugePages2Mi] = hugepagesVal
		pod.Spec.Containers[0].Resources.Limits[corev1.ResourceHugePagesPrefix+HugePages2Mi] = hugepagesVal

		return nil
	}

	return fmt.Errorf("pod %s does not have enough containers", pod.Name)
}

func RedefineFirstContainerWith1GiHugepages(pod *corev1.Pod, hugepages int) error {
	hugepagesVal := resource.MustParse(fmt.Sprintf("%d%s", hugepages, "Gi"))

	if len(pod.Spec.Containers) > 0 {
		pod.Spec.Containers[0].Resources.Requests[corev1.ResourceHugePagesPrefix+HugePages1Gi] = hugepagesVal
		pod.Spec.Containers[0].Resources.Limits[corev1.ResourceHugePagesPrefix+HugePages1Gi] = hugepagesVal

		return nil
	}

	return fmt.Errorf("pod %s does not have enough containers", pod.Name)
}

func RedefineSecondContainerWith1GHugepages(pod *corev1.Pod, hugepages int) error {
	hugepagesVal := resource.MustParse(fmt.Sprintf("%d%s", hugepages, "Gi"))

	if len(pod.Spec.Containers) > 1 {
		pod.Spec.Containers[1].Resources.Requests[corev1.ResourceHugePagesPrefix+HugePages1Gi] = hugepagesVal
		pod.Spec.Containers[1].Resources.Limits[corev1.ResourceHugePagesPrefix+HugePages1Gi] = hugepagesVal

		return nil
	}

	return fmt.Errorf("pod %s does not have enough containers", pod.Name)
}

// RedefineWithPostStart adds postStart to pod manifest.
func RedefineWithPostStart(pod *corev1.Pod) {
	for index := range pod.Spec.Containers {
		pod.Spec.Containers[index].Lifecycle = &corev1.Lifecycle{
			PostStart: &corev1.LifecycleHandler{
				Exec: &corev1.ExecAction{
					Command: []string{"ls"},
				},
			},
		}
	}
}

func RedefineWithContainerExecCommand(pod *corev1.Pod, commandArgs []string, containerIndex int) error {
	if len(pod.Spec.Containers) <= containerIndex {
		return fmt.Errorf("pod %s does not have enough containers", pod.Name)
	}

	pod.Spec.Containers[containerIndex].Command = commandArgs

	return nil
}
