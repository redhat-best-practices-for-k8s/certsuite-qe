package statefulset

import (
	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"
)

// DefineStatefulSet returns statefulset struct.
func DefineStatefulSet(statefulSetName string, namespace string,
	image string, label map[string]string) *v1.StatefulSet {
	return &v1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      statefulSetName,
			Namespace: namespace},
		Spec: v1.StatefulSetSpec{
			Replicas: pointer.Int32Ptr(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: label,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name:   "testpod-",
					Labels: label,
				},
				Spec: corev1.PodSpec{
					TerminationGracePeriodSeconds: pointer.Int64Ptr(0),
					Containers: []corev1.Container{
						{
							Name:    "test",
							Image:   image,
							Command: []string{"/bin/bash", "-c", "sleep INF"}}}}}}}
}

// RedefineWithReadinessProbe adds readiness probe to statefulSet manifest.
func RedefineWithReadinessProbe(statefulSet *v1.StatefulSet) {
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
func RedefineWithLivenessProbe(statefulSet *v1.StatefulSet) {
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
func RedefineWithStartUpProbe(statefulSet *v1.StatefulSet) {
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

func RedefineWithContainerSpecs(statefulSet *v1.StatefulSet, containerSpecs []corev1.Container) {
	statefulSet.Spec.Template.Spec.Containers = containerSpecs
}

func RedefineWithReplicaNumber(statefulSet *v1.StatefulSet, replicasNumber int32) {
	statefulSet.Spec.Replicas = pointer.Int32Ptr(replicasNumber)
}

func RedefineWithPriviledgedContainer(statefulSet *v1.StatefulSet) {
	for index := range statefulSet.Spec.Template.Spec.Containers {
		statefulSet.Spec.Template.Spec.Containers[index].SecurityContext = &corev1.SecurityContext{
			Privileged: pointer.Bool(true),
			RunAsUser:  pointer.Int64(0),
			Capabilities: &corev1.Capabilities{
				Add: []corev1.Capability{"ALL"}},
		}
	}
}

// RedefineWithPostStart adds postStart to statefulSet manifest.
func RedefineWithPostStart(statefulSet *v1.StatefulSet) {
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
