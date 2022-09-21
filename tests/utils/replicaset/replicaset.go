package replicaset

import (
	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"
)

// DefineReplicaSet returns replicaset struct.
func DefineReplicaSet(replicaSetName string, namespace string, image string, label map[string]string) *v1.ReplicaSet {
	return &v1.ReplicaSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      replicaSetName,
			Namespace: namespace},
		Spec: v1.ReplicaSetSpec{
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

// RedefineWithReplicaNumber redefines replicaSet with requested replica number.
func RedefineWithReplicaNumber(replicaSet *v1.ReplicaSet, replicasNumber int32) *v1.ReplicaSet {
	replicaSet.Spec.Replicas = pointer.Int32Ptr(replicasNumber)

	return replicaSet
}

func RedefineWithPVC(replicaSet *v1.ReplicaSet, name string, claim string) *v1.ReplicaSet {
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

	return replicaSet
}
