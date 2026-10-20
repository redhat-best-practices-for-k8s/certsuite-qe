package globalhelper

import (
	"context"

	egiClients "github.com/openshift-kni/eco-goinfra/pkg/clients"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	klog "k8s.io/klog/v2"
)

func CreateResourceQuota(quota *corev1.ResourceQuota) error {
	return createResourceQuota(GetEcoGoinfraClient(), quota)
}

func createResourceQuota(client *egiClients.Settings, quota *corev1.ResourceQuota) error {
	nsExist := namespaceExists(quota.Namespace, client)

	if !nsExist {
		return nil
	}

	_, err := client.ResourceQuotas(quota.Namespace).Create(context.TODO(), quota, metav1.CreateOptions{})

	if k8serrors.IsAlreadyExists(err) {
		klog.V(5).Infof("resource quota %s already exists", quota.Name)

		return nil
	}

	return err
}
