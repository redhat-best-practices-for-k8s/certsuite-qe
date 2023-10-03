package statefulset

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"
)

// DefineStatefulSet returns statefulset struct.
func DefineStatefulSet(statefulSetName string, namespace string,
	image string, label map[string]string) *appsv1.StatefulSet {
	return &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      statefulSetName,
			Namespace: namespace},
		Spec: appsv1.StatefulSetSpec{
			MinReadySeconds: 30,
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
