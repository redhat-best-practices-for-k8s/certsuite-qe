package service

import (
	"fmt"
	k8sv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// DefineService defines service struct
func DefineService(name string,
	namespace string,
	port int32,
	targetPort int32,
	protocol k8sv1.Protocol,
	labels map[string]string) *k8sv1.Service {

	return &k8sv1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: k8sv1.ServiceSpec{
			Selector: labels,
			Ports: []k8sv1.ServicePort{
				{
					Protocol: protocol,
					Port:     port,
					TargetPort: intstr.IntOrString{
						Type:   intstr.Int,
						IntVal: targetPort,
					},
				},
			},
		},
	}
}

// RedefineWithNodePort redefines service struct with NodePort
func RedefineWithNodePort(testService *k8sv1.Service) (*k8sv1.Service, error) {
	testService.Spec.Type = "NodePort"
	if len(testService.Spec.Ports) < 1 {
		return nil, fmt.Errorf("service does not have available ports")
	}
	testService.Spec.Ports[0].NodePort = testService.Spec.Ports[0].Port
	return testService, nil
}
