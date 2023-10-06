package globalhelper

import (
	"context"

	"github.com/operator-framework/api/pkg/operators/v1alpha1"
	v1alpha1typed "github.com/operator-framework/operator-lifecycle-manager/pkg/api/client/clientset/versioned/typed/operators/v1alpha1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func CreateInstallPlan(plan *v1alpha1.InstallPlan) error {
	return createInstallPlan(plan, GetAPIClient().OperatorsV1alpha1Interface)
}

func createInstallPlan(plan *v1alpha1.InstallPlan, opclient v1alpha1typed.OperatorsV1alpha1Interface) error {
	_, err := opclient.InstallPlans(plan.Namespace).Create(context.Background(), plan, metav1.CreateOptions{})
	if k8serrors.IsAlreadyExists(err) {
		return nil
	}

	return err
}
