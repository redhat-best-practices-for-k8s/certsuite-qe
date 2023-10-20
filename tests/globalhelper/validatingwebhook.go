package globalhelper

import (
	"context"
	"fmt"

	"github.com/golang/glog"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	admissionregistrationtypedv1 "k8s.io/client-go/kubernetes/typed/admissionregistration/v1"
)

func DeleteValidatingWebhookConfiguration(name string) error {
	return deleteValidatingWebhookConfiguration(GetAPIClient().K8sClient.AdmissionregistrationV1(), name)
}

func deleteValidatingWebhookConfiguration(client admissionregistrationtypedv1.AdmissionregistrationV1Interface, name string) error {
	err := client.ValidatingWebhookConfigurations().Delete(
		context.TODO(),
		name,
		metav1.DeleteOptions{},
	)
	if k8serrors.IsNotFound(err) {
		glog.V(5).Info(fmt.Sprintf("validating webhook configuration %s is not found", name))

		return nil
	} else if err != nil {
		return fmt.Errorf("failed to delete validating webhook configuration %q: %w", name, err)
	}

	return nil
}
