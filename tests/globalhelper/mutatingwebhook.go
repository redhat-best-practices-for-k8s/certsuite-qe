package globalhelper

import (
	"context"
	"fmt"

	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	admissionregistrationtypedv1 "k8s.io/client-go/kubernetes/typed/admissionregistration/v1"
	klog "k8s.io/klog/v2"
)

func DeleteMutatingWebhookConfiguration(name string) error {
	return deleteMutatingWebhookConfiguration(GetAPIClient().K8sClient.AdmissionregistrationV1(), name)
}

func deleteMutatingWebhookConfiguration(client admissionregistrationtypedv1.AdmissionregistrationV1Interface, name string) error {
	err := client.MutatingWebhookConfigurations().Delete(
		context.TODO(),
		name,
		metav1.DeleteOptions{},
	)
	if k8serrors.IsNotFound(err) {
		klog.V(5).Info(fmt.Sprintf("mutating webhook configuration %s is not found", name))

		return nil
	} else if err != nil {
		return fmt.Errorf("failed to delete mutating webhook configuration %q: %w", name, err)
	}

	return nil
}
