package daemonset

import (
	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"
)

type DaemonSet struct {
	v1.DaemonSet
}

func DefineDaemonSet(namespace string, image string, label map[string]string) *DaemonSet {
	return &DaemonSet{
		v1.DaemonSet{
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
								Command: []string{"/bin/bash", "-c", "sleep INF"}}}}}}}}
}

func (daemonSet *DaemonSet) RedefineDaemonSetWithNodeSelector(nodeSelector map[string]string) {
	daemonSet.Spec.Template.Spec.NodeSelector = nodeSelector
}
