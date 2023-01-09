package pod

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"
)

// DefinePod defines pod manifest based on given params.
func DefinePod(podName string, namespace string, image string, label map[string]string) *corev1.Pod {
	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      podName,
			Namespace: namespace,
			Labels:    label},
		Spec: corev1.PodSpec{
			TerminationGracePeriodSeconds: pointer.Int64(0),
			Containers: []corev1.Container{
				{
					Name:    "test",
					Image:   image,
					Command: []string{"/bin/bash", "-c", "sleep INF"}}}}}
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

// RedefineWithLivenessProbe adds liveness probe to pod manifest.
func RedefineWithLivenessProbe(pod *corev1.Pod) {
	for index := range pod.Spec.Containers {
		pod.Spec.Containers[index].LivenessProbe = &corev1.Probe{
			ProbeHandler: corev1.ProbeHandler{
				Exec: &corev1.ExecAction{
					Command: []string{"ls"},
				},
			},
		}
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

func RedefineWithPVC(pod *corev1.Pod, name string, claim string) {
	pod.Spec.Volumes = []corev1.Volume{
		{
			Name: name,
			VolumeSource: corev1.VolumeSource{
				PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
					ClaimName: claim,
				},
			},
		},
	}
}

func RedefineWithCPUResources(pod *corev1.Pod, limit string, req string) {
	for i := range pod.Spec.Containers {
		pod.Spec.Containers[i].Resources = corev1.ResourceRequirements{
			Limits: corev1.ResourceList{
				corev1.ResourceCPU: resource.MustParse(limit),
			},
			Requests: corev1.ResourceList{
				corev1.ResourceCPU: resource.MustParse(req),
			},
		}
	}
}

func RedefineWithRunTimeClass(pod *corev1.Pod, rtcName string) {
	pod.Spec.RuntimeClassName = pointer.String(rtcName)
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

// RedefineWithPodantiAffinity redefines pod with podAntiAffinity spec.
func RedefineWithPodantiAffinity(put *corev1.Pod, label map[string]string) {
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
		pod.Spec.Containers[i].Resources.Requests[corev1.ResourceHugePagesPrefix+"2Mi"] = hugepagesVal
		pod.Spec.Containers[i].Resources.Limits[corev1.ResourceHugePagesPrefix+"2Mi"] = hugepagesVal
	}
}

func RedefineFirstContainerWith2MiHugepages(pod *corev1.Pod, hugepages int) error {
	hugepagesVal := resource.MustParse(fmt.Sprintf("%d%s", hugepages, "Mi"))

	if len(pod.Spec.Containers) > 0 {
		pod.Spec.Containers[0].Resources.Requests[corev1.ResourceHugePagesPrefix+"2Mi"] = hugepagesVal
		pod.Spec.Containers[0].Resources.Limits[corev1.ResourceHugePagesPrefix+"2Mi"] = hugepagesVal

		return nil
	}

	return fmt.Errorf("pod %s does not have enough containers", pod.Name)
}

func RedefineSecondContainerWith1GHugepages(pod *corev1.Pod, hugepages int) error {
	hugepagesVal := resource.MustParse(fmt.Sprintf("%d%s", hugepages, "Gi"))

	if len(pod.Spec.Containers) > 1 {
		pod.Spec.Containers[1].Resources.Requests[corev1.ResourceHugePagesPrefix+"1Gi"] = hugepagesVal
		pod.Spec.Containers[1].Resources.Limits[corev1.ResourceHugePagesPrefix+"1Gi"] = hugepagesVal

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
