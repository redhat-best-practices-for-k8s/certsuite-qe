package utils

import (
	olmv1 "github.com/operator-framework/api/pkg/operators/v1"
	"github.com/operator-framework/api/pkg/operators/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// DefineOperatorGroup returns an operator group struct.
func DefineOperatorGroup(groupName string, namespace string, targetNamespace []string) *olmv1.OperatorGroup {
	return &olmv1.OperatorGroup{
		ObjectMeta: metav1.ObjectMeta{
			Name:      groupName,
			Namespace: namespace},
		Spec: olmv1.OperatorGroupSpec{
			TargetNamespaces: targetNamespace},
	}
}

// DefineSubscription returns a subscription struct.
func DefineSubscription(subName, namespace, channel, operatorName, catalogSource,
	catalogSourceNamespace, startingCSV string, installPlanApproval v1alpha1.Approval) *v1alpha1.Subscription {
	return &v1alpha1.Subscription{
		ObjectMeta: metav1.ObjectMeta{
			Name:      subName,
			Namespace: namespace,
		},
		Spec: &v1alpha1.SubscriptionSpec{
			Channel:                channel,
			Package:                operatorName,
			CatalogSource:          catalogSource,
			CatalogSourceNamespace: catalogSourceNamespace,
			StartingCSV:            startingCSV,
			InstallPlanApproval:    installPlanApproval,
		},
	}
}
