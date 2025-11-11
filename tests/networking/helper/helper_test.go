package helper

import (
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
)

func TestFindListIntersections(t *testing.T) {
	testCases := []struct {
		list1    []string
		list2    []string
		expected []string
	}{
		{
			list1:    []string{"a", "b", "c"},
			list2:    []string{"a", "b", "c"},
			expected: []string{"a", "b", "c"},
		},
		{
			list1:    []string{"a", "b", "c"},
			list2:    []string{"a", "b", "c", "d"},
			expected: []string{"a", "b", "c"},
		},
		{
			list1:    []string{"a", "b", "c", "d"},
			list2:    []string{"a", "b", "c"},
			expected: []string{"a", "b", "c"},
		},
		{
			list1:    []string{"a", "b", "c", "d"},
			list2:    []string{"e", "f", "g", "h"},
			expected: nil,
		},
	}

	for _, testCase := range testCases {
		assert.Equal(t, testCase.expected, findListIntersections(testCase.list1, testCase.list2))
	}
}

func TestCreateContainerSpecsFromContainerPorts(t *testing.T) {
	testCases := []struct {
		ports         []corev1.ContainerPort
		expectedPorts []corev1.ContainerPort
	}{
		{
			ports: []corev1.ContainerPort{
				{
					Name:          "port1",
					ContainerPort: 80,
				},
				{
					Name:          "port2",
					ContainerPort: 443,
				},
			},
			expectedPorts: []corev1.ContainerPort{
				{
					Name:          "port1",
					ContainerPort: 80,
				},
				{
					Name:          "port2",
					ContainerPort: 443,
				},
			},
		},
		{
			ports: []corev1.ContainerPort{
				{
					Name:          "port1",
					ContainerPort: 80,
				},
				{
					Name:          "port2",
					ContainerPort: 443,
				},
				{
					Name:          "port3",
					ContainerPort: 8080,
				},
			},
			expectedPorts: []corev1.ContainerPort{
				{
					Name:          "port1",
					ContainerPort: 80,
				},
				{
					Name:          "port2",
					ContainerPort: 443,
				},
				{
					Name:          "port3",
					ContainerPort: 8080,
				},
			},
		},
	}

	for _, testCase := range testCases {
		containers := createContainerSpecsFromContainerPorts(testCase.ports)
		for i, container := range containers {
			assert.Equal(t, testCase.expectedPorts[i], container.Ports[0])
		}
	}
}
