package deployment

import (
	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"
)

// DefineDeployment returns deployment struct.
func DefineDeployment(namespace string, image string, label map[string]string) *v1.Deployment {
	return &v1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "networkingput",
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
							Command: []string{"/bin/bash", "-c", "sleep INF"}}}}}}}
}

// RedefineWithContainersSecurityContextAll redefines deployment with extended permissions.
func RedefineWithContainersSecurityContextAll(deployment *v1.Deployment) *v1.Deployment {
	for index := range deployment.Spec.Template.Spec.Containers {
		deployment.Spec.Template.Spec.Containers[index].SecurityContext = &corev1.SecurityContext{
			Capabilities: &corev1.Capabilities{
				Add: []corev1.Capability{"ALL"}},
		}
	}

	return deployment
}

// RedefineWithLabels redefines deployment with additional label.
func RedefineWithLabels(deployment *v1.Deployment, label map[string]string) *v1.Deployment {
	newMap := make(map[string]string)
	for k, v := range deployment.Spec.Template.Labels {
		newMap[k] = v
	}

	for k, v := range label {
		newMap[k] = v
	}

	deployment.Spec.Template.Labels = newMap

	return deployment
}

// RedefineWithReplicaNumber redefines deployment with requested replica number.
func RedefineWithReplicaNumber(deployment *v1.Deployment, replicasNumber int32) *v1.Deployment {
	deployment.Spec.Replicas = pointer.Int32Ptr(replicasNumber)

	return deployment
}

// RedefineWithPreStopSpec redefines deployment with requested lifecycle/preStop spec.
func RedefineWithPreStopSpec(deployment *v1.Deployment, command []string) *v1.Deployment {
	for index := range deployment.Spec.Template.Spec.Containers {
		deployment.Spec.Template.Spec.Containers[index].Lifecycle = &corev1.Lifecycle{
			PreStop: &corev1.Handler{
				Exec: &corev1.ExecAction{
					Command: command,
				},
			},
		}
	}

	return deployment
}
