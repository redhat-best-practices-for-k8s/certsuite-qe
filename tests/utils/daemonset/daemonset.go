package daemonset

import (
	"fmt"

	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"
)

func DefineDaemonSet(namespace string, image string, label map[string]string) *v1.DaemonSet {
	return &v1.DaemonSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "daemonsetnetworkingput",
			Namespace: namespace},
		Spec: v1.DaemonSetSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: label,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name:   "testpod-",
					Labels: label,
				},
				Spec: corev1.PodSpec{
					NodeSelector:                  label,
					TerminationGracePeriodSeconds: pointer.Int64Ptr(0),
					Containers: []corev1.Container{
						{
							Name:    "test",
							Image:   image,
							Command: []string{"/bin/bash", "-c", "sleep INF"}}}}}}}
}

func RedefineDaemonSetWithNodeSelector(daemonSet *v1.DaemonSet, nodeSelector map[string]string) *v1.DaemonSet {
	daemonSet.Spec.Template.Spec.NodeSelector = nodeSelector

	return daemonSet
}

func RedefineWithPrivilegeAndHostNetwork(daemonSet *v1.DaemonSet) *v1.DaemonSet {
	daemonSet.Spec.Template.Spec.HostNetwork = true

	if daemonSet.Spec.Template.Spec.Containers[0].SecurityContext == nil {
		daemonSet.Spec.Template.Spec.Containers[0].SecurityContext = &corev1.SecurityContext{}
	}

	daemonSet.Spec.Template.Spec.Containers[0].SecurityContext.Privileged = pointer.BoolPtr(true)

	return daemonSet
}

func RedefineWithMultus(daemonSet *v1.DaemonSet, nadName string) *v1.DaemonSet {
	daemonSet.Spec.Template.Annotations = map[string]string{
		"k8s.v1.cni.cncf.io/networks": fmt.Sprintf(`[ { "name": "%s" } ]`, nadName),
	}

	return daemonSet
}
