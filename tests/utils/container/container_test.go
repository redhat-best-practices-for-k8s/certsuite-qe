package container

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
)

func TestCreateContainerSpecsFromContainerPorts(t *testing.T) {
	testCases := []struct {
		testName string
		ports    []corev1.ContainerPort
		image    string
		name     string
	}{
		{
			testName: "test1",
			ports: []corev1.ContainerPort{
				{
					Name:          "testName1",
					ContainerPort: 8080,
				},
			},
			image: "testImage1",
			name:  "testName1",
		},
		{
			testName: "test2",
			ports:    []corev1.ContainerPort{},
			image:    "testImage2",
			name:     "testName2",
		},
	}

	for _, testCase := range testCases {
		containerSpecs := CreateContainerSpecsFromContainerPorts(testCase.ports, testCase.image, testCase.name)
		assert.NotNil(t, containerSpecs)
		assert.Equal(t, len(testCase.ports), len(containerSpecs))

		for index, containerSpec := range containerSpecs {
			assert.Equal(t, fmt.Sprintf("%s-%d", testCase.name, index), containerSpec.Name)
			assert.Equal(t, testCase.image, containerSpec.Image)
			assert.Equal(t, []string{"/bin/bash", "-c", "sleep INF"}, containerSpec.Command)
			assert.Equal(t, []corev1.ContainerPort{testCase.ports[index]}, containerSpec.Ports)
		}
	}
}
