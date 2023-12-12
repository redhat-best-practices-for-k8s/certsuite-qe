package container

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
)

func CreateContainerSpecsFromContainerPorts(ports []corev1.ContainerPort, image, name string) []corev1.Container {
	numContainers := len(ports)
	containerSpecs := []corev1.Container{}

	for index := 0; index < numContainers; index++ {
		containerSpecs = append(containerSpecs,
			corev1.Container{
				Name:    fmt.Sprintf("%s-%d", name, index),
				Image:   image,
				Command: []string{"/bin/bash", "-c", "sleep INF"},
				Ports:   []corev1.ContainerPort{ports[index]},
			},
		)
	}

	return containerSpecs
}
