package daemonset

import (
	"fmt"

	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	nodev1 "k8s.io/api/node/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"
)

func DefineDaemonSet(namespace string, image string, label map[string]string, name string) *v1.DaemonSet {
	return &v1.DaemonSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
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
					TerminationGracePeriodSeconds: pointer.Int64Ptr(0),
					Containers: []corev1.Container{
						{
							Name:    "test",
							Image:   image,
							Command: []string{"/bin/bash", "-c", "sleep INF"}}}}}}}
}

// DefineDaemonSetWithContainerSpecs returns k8s statefulset using configurable
// labels and container specs.
func DefineDaemonSetWithContainerSpecs(name, namespace string, labels map[string]string,
	containerSpecs []corev1.Container) *v1.DaemonSet {
	return &v1.DaemonSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace},
		Spec: v1.DaemonSetSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Containers: containerSpecs,
				},
			},
		},
	}
}

func RedefineDaemonSetWithNodeSelector(daemonSet *v1.DaemonSet, nodeSelector map[string]string) *v1.DaemonSet {
	daemonSet.Spec.Template.Spec.NodeSelector = nodeSelector

	return daemonSet
}

func RedefineDaemonSetWithLabel(daemonSet *v1.DaemonSet, label map[string]string) *v1.DaemonSet {
	newMap := make(map[string]string)
	for k, v := range daemonSet.Spec.Template.Labels {
		newMap[k] = v
	}

	for k, v := range label {
		newMap[k] = v
	}

	daemonSet.Spec.Template.Labels = newMap

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

func RedefineWithImagePullPolicy(daemonSet *v1.DaemonSet, pullPolicy corev1.PullPolicy) *v1.DaemonSet {
	for index := range daemonSet.Spec.Template.Spec.Containers {
		daemonSet.Spec.Template.Spec.Containers[index].ImagePullPolicy = pullPolicy
	}

	return daemonSet
}

func RedefineWithContainerSpecs(daemonSet *v1.DaemonSet, containerSpecs []corev1.Container) *v1.DaemonSet {
	daemonSet.Spec.Template.Spec.Containers = containerSpecs

	return daemonSet
}

func RedefineWithPriviledgedContainer(daemonSet *v1.DaemonSet) *v1.DaemonSet {
	for index := range daemonSet.Spec.Template.Spec.Containers {
		daemonSet.Spec.Template.Spec.Containers[index].SecurityContext = &corev1.SecurityContext{
			Privileged: pointer.Bool(true),
			RunAsUser:  pointer.Int64(0),
			Capabilities: &corev1.Capabilities{
				Add: []corev1.Capability{"ALL"}},
		}
	}

	return daemonSet
}

func RedefineWithVolumeMount(daemonSet *v1.DaemonSet) *v1.DaemonSet {
	for index := range daemonSet.Spec.Template.Spec.Containers {
		daemonSet.Spec.Template.Spec.Containers[index].VolumeMounts = []corev1.VolumeMount{
			{
				Name:      "host",
				MountPath: "/host",
			},
		}
	}

	hostPathType := corev1.HostPathDirectory
	daemonSet.Spec.Template.Spec.Volumes = []corev1.Volume{
		{
			Name: "host",
			VolumeSource: corev1.VolumeSource{
				HostPath: &corev1.HostPathVolumeSource{
					Path: "/",
					Type: &hostPathType,
				},
			},
		},
	}

	return daemonSet
}

func RedefineWithResources(daemonSet *v1.DaemonSet, limit string, req string) {
	for i := range daemonSet.Spec.Template.Spec.Containers {
		daemonSet.Spec.Template.Spec.Containers[i].Resources = corev1.ResourceRequirements{
			Limits: corev1.ResourceList{
				corev1.ResourceCPU: resource.MustParse(limit),
			},
			Requests: corev1.ResourceList{
				corev1.ResourceCPU: resource.MustParse(req),
			},
		}
	}
}

func RedefineWithRunTimeClass(daemonSet *v1.DaemonSet, rtc *nodev1.RuntimeClass) {
	daemonSet.Spec.Template.Spec.RuntimeClassName = pointer.String(rtc.Name)
}
