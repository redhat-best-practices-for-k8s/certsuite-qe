package statefulset

import (
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/utils/infra"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"
)

// DefineStatefulSet returns statefulset struct.
func DefineStatefulSet(statefulSetName string, namespace string,
	image string, label map[string]string) *appsv1.StatefulSet {
	statefulSet := &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      statefulSetName,
			Namespace: namespace},
		Spec: appsv1.StatefulSetSpec{
			MinReadySeconds: 3,
			Replicas:        ptr.To[int32](1),
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
	RedefineWithInfrastructureTolerationsIfEnabled(statefulSet)

	return statefulSet
}

// RedefineWithReadinessProbe adds readiness probe to statefulSet manifest.
func RedefineWithReadinessProbe(statefulSet *appsv1.StatefulSet) {
	for index := range statefulSet.Spec.Template.Spec.Containers {
		statefulSet.Spec.Template.Spec.Containers[index].ReadinessProbe = &corev1.Probe{
			ProbeHandler: corev1.ProbeHandler{
				Exec: &corev1.ExecAction{
					Command: []string{"ls"},
				},
			},
		}
	}
}

// RedefineWithLivenessProbe adds liveness probe to statefulSet manifest.
func RedefineWithLivenessProbe(statefulSet *appsv1.StatefulSet) {
	for index := range statefulSet.Spec.Template.Spec.Containers {
		statefulSet.Spec.Template.Spec.Containers[index].LivenessProbe = &corev1.Probe{
			ProbeHandler: corev1.ProbeHandler{
				Exec: &corev1.ExecAction{
					Command: []string{"ls"},
				},
			},
		}
	}
}

// RedefineWithStartUpProbe adds startup probe to statefulSet manifest.
func RedefineWithStartUpProbe(statefulSet *appsv1.StatefulSet) {
	for index := range statefulSet.Spec.Template.Spec.Containers {
		statefulSet.Spec.Template.Spec.Containers[index].StartupProbe = &corev1.Probe{
			ProbeHandler: corev1.ProbeHandler{
				Exec: &corev1.ExecAction{
					Command: []string{"ls"},
				},
			},
		}
	}
}

// RedefineWithInfrastructureTolerations adds tolerations for common infrastructure taints
// that can occur in test/CI environments. This helps improve test reliability when
// nodes have transient resource pressure.
func RedefineWithInfrastructureTolerations(statefulSet *appsv1.StatefulSet) {
	infrastructureTolerations := []corev1.Toleration{
		{
			Key:      "node.kubernetes.io/disk-pressure",
			Operator: corev1.TolerationOpExists,
			Effect:   corev1.TaintEffectNoSchedule,
		},
		{
			Key:      "node.kubernetes.io/disk-pressure",
			Operator: corev1.TolerationOpExists,
			Effect:   corev1.TaintEffectNoExecute,
			// Tolerate for a reasonable time to allow disk cleanup
			TolerationSeconds: ptr.To[int64](300),
		},
		{
			Key:      "node.kubernetes.io/memory-pressure",
			Operator: corev1.TolerationOpExists,
			Effect:   corev1.TaintEffectNoSchedule,
		},
		{
			Key:      "node.kubernetes.io/memory-pressure",
			Operator: corev1.TolerationOpExists,
			Effect:   corev1.TaintEffectNoExecute,
			// Tolerate for a shorter time for memory pressure
			TolerationSeconds: ptr.To[int64](120),
		},
		{
			Key:      "node.kubernetes.io/pid-pressure",
			Operator: corev1.TolerationOpExists,
			Effect:   corev1.TaintEffectNoSchedule,
		},
		{
			Key:      "node.kubernetes.io/pid-pressure",
			Operator: corev1.TolerationOpExists,
			Effect:   corev1.TaintEffectNoExecute,
			// Tolerate briefly for PID pressure
			TolerationSeconds: ptr.To[int64](60),
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
	statefulSet.Spec.Template.Spec.Tolerations = append(statefulSet.Spec.Template.Spec.Tolerations, infrastructureTolerations...)
}

// RedefineWithCustomTolerations adds custom tolerations to the statefulset.
func RedefineWithCustomTolerations(statefulSet *appsv1.StatefulSet, tolerations []corev1.Toleration) {
	statefulSet.Spec.Template.Spec.Tolerations = append(statefulSet.Spec.Template.Spec.Tolerations, tolerations...)
}

// RedefineWithInfrastructureTolerationsIfEnabled conditionally adds infrastructure tolerations
// based on configuration. This is the recommended way to apply infrastructure tolerations.
func RedefineWithInfrastructureTolerationsIfEnabled(statefulSet *appsv1.StatefulSet) {
	if infra.ShouldEnableInfrastructureTolerations() {
		RedefineWithInfrastructureTolerations(statefulSet)
	}
}

func RedefineWithContainerSpecs(statefulSet *appsv1.StatefulSet, containerSpecs []corev1.Container) {
	statefulSet.Spec.Template.Spec.Containers = containerSpecs
}

func RedefineWithReplicaNumber(statefulSet *appsv1.StatefulSet, replicasNumber int32) {
	statefulSet.Spec.Replicas = ptr.To[int32](replicasNumber)
}

func RedefineWithPrivilegedContainer(statefulSet *appsv1.StatefulSet) {
	for index := range statefulSet.Spec.Template.Spec.Containers {
		statefulSet.Spec.Template.Spec.Containers[index].SecurityContext = &corev1.SecurityContext{
			Privileged: ptr.To[bool](true),
			RunAsUser:  ptr.To[int64](0),
			Capabilities: &corev1.Capabilities{
				Add: []corev1.Capability{"ALL"}},
		}
	}
}

// RedefineWithPostStart adds postStart to statefulSet manifest.
func RedefineWithPostStart(statefulSet *appsv1.StatefulSet) {
	for index := range statefulSet.Spec.Template.Spec.Containers {
		statefulSet.Spec.Template.Spec.Containers[index].Lifecycle = &corev1.Lifecycle{
			PostStart: &corev1.LifecycleHandler{
				Exec: &corev1.ExecAction{
					Command: []string{"ls"},
				},
			},
		}
	}
}

func RedefineWithContainersSecurityContextCaps(sts *appsv1.StatefulSet, add, drop []string) {
	var addedCaps, droppedCaps []corev1.Capability

	for _, cap := range add {
		addedCaps = append(addedCaps, corev1.Capability(cap))
	}

	for _, cap := range drop {
		droppedCaps = append(droppedCaps, corev1.Capability(cap))
	}

	for index := range sts.Spec.Template.Spec.Containers {
		sts.Spec.Template.Spec.Containers[index].SecurityContext = &corev1.SecurityContext{
			Privileged: ptr.To[bool](true),
			RunAsUser:  ptr.To[int64](0),
			Capabilities: &corev1.Capabilities{
				Add:  addedCaps,
				Drop: droppedCaps,
			},
		}
	}
}
