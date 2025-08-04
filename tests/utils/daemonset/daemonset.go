package daemonset

import (
	"fmt"

	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/utils/infra"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"
)

func DefineDaemonSet(namespace string, image string, label map[string]string, name string) *appsv1.DaemonSet {
	daemonSet := &appsv1.DaemonSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace},
		Spec: appsv1.DaemonSetSpec{
			MinReadySeconds: 3,
			Selector: &metav1.LabelSelector{
				MatchLabels: label,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name:   "testpod-",
					Labels: label,
				},
				Spec: corev1.PodSpec{
					TerminationGracePeriodSeconds: ptr.To[int64](0),
					Containers: []corev1.Container{
						{
							Name:    "test",
							Image:   image,
							Command: []string{"/bin/bash", "-c", "sleep INF"}}}}}}}

	// Automatically add infrastructure tolerations if enabled
	RedefineWithInfrastructureTolerationsIfEnabled(daemonSet)

	return daemonSet
}

// DefineDaemonSetWithContainerSpecs returns k8s statefulset using configurable
// labels and container specs.
func DefineDaemonSetWithContainerSpecs(name, namespace string, labels map[string]string,
	containerSpecs []corev1.Container) *appsv1.DaemonSet {
	daemonSet := &appsv1.DaemonSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace},
		Spec: appsv1.DaemonSetSpec{
			MinReadySeconds: 3,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					TerminationGracePeriodSeconds: ptr.To[int64](0),
					Containers:                    containerSpecs,
				},
			},
		},
	}

	// Automatically add infrastructure tolerations if enabled
	RedefineWithInfrastructureTolerationsIfEnabled(daemonSet)

	return daemonSet
}

func RedefineDaemonSetWithNodeSelector(daemonSet *appsv1.DaemonSet, nodeSelector map[string]string) {
	daemonSet.Spec.Template.Spec.NodeSelector = nodeSelector
}

// RedefineWithInfrastructureTolerations adds tolerations for common infrastructure taints
// that can occur in test/CI environments. This helps improve test reliability when
// nodes have transient resource pressure.
func RedefineWithInfrastructureTolerations(daemonSet *appsv1.DaemonSet) {
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
	daemonSet.Spec.Template.Spec.Tolerations = append(daemonSet.Spec.Template.Spec.Tolerations, infrastructureTolerations...)
}

// RedefineWithCustomTolerations adds custom tolerations to the daemonset.
func RedefineWithCustomTolerations(daemonSet *appsv1.DaemonSet, tolerations []corev1.Toleration) {
	daemonSet.Spec.Template.Spec.Tolerations = append(daemonSet.Spec.Template.Spec.Tolerations, tolerations...)
}

// RedefineWithInfrastructureTolerationsIfEnabled conditionally adds infrastructure tolerations
// based on configuration. This is the recommended way to apply infrastructure tolerations.
func RedefineWithInfrastructureTolerationsIfEnabled(daemonSet *appsv1.DaemonSet) {
	if infra.ShouldEnableInfrastructureTolerations() {
		RedefineWithInfrastructureTolerations(daemonSet)
	}
}

func RedefineWithLabel(daemonSet *appsv1.DaemonSet, label map[string]string) {
	newMap := make(map[string]string)
	for k, v := range daemonSet.Spec.Template.Labels {
		newMap[k] = v
	}

	for k, v := range label {
		newMap[k] = v
	}

	daemonSet.Spec.Template.Labels = newMap
}

func RedefineWithPrivilegeAndHostNetwork(daemonSet *appsv1.DaemonSet) {
	daemonSet.Spec.Template.Spec.HostNetwork = true

	if daemonSet.Spec.Template.Spec.Containers[0].SecurityContext == nil {
		daemonSet.Spec.Template.Spec.Containers[0].SecurityContext = &corev1.SecurityContext{}
	}

	daemonSet.Spec.Template.Spec.Containers[0].SecurityContext = &corev1.SecurityContext{
		Privileged: ptr.To[bool](true),
		RunAsUser:  ptr.To[int64](0),
		Capabilities: &corev1.Capabilities{
			Add: []corev1.Capability{"ALL"}},
	}
}

func RedefineWithMultus(daemonSet *appsv1.DaemonSet, nadName string) {
	daemonSet.Spec.Template.Annotations = map[string]string{
		"k8s.v1.cni.cncf.io/networks": fmt.Sprintf(`[ { "name": "%s" } ]`, nadName),
	}
}

func RedefineWithImagePullPolicy(daemonSet *appsv1.DaemonSet, pullPolicy corev1.PullPolicy) {
	for index := range daemonSet.Spec.Template.Spec.Containers {
		daemonSet.Spec.Template.Spec.Containers[index].ImagePullPolicy = pullPolicy
	}
}

func RedefineWithContainerSpecs(daemonSet *appsv1.DaemonSet, containerSpecs []corev1.Container) {
	daemonSet.Spec.Template.Spec.Containers = containerSpecs
}

func RedefineWithPrivilegedContainer(daemonSet *appsv1.DaemonSet) {
	for index := range daemonSet.Spec.Template.Spec.Containers {
		daemonSet.Spec.Template.Spec.Containers[index].SecurityContext = &corev1.SecurityContext{
			Privileged: ptr.To[bool](true),
			RunAsUser:  ptr.To[int64](0),
			Capabilities: &corev1.Capabilities{
				Add: []corev1.Capability{"ALL"}},
		}
	}
}

func RedefineWithVolumeMount(daemonSet *appsv1.DaemonSet) {
	for index := range daemonSet.Spec.Template.Spec.Containers {
		daemonSet.Spec.Template.Spec.Containers[index].VolumeMounts = []corev1.VolumeMount{
			{
				Name:      "host",
				MountPath: "/host",
			},
		}
	}

	hostPathType := corev1.HostPathDirectory
	daemonSet.Spec.Template.Spec.Volumes = []corev1.Volume{
		{
			Name: "host",
			VolumeSource: corev1.VolumeSource{
				HostPath: &corev1.HostPathVolumeSource{
					Path: "/",
					Type: &hostPathType,
				},
			},
		},
	}
}

func RedefineWithCPUResources(daemonSet *appsv1.DaemonSet, limit string, req string) {
	for i := range daemonSet.Spec.Template.Spec.Containers {
		daemonSet.Spec.Template.Spec.Containers[i].Resources = corev1.ResourceRequirements{
			Limits: corev1.ResourceList{
				corev1.ResourceCPU: resource.MustParse(limit),
			},
			Requests: corev1.ResourceList{
				corev1.ResourceCPU: resource.MustParse(req),
			},
		}
	}
}

func RedefineWithRunTimeClass(daemonSet *appsv1.DaemonSet, rtcName string) {
	daemonSet.Spec.Template.Spec.RuntimeClassName = ptr.To[string](rtcName)
}
