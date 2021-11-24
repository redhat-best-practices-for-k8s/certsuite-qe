package deployment

import (
	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"
)

type Deployment struct {
	v1.Deployment
}

// DefineDeployment returns deployment struct
func DefineDeployment(namespace string, image string, label map[string]string) *Deployment {
	return &Deployment{
		v1.Deployment{
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
								Command: []string{"/bin/bash", "-c", "sleep INF"}}}}}}}}
}

// RedefineWithContainersSecurityContextAll redefines deployment with extended permissions
func (deployment *Deployment) RedefineWithContainersSecurityContextAll() {
	for index := range deployment.Spec.Template.Spec.Containers {
		deployment.Spec.Template.Spec.Containers[index].SecurityContext = &corev1.SecurityContext{
			Capabilities: &corev1.Capabilities{
				Add: []corev1.Capability{"ALL"}},
		}
	}
}

// RedefineWithLabels redefines deployment with additional label
func (deployment *Deployment) RedefineWithLabels(label map[string]string) {
	newMap := make(map[string]string)
	for k, v := range deployment.Spec.Template.Labels {
		newMap[k] = v
	}
	for k, v := range label {
		newMap[k] = v
	}
	deployment.Spec.Template.Labels = newMap
}

// RedefineWithReplicaNumber redefines deployment with requested replica number
func (deployment *Deployment) RedefineWithReplicaNumber(replicasNumber int32) {
	deployment.Spec.Replicas = pointer.Int32Ptr(replicasNumber)
}
