package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
)

func TestDefineService(t *testing.T) {
	singleStack := corev1.IPFamilyPolicySingleStack
	testService := DefineService("testService", "testNamespace", 80, 8080, "TCP",
		map[string]string{"app": "test"}, []corev1.IPFamily{corev1.IPv4Protocol}, &singleStack)
	assert.NotNil(t, testService)
	assert.Equal(t, "testService", testService.Name)
	assert.Equal(t, "testNamespace", testService.Namespace)
	assert.Equal(t, int32(80), testService.Spec.Ports[0].Port)
	assert.Equal(t, int32(8080), testService.Spec.Ports[0].TargetPort.IntVal)
	assert.Equal(t, "test", testService.Spec.Selector["app"])
	assert.Equal(t, corev1.IPv4Protocol, testService.Spec.IPFamilies[0])
	assert.Equal(t, corev1.IPFamilyPolicySingleStack, *testService.Spec.IPFamilyPolicy)
}

func TestRedefineWithNodePort(t *testing.T) {
	singleStack := corev1.IPFamilyPolicySingleStack
	testService := DefineService("testService", "testNamespace", 80, 8080, "TCP",
		map[string]string{"app": "test"}, []corev1.IPFamily{corev1.IPv4Protocol}, &singleStack)
	testService, err := RedefineWithNodePort(testService)
	assert.Nil(t, err)
	assert.NotNil(t, testService)
	assert.Equal(t, corev1.ServiceTypeNodePort, testService.Spec.Type)
	assert.Equal(t, int32(80), testService.Spec.Ports[0].NodePort)
}
