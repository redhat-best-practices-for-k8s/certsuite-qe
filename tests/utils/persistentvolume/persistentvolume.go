package persistentvolume

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// DefinePersistentVolume defines a persistent volume manifest based on given params.
func DefinePersistentVolume(pvName, pvcName, pvcNamespace string) *corev1.PersistentVolume {
	return &corev1.PersistentVolume{
		ObjectMeta: metav1.ObjectMeta{
			Name: pvName,
		},
		Spec: corev1.PersistentVolumeSpec{
			AccessModes: []corev1.PersistentVolumeAccessMode{
				corev1.ReadWriteMany,
			},
			ClaimRef: &corev1.ObjectReference{
				Namespace: pvcNamespace,
				Name:      pvcName,
			},

			PersistentVolumeReclaimPolicy: corev1.PersistentVolumeReclaimRetain,
			PersistentVolumeSource:        corev1.PersistentVolumeSource{Local: &corev1.LocalVolumeSource{Path: "/tmp"}},
			NodeAffinity: &corev1.VolumeNodeAffinity{Required: &corev1.NodeSelector{
				NodeSelectorTerms: []corev1.NodeSelectorTerm{{
					MatchExpressions: []corev1.NodeSelectorRequirement{
						{
							Key:      "kubernetes.io/hostname",
							Operator: corev1.NodeSelectorOpExists,
						},
					},
				},
				},
			},
			},
			Capacity: corev1.ResourceList{corev1.ResourceStorage: resource.MustParse("10Gi")},
		},
	}
}

func RedefineWithPVReclaimPolicy(pv *corev1.PersistentVolume, policy corev1.PersistentVolumeReclaimPolicy) {
	pv.Spec.PersistentVolumeReclaimPolicy = policy
}

func RedefineWithStorageClass(pv *corev1.PersistentVolume, storageClassName string) {
	pv.Spec.StorageClassName = storageClassName
}
