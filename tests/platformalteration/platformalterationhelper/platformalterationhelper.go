package platformalterationhelper

import (
	"time"

	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/utils/pointer"
)

const WaitingTime = 5 * time.Minute

// DefineDeploymentWithPriviledgedContainer returns deployment struct with priviledged container.
func DefineDeploymentWithPriviledgedContainer(deploymentName string, namespace string, image string,
	label map[string]string) *v1.Deployment {
	return &v1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      deploymentName,
			Namespace: namespace},
		Spec: v1.DeploymentSpec{
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
							Command: []string{"/bin/bash", "-c", "sleep INF"},
							SecurityContext: &corev1.SecurityContext{
								Privileged:   pointer.Bool(true),
								RunAsUser:    pointer.Int64(0),
								Capabilities: &corev1.Capabilities{Add: []corev1.Capability{"ALL"}},
							},
						},
					},
				},
			},
		},
	}
}

// DefineStatefulSetWithPriviledgedContainer returns statefulset struct with priviledged container.
func DefineStatefulSetWithPriviledgedContainer(statefulSetName string, namespace string,
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
							Command: []string{"/bin/bash", "-c", "sleep INF"},
							SecurityContext: &corev1.SecurityContext{
								Privileged:   pointer.Bool(true),
								RunAsUser:    pointer.Int64(0),
								Capabilities: &corev1.Capabilities{Add: []corev1.Capability{"ALL"}},
							},
						},
					},
				},
			},
		},
	}
}
