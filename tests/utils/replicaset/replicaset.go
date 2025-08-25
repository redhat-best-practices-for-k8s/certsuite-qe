package replicaset

import (
	"github.com/redhat-best-practices-for-k8s/certsuite-qe/tests/utils/infra"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"
)

// DefineReplicaSet returns replicaset struct.
func DefineReplicaSet(replicaSetName string, namespace string, image string, label map[string]string) *appsv1.ReplicaSet {
	replicaSet := &appsv1.ReplicaSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      replicaSetName,
			Namespace: namespace},
		Spec: appsv1.ReplicaSetSpec{
			Replicas:        ptr.To[int32](1),
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
	RedefineWithInfrastructureTolerationsIfEnabled(replicaSet)

	return replicaSet
}

// RedefineWithReplicaNumber redefines replicaSet with requested replica number.
func RedefineWithReplicaNumber(replicaSet *appsv1.ReplicaSet, replicasNumber int32) {
	replicaSet.Spec.Replicas = ptr.To[int32](replicasNumber)
}

func RedefineWithPVC(replicaSet *appsv1.ReplicaSet, name string, claim string) {
	replicaSet.Spec.Template.Spec.Volumes = []corev1.Volume{
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

// RedefineWithInfrastructureTolerations adds tolerations for common infrastructure taints
// that can occur in test/CI environments. This helps improve test reliability when
// nodes have transient resource pressure.
func RedefineWithInfrastructureTolerations(replicaSet *appsv1.ReplicaSet) {
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
	replicaSet.Spec.Template.Spec.Tolerations = append(replicaSet.Spec.Template.Spec.Tolerations, infrastructureTolerations...)
}

// RedefineWithCustomTolerations adds custom tolerations to the replicaset.
func RedefineWithCustomTolerations(replicaSet *appsv1.ReplicaSet, tolerations []corev1.Toleration) {
	replicaSet.Spec.Template.Spec.Tolerations = append(replicaSet.Spec.Template.Spec.Tolerations, tolerations...)
}

// RedefineWithInfrastructureTolerationsIfEnabled conditionally adds infrastructure tolerations
// based on configuration. This is the recommended way to apply infrastructure tolerations.
func RedefineWithInfrastructureTolerationsIfEnabled(replicaSet *appsv1.ReplicaSet) {
	if infra.ShouldEnableInfrastructureTolerations() {
		RedefineWithInfrastructureTolerations(replicaSet)
	}
}
