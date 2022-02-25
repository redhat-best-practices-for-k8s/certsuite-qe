package globalhelper

import (
	"context"
	"fmt"

	"github.com/golang/glog"
	olmv1 "github.com/operator-framework/api/pkg/operators/v1"
	"github.com/operator-framework/api/pkg/operators/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/test-network-function/cnfcert-tests-verification/tests/utils/namespaces"
)

func DeployOperator(namespace string, operatorGroup *olmv1.OperatorGroup, subscription *v1alpha1.Subscription) error {
	err := namespaces.Create(namespace, APIClient)
	if err != nil {
		return err
	}
	err = APIClient.Create(context.TODO(),
		&olmv1.OperatorGroup{
			ObjectMeta: metav1.ObjectMeta{
				Name:      operatorGroup.Name,
				Namespace: operatorGroup.Namespace},
			Spec: olmv1.OperatorGroupSpec{
				TargetNamespaces: operatorGroup.Spec.TargetNamespaces},
		},
	)
	if err != nil {
		return fmt.Errorf("can not deploy operatorGroup %w", err)
	}
	err = APIClient.Create(context.TODO(),
		&v1alpha1.Subscription{
			ObjectMeta: metav1.ObjectMeta{
				Name:      subscription.Name,
				Namespace: subscription.Namespace,
			},
			Spec: &v1alpha1.SubscriptionSpec{
				Channel:                subscription.Spec.Channel,
				Package:                subscription.Spec.Package,
				CatalogSource:          subscription.Spec.CatalogSource,
				CatalogSourceNamespace: subscription.Spec.CatalogSourceNamespace,
			},
		},
	)
	if err != nil {
		return fmt.Errorf("can not install Subscription %w", err)
	}
	return nil
}

// IsOperatorInstalled validates if the given operator is deployed on the given cluster.
func IsOperatorInstalled(namespace string, operatorDeploymentName string) error {
	glog.V(5).Info(fmt.Sprintf("Validate that operator namespace: %s exists", namespace))

	namespaceExists, err := namespaces.Exists(namespace, APIClient)
	if !namespaceExists && err == nil {
		return fmt.Errorf("operator namespace %s doesn't exist", namespace)
	}

	glog.V(5).Info(fmt.Sprintf("Validate that operator's deployment %s exists", operatorDeploymentName))
	operatorInstalled, err := IsDeploymentInstalled(
		APIClient, namespace, operatorDeploymentName)

	if err != nil {
		return err
	}

	if !operatorInstalled {
		return fmt.Errorf("%s operator's deployment is not installed", operatorDeploymentName)
	}

	return nil
}
