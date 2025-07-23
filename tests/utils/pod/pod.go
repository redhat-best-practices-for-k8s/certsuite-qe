package pod

import (
	"fmt"
	"os"
	"strings"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/utils/ptr"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Package pod provides utilities for creating and configuring Kubernetes pods
// for testing purposes.
//
// Infrastructure Tolerations:
// This package provides functionality to handle infrastructure taints that can
// occur in test/CI environments, such as disk pressure, memory pressure, etc.
//
// Usage Examples:
//
// 1. Basic pod with infrastructure tolerations (recommended for CI):
//    pod := DefinePod("test-pod", "test-ns", "nginx", map[string]string{"app": "test"})
//    RedefineWithInfrastructureTolerationsIfEnabled(pod) // Checks env var automatically
//
// 2. Always apply infrastructure tolerations:
//    pod := DefinePod("test-pod", "test-ns", "nginx", map[string]string{"app": "test"})
//    RedefineWithInfrastructureTolerations(pod)
//
// 3. Custom tolerations for specific scenarios:
//    customTolerations := []corev1.Toleration{{
//        Key:      "custom-taint",
//        Operator: corev1.TolerationOpEqual,
//        Value:    "custom-value",
//        Effect:   corev1.TaintEffectNoExecute,
//    }}
//    RedefineWithCustomTolerations(pod, customTolerations)
//
// Environment Configuration:
// Infrastructure tolerations are now enabled by default.
//
// To disable infrastructure tolerations, set the environment variable `ENABLE_INFRASTRUCTURE_TOLERATIONS` to `false`.
//
// Example:
//    export ENABLE_INFRASTRUCTURE_TOLERATIONS=false
//
// This will cause `RedefineWithInfrastructureTolerationsIfEnabled` to not apply infrastructure tolerations.

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

// RedefineWithInfrastructureTolerations adds tolerations for common infrastructure taints
// that can occur in test/CI environments. This helps improve test reliability when
// nodes have transient resource pressure.
func RedefineWithInfrastructureTolerations(pod *corev1.Pod) {
	infrastructureTolerations := []corev1.Toleration{
		{
			Key:      "node.kubernetes.io/disk-pressure",
			Operator: corev1.TolerationOpExists,
			Effect:   corev1.TaintEffectNoSchedule,
		},
		{
			Key:      "node.kubernetes.io/memory-pressure",
			Operator: corev1.TolerationOpExists,
			Effect:   corev1.TaintEffectNoSchedule,
		},
		{
			Key:      "node.kubernetes.io/pid-pressure",
			Operator: corev1.TolerationOpExists,
			Effect:   corev1.TaintEffectNoSchedule,
		},
		{
			Key:      "node.kubernetes.io/network-unavailable",
			Operator: corev1.TolerationOpExists,
			Effect:   corev1.TaintEffectNoSchedule,
		},
		{
			Key:      "node.kubernetes.io/unreachable",
			Operator: corev1.TolerationOpExists,
			Effect:   corev1.TaintEffectNoExecute,
			// Tolerate for a short time to allow for transient network issues
			TolerationSeconds: ptr.To[int64](30),
		},
		{
			Key:      "node.kubernetes.io/not-ready",
			Operator: corev1.TolerationOpExists,
			Effect:   corev1.TaintEffectNoExecute,
			// Tolerate for a short time to allow for node startup
			TolerationSeconds: ptr.To[int64](30),
		},
	}

	// Append to existing tolerations rather than replacing them
	pod.Spec.Tolerations = append(pod.Spec.Tolerations, infrastructureTolerations...)
}

// RedefineWithCustomTolerations adds custom tolerations to the pod.
func RedefineWithCustomTolerations(pod *corev1.Pod, tolerations []corev1.Toleration) {
	pod.Spec.Tolerations = append(pod.Spec.Tolerations, tolerations...)
}

// RedefineWithInfrastructureTolerationsIfEnabled conditionally adds infrastructure tolerations
// based on configuration. This is the recommended way to apply infrastructure tolerations.
func RedefineWithInfrastructureTolerationsIfEnabled(pod *corev1.Pod) {
	if shouldEnableInfrastructureTolerations() {
		RedefineWithInfrastructureTolerations(pod)
	}
}

// shouldEnableInfrastructureTolerations checks if infrastructure tolerations should be enabled
// based on environment configuration.
func shouldEnableInfrastructureTolerations() bool {
	enabled := os.Getenv("ENABLE_INFRASTRUCTURE_TOLERATIONS")
	if enabled == "" {
		return true
	}

	return strings.ToLower(enabled) == "true"
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
		// Check if the maps are initialized
		if pod.Spec.Containers[0].Resources.Requests == nil {
			pod.Spec.Containers[0].Resources.Requests = make(map[corev1.ResourceName]resource.Quantity)
		}

		if pod.Spec.Containers[0].Resources.Limits == nil {
			pod.Spec.Containers[0].Resources.Limits = make(map[corev1.ResourceName]resource.Quantity)
		}

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
