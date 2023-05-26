package container

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"

	corev1 "k8s.io/api/core/v1"
)

// SelectEngine check what container engine is present on the machine.
func SelectEngine() (string, error) {
	for _, containerEngine := range []string{"docker", "podman"} {
		containerEngineCMD := exec.Command(containerEngine)
		directoryName, _ := path.Split(containerEngineCMD.Path)

		if directoryName != "" {
			if strings.Contains(containerEngineCMD.String(), "docker") {
				err := validateDockerDaemonRunning()
				if err != nil {
					return "", err
				}
			}

			return containerEngine, nil
		}
	}

	return "nil", fmt.Errorf("no container Engine present on host machine")
}

func validateDockerDaemonRunning() error {
	// To make it run on MacOs or Windows
	if _, isNonLinuxEnv := os.LookupEnv("NON_LINUX_ENV"); isNonLinuxEnv {
		return nil
	}

	isDaemonRunning := exec.Command("systemctl", "is-active", "--quiet", "docker")

	if isDaemonRunning.Run() != nil {
		return fmt.Errorf("docker daemon is not active on host")
	}

	return nil
}

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
