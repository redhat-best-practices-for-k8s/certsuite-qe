package globalhelper

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func GetServicesFromNamespace(namespace string) (services []corev1.Service, err error) {
	return getServicesFromNamespace(namespace, GetAPIClient().K8sClient)
}

func getServicesFromNamespace(namespace string, client kubernetes.Interface) ([]corev1.Service, error) {
	services, err := client.CoreV1().Services(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	return services.Items, nil
}
