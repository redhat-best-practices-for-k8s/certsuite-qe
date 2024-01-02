package globalhelper

import (
	"context"
	"fmt"

	"github.com/golang/glog"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func CreateResourceQuota(quota *corev1.ResourceQuota) error {
	return createResourceQuota(GetAPIClient().K8sClient, quota)
}

func createResourceQuota(client kubernetes.Interface, quota *corev1.ResourceQuota) error {
	nsExist, err := namespaceExists(quota.Namespace, client)
	if err != nil {
		return err
	}

	if !nsExist {
		return nil
	}

	_, err1 := client.CoreV1().ResourceQuotas(quota.Namespace).Create(context.TODO(), quota, metav1.CreateOptions{})

	if k8serrors.IsAlreadyExists(err1) {
		glog.V(5).Info(fmt.Sprintf("resource quota %s already exists", quota.Name))

		return nil
	}

	return err1
}
