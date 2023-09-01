package replicaset

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"
)

// DefineReplicaSet returns replicaset struct.
func DefineReplicaSet(replicaSetName string, namespace string, image string, label map[string]string) *appsv1.ReplicaSet {
	return &appsv1.ReplicaSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      replicaSetName,
			Namespace: namespace},
		Spec: appsv1.ReplicaSetSpec{
			Replicas:        pointer.Int32(1),
			MinReadySeconds: 30,
			Selector: &metav1.LabelSelector{
				MatchLabels: label,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name:   "testpod-",
					Labels: label,
				},
				Spec: corev1.PodSpec{
					TerminationGracePeriodSeconds: pointer.Int64(0),
					Containers: []corev1.Container{
						{
							Name:    "test",
							Image:   image,
							Command: []string{"/bin/bash", "-c", "sleep INF"}}}}}}}
}

// RedefineWithReplicaNumber redefines replicaSet with requested replica number.
func RedefineWithReplicaNumber(replicaSet *appsv1.ReplicaSet, replicasNumber int32) {
	replicaSet.Spec.Replicas = pointer.Int32(replicasNumber)
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
