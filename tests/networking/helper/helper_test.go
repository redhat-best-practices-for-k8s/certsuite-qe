package helper

import (
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

func TestDefineDpdkPod(t *testing.T) {
	testCases := []struct {
		name      string
		namespace string
		expected  *corev1.Pod
	}{
		{
			name:      "test-pod",
			namespace: "test-namespace",
			expected: &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "test-pod",
					Namespace: "test-namespace",
					Labels: map[string]string{
						"app": "test-pod",
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "app-container",
							Image: "registry.redhat.io/openshift4/dpdk-base-rhel8:v4.9",
							Ports: []corev1.ContainerPort{
								{
									Name:          "port1",
									ContainerPort: 80,
								},
							},
							Resources: corev1.ResourceRequirements{
								Limits: corev1.ResourceList{
									"intel.com/intel_sriov_netdevice": resource.MustParse("1"),
								},
							},
						},
					},
					NodeSelector: nil,
					HostNetwork:  false,
					HostPID:      false,
					HostIPC:      false,
				},
			},
		},
	}

	for _, testCase := range testCases {
		pod := DefineDpdkPod(testCase.name, testCase.namespace)
		assert.Equal(t, testCase.expected.Spec.Containers[0].Name, pod.Spec.Containers[0].Name)
		assert.Equal(t, testCase.expected.Spec.Containers[0].Image, pod.Spec.Containers[0].Image)
		assert.Equal(t, testCase.expected.Spec.NodeSelector, pod.Spec.NodeSelector)
		assert.Equal(t, testCase.expected.Spec.HostNetwork, pod.Spec.HostNetwork)
		assert.Equal(t, testCase.expected.Spec.HostPID, pod.Spec.HostPID)
		assert.Equal(t, testCase.expected.Spec.HostIPC, pod.Spec.HostIPC)
	}
}
