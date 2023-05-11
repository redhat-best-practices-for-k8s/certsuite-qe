package service

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// DefineService defines service struct.
func DefineService(name string,
	namespace string,
	port int32,
	targetPort int32,
	protocol corev1.Protocol,
	labels map[string]string,
	ipFamilies []corev1.IPFamily,
	ipFamilyPolicy *corev1.IPFamilyPolicy) *corev1.Service {
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: corev1.ServiceSpec{
			Selector: labels,
			Ports: []corev1.ServicePort{
				{
					Protocol: protocol,
					Port:     port,
					TargetPort: intstr.IntOrString{
						Type:   intstr.Int,
						IntVal: targetPort,
					},
				},
			},
			IPFamilies:     ipFamilies,
			IPFamilyPolicy: ipFamilyPolicy,
		},
	}
}

// RedefineWithNodePort redefines service struct with NodePort.
func RedefineWithNodePort(testService *corev1.Service) (*corev1.Service, error) {
	testService.Spec.Type = "NodePort"
	if len(testService.Spec.Ports) < 1 {
		return nil, fmt.Errorf("service does not have available ports")
	}

	testService.Spec.Ports[0].NodePort = testService.Spec.Ports[0].Port

	return testService, nil
}
