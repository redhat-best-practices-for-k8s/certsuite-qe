package pod

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"
)

// DefinePod defines pod manifest based on given params.
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

// RedefinePodWithLabel adds label to given pod manifest.
func RedefinePodWithLabel(pod *corev1.Pod, label map[string]string) {
	newMap := make(map[string]string)
	for k, v := range pod.Labels {
		newMap[k] = v
	}

	for k, v := range label {
		newMap[k] = v
	}

	pod.Labels = newMap
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
