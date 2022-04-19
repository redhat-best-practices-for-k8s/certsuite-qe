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

// RedefineWithReplicaNumber redefines statefulSet with requested replica number.
func RedefineWithReplicaNumber(statefulSet *v1.StatefulSet, replicasNumber int32) *v1.StatefulSet {
	statefulSet.Spec.Replicas = pointer.Int32Ptr(replicasNumber)

	return statefulSet
}

func RedefineWithLivenessProbe(statefulSet *v1.StatefulSet) *v1.StatefulSet {
	for index := range statefulSet.Spec.Template.Spec.Containers {
		statefulSet.Spec.Template.Spec.Containers[index].LivenessProbe = &corev1.Probe{
			Handler: corev1.Handler{
				Exec: &corev1.ExecAction{
					Command: []string{"ls"},
				},
			},
		}
	}

	return statefulSet
}
